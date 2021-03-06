// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gate_test

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"gate.computer/gate/image"
	internalbuild "gate.computer/gate/internal/build"
	"gate.computer/gate/packet"
	"gate.computer/gate/runtime"
	"gate.computer/gate/runtime/abi"
	"gate.computer/gate/runtime/container"
	"gate.computer/gate/service"
	"gate.computer/gate/service/origin"
	"gate.computer/gate/snapshot"
	"gate.computer/gate/trap"
	"gate.computer/wag/binding"
	"gate.computer/wag/compile"
	"gate.computer/wag/object"
	objectabi "gate.computer/wag/object/abi"
	"gate.computer/wag/object/stack/stacktrace"
	"gate.computer/wag/wa"
	"github.com/tsavola/mu"
)

const (
	maxTextSize     = 32 * 1024 * 1024
	stackSize       = wa.PageSize
	memorySizeLimit = 128 * 1024 * 1024
)

type executor struct {
	*runtime.Executor
	closed bool
}

func (test *executor) Close() error {
	test.closed = true
	return test.Executor.Close()
}

func newExecutor() (tester *executor) {
	actual, err := runtime.NewExecutor(&runtime.Config{
		Container: container.Config{
			Namespace: testNamespaceConfig,
			ExecDir:   testExecDir,
		},
	})
	if err != nil {
		panic(err)
	}

	tester = &executor{Executor: actual}

	go func() {
		<-tester.Dead()
		time.Sleep(time.Second)
		if !tester.closed {
			time.Sleep(time.Second)
			panic("executor died")
		}
	}()

	return
}

type serviceRegistry struct {
	origin   io.Writer
	originMu *sync.Mutex
}

func (services serviceRegistry) StartServing(ctx context.Context, config runtime.ServiceConfig, _ []snapshot.Service, send chan<- packet.Buf, recv <-chan packet.Buf,
) (runtime.ServiceDiscoverer, []runtime.ServiceState, <-chan error, error) {
	d := new(serviceDiscoverer)

	go func() {
		var originInstance service.Instance

		for p := range recv {
			var name string
			d.nameMu.Guard(func() {
				name = d.names[p.Code()]
			})

			switch name {
			case "origin":
				if originInstance == nil {
					connector := origin.New(nil)
					go func() {
						defer connector.Close()

						if services.originMu != nil {
							services.originMu.Lock()
							defer services.originMu.Unlock()
						}

						if f := connector.Connect(context.Background()); f != nil {
							f(context.Background(), bytes.NewReader(nil), services.origin)
						}
					}()

					var err error
					originInstance, err = connector.CreateInstance(ctx, service.InstanceConfig{
						Service: packet.Service{
							Code:        p.Code(),
							MaxSendSize: config.MaxSendSize,
						},
					}, nil)
					if err != nil {
						panic(err)
					}
					defer originInstance.Shutdown(context.Background())

					if err := originInstance.Ready(ctx); err != nil {
						panic(err)
					}
					if err := originInstance.Start(ctx, send, func(e error) { panic(e) }); err != nil {
						panic(err)
					}
				}

				if err := originInstance.Handle(ctx, send, p); err != nil {
					panic(err)
				}
			}
		}
	}()

	return d, nil, make(chan error), nil
}

type serviceDiscoverer struct {
	services []runtime.ServiceState
	nameMu   mu.Mutex
	names    []string
}

func (d *serviceDiscoverer) Discover(ctx context.Context, names []string) ([]runtime.ServiceState, error) {
	for _, name := range names {
		var s runtime.ServiceState

		switch name {
		case "origin":
			s.SetAvail()
		}

		d.services = append(d.services, s)

		d.nameMu.Guard(func() {
			d.names = append(d.names, name)
		})
	}

	return d.services, nil
}

func (d *serviceDiscoverer) NumServices() int {
	return len(d.services)
}

func (*serviceDiscoverer) Shutdown(context.Context) error {
	return nil
}

func (*serviceDiscoverer) Suspend(context.Context) ([]snapshot.Service, error) {
	return nil, nil
}

var testFS *image.Filesystem

func init() {
	dir := os.Getenv("GATE_TEST_FILESYSTEM")
	if dir == "" {
		d := "testdata/filesystem"
		if _, err := os.Stat(d); err == nil {
			dir = d
		} else if !os.IsNotExist(err) {
			panic(err)
		}
	}

	if dir != "" {
		if err := os.RemoveAll(path.Join(dir, "program")); err != nil {
			panic(err)
		}
		if err := os.RemoveAll(path.Join(dir, "instance")); err != nil {
			panic(err)
		}
		fs, err := image.NewFilesystem(dir)
		if err != nil {
			panic(err)
		}
		testFS = fs
	}
}

func prepareBuild(exec *executor, storage image.Storage, config compile.Config, wasm []byte, moduleSize int, codeMap *object.CallMap,
) (r *bytes.Reader, mod compile.Module, build *image.Build) {
	r = bytes.NewReader(wasm)

	mod, err := compile.LoadInitialSections(&compile.ModuleConfig{Config: config}, r)
	if err != nil {
		panic(err)
	}

	build, err = image.NewBuild(storage, moduleSize, maxTextSize, codeMap, true)
	if err != nil {
		panic(err)
	}

	if err := binding.BindImports(&mod, build.ImportResolver()); err != nil {
		panic(err)
	}

	return
}

func buildInstance(exec *executor, storage image.Storage, codeMap *object.CallMap, wasm []byte, moduleSize int, function string, persistent bool,
) (prog *image.Program, inst *image.Instance, mod compile.Module) {
	var config compile.Config
	var sectionMap image.SectionMap

	if persistent {
		config.SectionMapper = sectionMap.Mapper()
	}

	r, mod, build := prepareBuild(exec, storage, config, wasm, moduleSize, codeMap)
	defer build.Close()

	if persistent {
		if _, err := build.ModuleWriter().Write(wasm); err != nil {
			panic(err)
		}
	}

	var codeConfig = &compile.CodeConfig{
		Text:   build.TextBuffer(),
		Mapper: codeMap,
		Config: config,
	}

	err := compile.LoadCodeSection(codeConfig, r, mod, abi.Library())
	if err != nil {
		panic(err)
	}

	// dump.Text(os.Stderr, codeConfig.Text.Bytes(), 0, codeMap.FuncAddrs, nil)

	entryIndex, err := internalbuild.ResolveEntryFunc(mod, function, false)
	if err != nil {
		panic(err)
	}

	if err := build.FinishText(stackSize, 0, mod.GlobalsSize(), mod.InitialMemorySize()); err != nil {
		panic(err)
	}

	var dataConfig = &compile.DataConfig{
		GlobalsMemory:   build.GlobalsMemoryBuffer(),
		MemoryAlignment: build.MemoryAlignment(),
		Config:          config,
	}

	if err := compile.LoadDataSection(dataConfig, r, mod); err != nil {
		panic(err)
	}

	if persistent {
		err = compile.LoadCustomSections(&config, r)
		if err != nil {
			return
		}
	}

	startIndex := -1
	if index, found := mod.StartFunc(); found {
		startIndex = int(index)
	}

	prog, err = build.FinishProgram(sectionMap, mod, startIndex, true, nil, 0)
	if err != nil {
		panic(err)
	}

	memLimit := mod.MemorySizeLimit()
	if memLimit < 0 || memLimit > memorySizeLimit {
		memLimit = memorySizeLimit
	}

	inst, err = build.FinishInstance(prog, memLimit, entryIndex)
	if err != nil {
		panic(err)
	}

	return
}

func startInstance(ctx context.Context, t *testing.T, storage image.Storage, wasm []byte, function string, debugOut io.Writer,
) (*executor, *image.Program, *image.Instance, *runtime.Process, *object.CallMap, compile.Module) {
	var err error

	executor := newExecutor()
	defer func() {
		if err != nil {
			executor.Close()
		}
	}()

	codeMap := new(object.CallMap)

	prog, inst, mod := buildInstance(executor, storage, codeMap, wasm, len(wasm), function, true)
	defer func() {
		if err != nil {
			prog.Close()
			inst.Close()
		}
	}()

	err = prog.Store(fmt.Sprint(crc32.ChecksumIEEE(wasm)))
	if err != nil {
		t.Fatal(err)
	}

	proc, err := executor.NewProcess(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err != nil {
			proc.Kill()
		}
	}()

	policy := runtime.ProcessPolicy{
		TimeResolution: time.Microsecond,
		DebugLog:       debugOut,
	}

	err = proc.Start(prog, inst, policy)
	if err != nil {
		t.Fatal(err)
	}

	return executor, prog, inst, proc, codeMap, mod
}

func runProgram(t *testing.T, wasm []byte, function string, debug io.Writer, expectTrap trap.ID) (output bytes.Buffer) {
	t.Helper()

	ctx := context.Background()

	executor, prog, inst, proc, textMap, mod := startInstance(ctx, t, image.Memory, wasm, function, debug)
	defer proc.Kill()
	defer inst.Close()
	defer prog.Close()
	defer executor.Close()

	var outputMu sync.Mutex

	result, trapID, err := proc.Serve(ctx, serviceRegistry{&output, &outputMu}, nil)
	if err != nil {
		t.Errorf("run error: %v", err)
	} else {
		if trapID != expectTrap {
			t.Errorf("run %v", trapID)
		}
		if trapID == trap.Exit && result.Value() != runtime.ResultSuccess {
			t.Errorf("run result: %s", result)
		}
		if testing.Verbose() {
			trace, err := inst.Stacktrace(textMap, mod.FuncTypes())
			if err == nil {
				err = stacktrace.Fprint(os.Stderr, trace, mod.FuncTypes(), nil, nil)
			}
			if err != nil {
				t.Error(err)
			}
		}
	}

	outputMu.Lock()
	defer outputMu.Unlock()

	if s := output.String(); len(s) > 0 {
		t.Logf("output: %q", s)
	}
	return
}

func TestRunNop(t *testing.T) {
	runProgram(t, wasmNop, "", nil, trap.Exit)
}

func testRunHello(t *testing.T, debug io.Writer) {
	output := runProgram(t, wasmHello, "greet", debug, trap.Exit)
	if s := output.String(); s != "hello, world\n" {
		t.Errorf("%q", s)
	}
}

func TestRunHello(t *testing.T) {
	testRunHello(t, os.Stdout)
}

func TestRunHelloNoDebug(t *testing.T) {
	testRunHello(t, nil)
}

func TestRunHelloDebug(t *testing.T) {
	var debug bytes.Buffer
	runProgram(t, wasmHelloDebug, "debug", &debug, trap.Exit)
	s := debug.String()
	t.Logf("debug: %q", s)
	if s != "hello, world\n" {
		t.Errorf("%q", s)
	}
}

func TestRunHelloDebugNoDebug(t *testing.T) {
	runProgram(t, wasmHelloDebug, "debug", nil, trap.Exit)
}

func TestRunSuspendMem(t *testing.T) {
	testRunSuspend(t, image.Memory, objectabi.TextAddrNoFunction)
}

func TestRunSuspendFS(t *testing.T) {
	if testFS == nil {
		t.Skip("test filesystem not specified")
	}

	testRunSuspend(t, testFS, objectabi.TextAddrResume)
}

func TestRunSuspendPersistMem(t *testing.T) {
	if testFS == nil {
		t.Skip("test filesystem not specified")
	}

	s := image.CombinedStorage(testFS, image.PersistentMemory(testFS))
	testRunSuspend(t, s, objectabi.TextAddrResume)
}

func testRunSuspend(t *testing.T, storage image.Storage, expectInitRoutine uint32) {
	ctx := context.Background()

	executor, prog, inst, proc, codeMap, mod := startInstance(ctx, t, storage, wasmSuspend, "loop", os.Stdout)
	defer proc.Kill()
	defer inst.Close()
	defer prog.Close()
	defer executor.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	exit, trapID, err := proc.Serve(ctx, serviceRegistry{}, nil)
	if err != nil {
		t.Errorf("run error: %v", err)
	} else if trapID == 0 {
		t.Errorf("run exit: %d", exit)
	} else if trapID != trap.Suspended {
		t.Errorf("run %v", trapID)
	}

	if err := inst.CheckMutation(); err != nil {
		t.Errorf("instance state: %v", err)
	}

	if err := inst.Store(t.Name(), t.Name(), prog); err != nil {
		t.Fatal(err)
	}

	if false {
		trace, err := inst.Stacktrace(codeMap, mod.FuncTypes())
		if err != nil {
			t.Fatal(err)
		}

		if len(trace) > 0 {
			stacktrace.Fprint(os.Stderr, trace, mod.FuncTypes(), nil, nil)
		}
	}
}

func TestRandomSeed(t *testing.T) {
	values := make([][2]uint64, 10)

	for i := 0; i < len(values); i++ {
		var debug bytes.Buffer
		runProgram(t, wasmRandomSeed, "dump", &debug, trap.Exit)
		for j, s := range strings.Split(debug.String(), " ") {
			n, err := strconv.ParseUint(strings.TrimSpace(s), 16, 64)
			if err != nil {
				t.Fatal(err)
			}
			values[i][j] = n
		}
	}

	for i := 0; i < len(values); i++ {
		for j := 0; j < len(values); j++ {
			if i != j && values[i] == values[j] {
				t.Fatal(values[i])
			}
		}
	}
}

func TestRandomDeficiency(t *testing.T) {
	testRandomDeficiency(t, "toomuch")
}

func TestRandomDeficiency2(t *testing.T) {
	testRandomDeficiency(t, "toomuch2")
}

func testRandomDeficiency(t *testing.T, function string) {
	var debug bytes.Buffer
	runProgram(t, wasmRandomSeed, function, &debug, trap.ABIDeficiency)
	if s := debug.String(); s != "ping\n" {
		t.Error(s)
	}
}

func TestTime(t *testing.T) {
	runProgram(t, wasmTime, "check", os.Stderr, trap.Exit)
}
