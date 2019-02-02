// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"bufio"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"os"

	"github.com/tsavola/gate/entry"
	"github.com/tsavola/gate/image"
	"github.com/tsavola/gate/image/metadata"
	"github.com/tsavola/gate/internal/error/resourcelimit"
	"github.com/tsavola/gate/runtime"
	"github.com/tsavola/gate/runtime/abi"
	"github.com/tsavola/gate/server/event"
	"github.com/tsavola/gate/server/internal/error/failrequest"
	"github.com/tsavola/wag/compile"
	"github.com/tsavola/wag/object"
	"github.com/tsavola/wag/object/stack"
	"github.com/tsavola/wag/section"
	"github.com/tsavola/wag/wa"
)

var pageSize = os.Getpagesize()

var (
	newHash      = sha512.New384
	hashEncoding = base64.URLEncoding
)

func validateHashBytes(hash1 string, digest2 []byte) (err error) {
	digest1, err := hashEncoding.DecodeString(hash1)
	if err != nil {
		return
	}

	if subtle.ConstantTimeCompare(digest1, digest2) != 1 {
		err = failrequest.New(event.FailRequest_ModuleHashMismatch, "module hash does not match content")
		return
	}

	return
}

func validateHashContent(hash1 string, r io.Reader) (err error) {
	hash2 := newHash()

	_, err = io.Copy(hash2, r)
	if err != nil {
		err = wrapContentError(err)
		return
	}

	return validateHashBytes(hash1, hash2.Sum(nil))
}

type binary struct {
	codeMap object.CallMap
	archive image.Archive
	*image.ArchiveManifest

	module image.Module // TODO: separate reference counting?

	// Protected by Server.lock:
	refCount int
}

// ref must be called with Server.lock held.
func (bin *binary) ref() *binary {
	bin.refCount++
	return bin
}

// unref must be called with Server.lock held.
func (bin *binary) unref() {
	bin.refCount--
	if bin.refCount == 0 {
		bin.archive.Close()
		bin.module.Close()
	}
}

type program struct {
	hash    string
	bin     *binary
	routine runtime.InitRoutine

	module image.Module // Non-nil overrides bin.module.

	// Protected by Server.lock:
	refCount int
}

// compileProgram returns Executable if ExecutableRef is provided.
// InstancePolicy must be provided with ExecutableRef.  Entry name can be
// provided only with ExecutableRef.
func compileProgram(ctx context.Context, ref image.ExecutableRef, instPolicy *InstancePolicy, progPolicy *ProgramPolicy, storage image.Storage, allegedHash string, content io.ReadCloser, contentSize int, entryName string,
) (exe *image.Executable, prog *program, err error) {
	defer func() {
		if content != nil {
			content.Close()
		}
	}()

	moduleStore, err := storage.CreateModule(ctx, contentSize)
	if err != nil {
		return
	}
	defer moduleStore.Close()

	var actualHash = newHash()
	var r = bufio.NewReader(io.TeeReader(io.TeeReader(content, moduleStore.Writer), actualHash))

	var sectionMap = section.NewMap()
	var sectionLoaders = make(section.CustomLoaders)
	var sectionConfig = compile.Config{
		SectionMapper:       sectionMap.Mapper(),
		CustomSectionLoader: sectionLoaders.Load,
	}

	sectionLoaders["gate.stack"] = func(_ string, _ section.Reader, _ uint32) error {
		return errors.New("gate.stack section appears too early in wasm module")
	}

	module, err := compile.LoadInitialSections(&compile.ModuleConfig{Config: sectionConfig}, r)
	if err != nil {
		err = failrequest.Tag(event.FailRequest_ModuleError, err)
		return
	}

	var stackSize int
	var memorySize = module.InitialMemorySize()
	var maxMemorySize int

	if instPolicy == nil {
		// Building to storage.
		maxMemorySize = memorySize
	} else {
		stackSize = instPolicy.StackSize

		maxMemorySize = module.MemorySizeLimit()
		if maxMemorySize > instPolicy.MaxMemorySize {
			maxMemorySize = roundSize(instPolicy.MaxMemorySize, wa.PageSize)
		}

		if memorySize > maxMemorySize {
			err = resourcelimit.New("initial program memory size exceeds instance memory size limit")
			return
		}
	}

	err = abi.BindImports(module)
	if err != nil {
		return
	}

	var entryIndex uint32

	if entryName != "" {
		entryIndex, err = entry.FuncIndex(module, entryName)
		if err != nil {
			return
		}
	}

	if ref == nil {
		back, ok := storage.(image.BackingStore)
		if !ok {
			back = image.Memory
		}

		ref, err = image.NewExecutableRef(back)
		if err != nil {
			return
		}
		defer ref.Close()
	}

	build := image.NewBuild(ref)
	defer build.Close()

	buildConfig := new(image.BuildConfig)
	buildConfig.MaxTextSize = progPolicy.MaxTextSize

	err = build.Configure(buildConfig)
	if err != nil {
		return
	}

	var routine = runtime.InitStart
	var codeMap = new(object.CallMap)
	var codeConfig = &compile.CodeConfig{
		Text:   build.TextBuffer(),
		Mapper: codeMap,
		Config: sectionConfig,
	}

	err = compile.LoadCodeSection(codeConfig, r, module)
	if err != nil {
		err = failrequest.Tag(event.FailRequest_ModuleError, err)
		return
	}

	buildConfig.MaxTextSize = build.TextSize()
	buildConfig.StackSize = stackSize
	buildConfig.GlobalsSize = module.GlobalsSize()
	buildConfig.MemorySize = memorySize
	buildConfig.MaxMemorySize = maxMemorySize

	err = build.Configure(buildConfig)
	if err != nil {
		return
	}

	if stackSize != 0 {
		var addr uint32
		if entryName != "" {
			addr = codeMap.FuncAddrs[entryIndex]
		}
		build.SetupEntryStackFrame(addr)
	}

	sectionLoaders["gate.stack"] = func(_ string, r section.Reader, payloadLen uint32) error {
		if entryName != "" {
			return errors.New("entry function specified with suspended program")
		}

		routine = runtime.InitResume

		if _, err := r.Read(build.StackBytes()); err != nil {
			return err
		}

		switch _, err := r.ReadByte(); err {
		case io.EOF:
			return nil

		case nil:
			return resourcelimit.New("suspended program stack size exceeds call stack size limit")

		default:
			return err
		}
	}

	var dataConfig = &compile.DataConfig{
		GlobalsMemory:   build.GlobalsMemoryBuffer(),
		MemoryAlignment: pageSize,
		Config:          sectionConfig,
	}

	err = compile.LoadDataSection(dataConfig, r, module)
	if err != nil {
		err = failrequest.Tag(event.FailRequest_ModuleError, err)
		return
	}

	_, err = io.Copy(ioutil.Discard, r)
	if err != nil {
		err = wrapContentError(err)
		return
	}

	err = content.Close()
	content = nil
	if err != nil {
		err = wrapContentError(err)
		return
	}

	var actualDigest = actualHash.Sum(nil)

	if allegedHash != "" {
		err = validateHashBytes(allegedHash, actualDigest)
		if err != nil {
			return
		}
	}

	var hash = hashEncoding.EncodeToString(actualDigest)

	storedModule, err := moduleStore.Module(hash)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			storedModule.Close()
		}
	}()

	exe, err = build.Executable()
	if err != nil {
		return
	}
	defer func() {
		if err != nil && exe != nil {
			exe.Close()
			exe = nil
		}
	}()

	var meta = metadata.New(module, sectionMap, codeMap)
	var archive image.Archive

	if instPolicy == nil {
		archive, err = exe.StoreThis(ctx, hash, meta, storage)
		exe.Close()
		exe = nil
	} else {
		archive, err = exe.StoreCopy(ctx, hash, meta, storage)
	}
	if err != nil {
		return
	}

	prog = &program{
		hash: hash,
		bin: &binary{
			codeMap:         *codeMap,
			archive:         archive,
			ArchiveManifest: archive.Manifest(),
			module:          storedModule,
			refCount:        1,
		},
		routine:  routine,
		refCount: 1,
	}
	return
}

// ref must be called with Server.lock held.
func (prog *program) ref() *program {
	prog.refCount++
	return prog
}

// unref must be called with Server.lock held.
func (prog *program) unref() (final bool) {
	prog.refCount--
	if prog.refCount == 0 {
		prog.bin.unref()
		if prog.module != nil {
			prog.module.Close()
		}
		final = true
	}
	return
}

func (prog *program) getEntryAddr(name string) (addr uint32, err error) {
	if name != "" {
		addr, err = entry.FuncAddr(prog.bin.EntryAddrs, name)
	}
	return
}

func (prog *program) loadExecutable(ctx context.Context, ref image.ExecutableRef, instPolicy *InstancePolicy, entryAddr uint32,
) (exe *image.Executable, err error) {
	config := &image.Config{
		MaxTextSize:   math.MaxInt32, // Policy was enforced when program was acquired.
		StackSize:     instPolicy.StackSize,
		MaxMemorySize: roundSize(instPolicy.MaxMemorySize, wa.PageSize),
	}

	return image.LoadExecutable(ctx, ref, config, prog.bin.archive, stack.EntryFrame(entryAddr, nil))
}

func roundSize(n, pageSize int) int {
	mask := pageSize - 1
	return (n + mask) &^ mask
}
