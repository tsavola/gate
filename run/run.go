package run

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/tsavola/wag"
	"github.com/tsavola/wag/types"
	"github.com/tsavola/wag/wasm"

	"github.com/tsavola/gate/internal/memfd"
)

var (
	pageSize = uint32(os.Getpagesize())
)

func roundToPage(size int) uint32 {
	mask := pageSize - 1
	return (uint32(size) + mask) &^ mask
}

type envFunc struct {
	addr uint64
	sig  types.Function
}

type Environment struct {
	executor string
	loader   *os.File
	funcs    map[string]envFunc
}

func NewEnvironment(executor, loader, loaderSymbols string) (env *Environment, err error) {
	symbolFile, err := os.Open(loaderSymbols)
	if err != nil {
		return
	}
	defer symbolFile.Close()

	funcs := make(map[string]envFunc)

	for {
		var (
			name string
			addr uint64
			n    int
		)

		n, err = fmt.Fscanf(symbolFile, "%x T %s\n", &addr, &name)
		if err != nil {
			if err == io.EOF && n == 0 {
				break
			}
			return
		}
		if n != 2 {
			err = fmt.Errorf("%s: parse error", loaderSymbols)
			return
		}

		switch name {
		case "__gate_get_abi_version", "__gate_get_max_packet_size":
			funcs[name] = envFunc{addr, types.Function{
				Result: types.I32,
			}}

		case "__gate_func_ptr":
			funcs[name] = envFunc{addr, types.Function{
				Args:   []types.T{types.I32},
				Result: types.I32,
			}}

		case "__gate_exit":
			funcs[name] = envFunc{addr, types.Function{
				Args: []types.T{types.I32},
			}}

		case "__gate_recv_full", "__gate_send_full":
			funcs[name] = envFunc{addr, types.Function{
				Args: []types.T{types.I32, types.I32},
			}}
		}
	}

	loaderFile, err := os.Open(loader)
	if err != nil {
		return
	}

	env = &Environment{
		executor: executor,
		loader:   loaderFile,
		funcs:    funcs,
	}
	return
}

func (env *Environment) ImportFunction(module, field string, sig types.Function) (variadic bool, addr uint64, err error) {
	if module == "env" {
		if f, found := env.funcs[field]; found {
			if !f.sig.Equal(sig) {
				err = fmt.Errorf("function %s %s imported with wrong signature: %s", field, f.sig, sig)
				return
			}

			addr = f.addr
			return
		}
	}

	err = fmt.Errorf("imported function not found: %s %s %s", module, field, sig)
	return
}

func (env *Environment) ImportGlobal(module, field string, t types.T) (value uint64, err error) {
	if module == "env" {
		switch field {
		case "__gate_abi_version":
			value = abiVersion
			return

		case "__gate_max_packet_size":
			value = maxPacketSize
			return
		}
	}

	err = fmt.Errorf("imported global not found: %s %s %s", module, field, t)
	return
}

type payloadInfo struct {
	PageSize       uint32
	RODataSize     uint32
	TextSize       uint32
	MemoryOffset   uint32
	InitMemorySize uint32
	GrowMemorySize uint32
	StackSize      uint32
}

type Payload struct {
	maps *os.File
	info payloadInfo
}

func NewPayload(m *wag.Module, growMemorySize wasm.MemorySize, stackSize int32) (payload *Payload, err error) {
	initMemorySize, _ := m.MemoryLimits()

	if initMemorySize > growMemorySize {
		err = fmt.Errorf("initial memory size %d exceeds maximum memory size %d", initMemorySize, growMemorySize)
		return
	}

	roData := m.ROData()
	text := m.Text()
	data, memoryOffset := m.Data()

	fd, err := memfd.Create("payload", memfd.CLOEXEC|memfd.ALLOW_SEALING)
	if err != nil {
		return
	}

	f := os.NewFile(uintptr(fd), "memfd")

	_, err = f.Write(roData)
	if err != nil {
		f.Close()
		return
	}

	roDataSize := roundToPage(len(roData))

	_, err = f.WriteAt(text, int64(roDataSize))
	if err != nil {
		f.Close()
		return
	}

	textSize := roundToPage(len(text))

	_, err = f.WriteAt(data, int64(roDataSize)+int64(textSize))
	if err != nil {
		f.Close()
		return
	}

	globalsMemorySize := roundToPage(memoryOffset + int(growMemorySize))
	totalSize := int64(roDataSize) + int64(textSize) + int64(globalsMemorySize) + int64(stackSize)

	err = f.Truncate(totalSize)
	if err != nil {
		f.Close()
		return
	}

	_, err = memfd.Fcntl(fd, memfd.F_ADD_SEALS, memfd.F_SEAL_SHRINK|memfd.F_SEAL_GROW)
	if err != nil {
		f.Close()
		return
	}

	payload = &Payload{
		maps: f,
		info: payloadInfo{
			PageSize:       pageSize,
			RODataSize:     roDataSize,
			TextSize:       textSize,
			MemoryOffset:   uint32(memoryOffset),
			InitMemorySize: uint32(initMemorySize),
			GrowMemorySize: uint32(growMemorySize),
			StackSize:      uint32(stackSize),
		},
	}
	return
}

func (payload *Payload) Close() (err error) {
	err = payload.maps.Close()
	payload.maps = nil
	return
}

func (payload *Payload) DumpGlobalsMemoryStack(w io.Writer) (err error) {
	fd := int(payload.maps.Fd())

	dataMapOffset := int64(payload.info.RODataSize) + int64(payload.info.TextSize)

	globalsMemorySize := payload.info.MemoryOffset + payload.info.GrowMemorySize
	dataSize := int(globalsMemorySize) + int(payload.info.StackSize)

	data, err := syscall.Mmap(fd, dataMapOffset, dataSize, syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)

	buf := data[:payload.info.MemoryOffset]
	fmt.Fprintf(w, "--- GLOBALS (%d kB) ---\n", len(buf)/1024)
	for i := 0; len(buf) > 0; i += 8 {
		fmt.Fprintf(w, "%08x: %x\n", i, buf[0:8])
		buf = buf[8:]
	}

	buf = data[payload.info.MemoryOffset : payload.info.MemoryOffset+globalsMemorySize]
	fmt.Fprintf(w, "--- MEMORY (%d kB) ---\n", len(buf)/1024)
	for i := 0; len(buf) > 0; i += 32 {
		fmt.Fprintf(w, "%08x: %x %x %x %x\n", i, buf[0:8], buf[8:16], buf[16:24], buf[24:32])
		buf = buf[32:]
	}

	buf = data[globalsMemorySize:]
	fmt.Fprintf(w, "--- STACK (%d kB) ---\n", len(buf)/1024)
	for i := 0; len(buf) > 0; i += 32 {
		fmt.Fprintf(w, "%08x: %x %x %x %x\n", i, buf[0:8], buf[8:16], buf[16:24], buf[24:32])
		buf = buf[32:]
	}

	fmt.Fprintf(w, "---\n")
	return
}

func Run(env *Environment, payload *Payload) (output []byte, err error) {
	cmd := exec.Cmd{
		Path: env.executor,
		Args: []string{},
		Env:  []string{},
		Dir:  "/",
		ExtraFiles: []*os.File{
			payload.maps,
		},
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig:  syscall.SIGKILL,
			Cloneflags: syscall.CLONE_NEWUSER,
		},
	}

	cmd.Stderr = env.loader

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return
	}

	err = cmd.Start()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return
	}

	err = binary.Write(stdin, nativeEndian, payload.info)
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return
	}

	output, _ = ioutil.ReadAll(stdout)

	err = cmd.Wait()
	if err != nil {
		return
	}

	if !cmd.ProcessState.Success() {
		err = errors.New(cmd.ProcessState.String())
	}
	return
}
