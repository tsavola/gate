// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"io"
	"log"
	goruntime "runtime"

	"gate.computer/gate/image"
	"gate.computer/gate/internal/error/public"
	"gate.computer/gate/internal/error/resourcelimit"
	"gate.computer/gate/internal/principal"
	"gate.computer/gate/runtime"
	"gate.computer/gate/scope"
	"gate.computer/gate/server/api"
	"gate.computer/gate/server/detail"
	"gate.computer/gate/server/event"
	"gate.computer/gate/server/internal/error/resourcenotfound"
	"gate.computer/wag/object/stack"
)

const ErrServerClosed = public.Err("server closed")

var errAnonymous = AccessUnauthorized("anonymous access not supported")

type progPolicy struct {
	res  ResourcePolicy
	prog ProgramPolicy
}

type instPolicy struct {
	res  ResourcePolicy
	inst InstancePolicy
}

type instProgPolicy struct {
	res  ResourcePolicy
	prog ProgramPolicy
	inst InstancePolicy
}

type privateConfig struct {
	Config
}

type Server struct {
	privateConfig

	mu        serverMutex
	programs  map[string]*program
	accounts  map[principal.RawKey]*account
	anonymous map[*Instance]struct{}
}

func New(ctx context.Context, config *Config) (_ *Server, err error) {
	defer func() { err = asError(recover()) }()

	s := &Server{
		programs:  make(map[string]*program),
		accounts:  make(map[principal.RawKey]*account),
		anonymous: make(map[*Instance]struct{}),
	}

	if config != nil {
		s.Config = *config
	}
	if s.ImageStorage == nil {
		s.ImageStorage = image.Memory
	}
	if s.Monitor == nil {
		s.Monitor = defaultMonitor
	}
	if !s.Configured() {
		panic("incomplete server configuration")
	}

	progs, err := s.ImageStorage.Programs()
	_check(err)

	insts, err := s.ImageStorage.Instances()
	_check(err)

	shutdown := s.Shutdown
	defer func() {
		if shutdown != nil {
			shutdown(context.Background())
		}
	}()

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	var owner *account
	if id := principal.ContextID(ctx); id != nil {
		owner = newAccount(id)
		s.accounts[principal.Raw(id)] = owner
	}

	for _, id := range progs {
		s._loadProgramDuringInit(lock, owner, id)
	}

	for _, key := range insts {
		s._loadInstanceDuringInit(lock, key)
	}

	shutdown = nil
	return s, nil
}

func (s *Server) _loadProgramDuringInit(lock serverLock, owner *account, progID string) {
	image, err := s.ImageStorage.LoadProgram(progID)
	_check(err)
	if image == nil { // Race condition with human operator?
		return
	}
	defer closeProgramImage(&image)

	buffers, err := image.LoadBuffers()
	_check(err)

	prog := newProgram(progID, image, buffers, true)
	image = nil

	if owner != nil {
		owner.ensureProgramRef(lock, prog, nil)
	}

	s.programs[progID] = prog
}

func (s *Server) _loadInstanceDuringInit(lock serverLock, key string) {
	image, err := s.ImageStorage.LoadInstance(key)
	_check(err)
	if image == nil { // Race condition with human operator?
		return
	}
	defer closeInstanceImage(&image)

	pri, instID := _parseInstanceStorageKey(key)
	acc := s.ensureAccount(lock, pri)

	// TODO: restore instance
	log.Printf("TODO: load account %s instance %s (%s)", acc.ID, instID, image.Trap())
}

func (s *Server) Shutdown(ctx context.Context) error {
	var (
		accInsts  []*Instance
		anonInsts map[*Instance]struct{}
	)
	s.mu.Guard(func(lock serverLock) {
		progs := s.programs
		s.programs = nil

		for _, prog := range progs {
			prog.unref(lock)
		}

		accs := s.accounts
		s.accounts = nil

		for _, acc := range accs {
			for _, x := range acc.shutdown(lock) {
				accInsts = append(accInsts, x.inst)
				x.prog.unref(lock)
			}
		}

		anonInsts = s.anonymous
		s.anonymous = nil
	})

	for _, inst := range accInsts {
		inst.suspend()
	}
	for inst := range anonInsts {
		inst.Kill()
	}

	var aborted bool

	for _, inst := range accInsts {
		if inst.Wait(ctx).State == api.StateRunning {
			aborted = true
		}
	}
	for inst := range anonInsts {
		inst.Wait(ctx)
	}

	if aborted {
		return ctx.Err()
	}
	return nil
}

func (*Server) Features(ctx context.Context) (*api.Features, error) {
	return &api.Features{Scope: scope.Names()}, nil
}

func (s *Server) UploadModule(ctx context.Context, upload *ModuleUpload, know *api.ModuleOptions) (module string, err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)

	policy := new(progPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeProgram(ctx, &policy.res, &policy.prog))

	if upload.Length > int64(policy.prog.MaxModuleSize) {
		panic(resourcelimit.New("module size limit exceeded"))
	}

	// TODO: check resource policy

	if upload.Hash != "" && s._loadKnownModule(ctx, policy, upload, know) {
		return upload.Hash, nil
	}

	return s._loadUnknownModule(ctx, policy, upload, know), nil
}

func (s *Server) SourceModule(ctx context.Context, source *ModuleSource, know *api.ModuleOptions) (module string, err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)

	policy := new(progPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeProgramSource(ctx, &policy.res, &policy.prog, source.Source))

	stream, length, err := source.Source.OpenURI(ctx, source.URI, policy.prog.MaxModuleSize)
	_check(err)
	if stream == nil {
		if length > 0 {
			panic(resourcelimit.New("program size limit exceeded"))
		}
		panic(resourcenotfound.ErrModule)
	}

	upload := &ModuleUpload{
		Stream: stream,
		Length: length,
	}
	defer upload.Close()

	return s._loadUnknownModule(ctx, policy, upload, know), nil
}

func (s *Server) _loadKnownModule(ctx context.Context, policy *progPolicy, upload *ModuleUpload, know *api.ModuleOptions) bool {
	prog := s._refProgram(upload.Hash, upload.Length)
	if prog == nil {
		return false
	}
	defer s.unrefProgram(&prog)
	progID := prog.id

	upload._validate()

	if prog.image.TextSize() > policy.prog.MaxTextSize {
		panic(resourcelimit.New("program code size limit exceeded"))
	}

	s._registerProgramRef(ctx, prog, know)
	prog = nil

	s.monitor(&event.ModuleUploadExist{
		Ctx:    ContextDetail(ctx),
		Module: progID,
	})

	return true
}

func (s *Server) _loadUnknownModule(ctx context.Context, policy *progPolicy, upload *ModuleUpload, know *api.ModuleOptions) string {
	prog, _ := _buildProgram(s.ImageStorage, &policy.prog, nil, upload, "")
	defer s.unrefProgram(&prog)
	progID := prog.id

	redundant := s._registerProgramRef(ctx, prog, know)
	prog = nil

	if redundant {
		s.monitor(&event.ModuleUploadExist{
			Ctx:      ContextDetail(ctx),
			Module:   progID,
			Compiled: true,
		})
	} else {
		s.monitor(&event.ModuleUploadNew{
			Ctx:    ContextDetail(ctx),
			Module: progID,
		})
	}

	return progID
}

func (s *Server) NewInstance(ctx context.Context, module string, launch *api.LaunchOptions, invoke *InvokeOptions) (_ *Instance, err error) {
	defer func() { err = asError(recover()) }()
	launch = _prepareLaunchOptions(launch)

	policy := new(instPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeInstance(ctx, &policy.res, &policy.inst))

	acc := s._checkAccountInstanceID(ctx, launch.Instance)
	if acc == nil {
		panic(errAnonymous)
	}

	prog := s.mu.GuardProgram(func(lock serverLock) *program {
		prog := s.programs[module]
		if prog == nil {
			return nil
		}

		return acc.refProgram(lock, prog)
	})
	if prog == nil {
		panic(resourcenotfound.ErrModule)
	}
	defer s.unrefProgram(&prog)

	funcIndex, err := prog.image.ResolveEntryFunc(launch.Function, false)
	_check(err)

	// TODO: check resource policy (text/stack/memory/max-memory size etc.)

	instImage, err := image.NewInstance(prog.image, policy.inst.MaxMemorySize, policy.inst.StackSize, funcIndex)
	_check(err)
	defer closeInstanceImage(&instImage)

	ref := &api.ModuleOptions{}
	inst, prog, _ := s._registerProgramRefInstance(ctx, acc, prog, instImage, &policy.inst, ref, launch, invoke)
	instImage = nil

	s._runOrDeleteInstance(ctx, inst, prog, launch.Function)
	prog = nil

	s.monitor(&event.InstanceCreateKnown{
		Ctx:    ContextDetail(ctx),
		Create: newInstanceCreateEvent(inst.ID, module, launch),
	})

	return inst, nil
}

func (s *Server) UploadModuleInstance(ctx context.Context, upload *ModuleUpload, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) (_ *Instance, err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)
	launch = _prepareLaunchOptions(launch)

	policy := new(instProgPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeProgramInstance(ctx, &policy.res, &policy.prog, &policy.inst))

	acc := s._checkAccountInstanceID(ctx, launch.Instance)
	_, inst := s._loadModuleInstance(ctx, acc, policy, upload, know, launch, invoke)
	return inst, nil
}

func (s *Server) SourceModuleInstance(ctx context.Context, source *ModuleSource, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) (module string, _ *Instance, err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)
	launch = _prepareLaunchOptions(launch)

	policy := new(instProgPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeProgramInstanceSource(ctx, &policy.res, &policy.prog, &policy.inst, source.Source))

	acc := s._checkAccountInstanceID(ctx, launch.Instance)

	stream, length, err := source.Source.OpenURI(ctx, source.URI, policy.prog.MaxModuleSize)
	_check(err)
	if stream == nil {
		if length > 0 {
			panic(resourcelimit.New("program size limit exceeded"))
		}
		panic(resourcenotfound.ErrModule)
	}

	upload := &ModuleUpload{
		Stream: stream,
		Length: length,
	}
	defer upload.Close()

	module, inst := s._loadModuleInstance(ctx, acc, policy, upload, know, launch, invoke)
	return module, inst, nil
}

func (s *Server) _loadModuleInstance(ctx context.Context, acc *account, policy *instProgPolicy, upload *ModuleUpload, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) (string, *Instance) {
	if upload.Length > int64(policy.prog.MaxModuleSize) {
		panic(resourcelimit.New("module size limit exceeded"))
	}

	// TODO: check resource policy

	if upload.Hash != "" {
		inst := s._loadKnownModuleInstance(ctx, acc, policy, upload, know, launch, invoke)
		if inst != nil {
			return upload.Hash, inst
		}
	}

	return s._loadUnknownModuleInstance(ctx, acc, policy, upload, know, launch, invoke)
}

func (s *Server) _loadKnownModuleInstance(ctx context.Context, acc *account, policy *instProgPolicy, upload *ModuleUpload, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) *Instance {
	prog := s._refProgram(upload.Hash, upload.Length)
	if prog == nil {
		return nil
	}
	defer s.unrefProgram(&prog)
	progID := prog.id

	upload._validate()

	if prog.image.TextSize() > policy.prog.MaxTextSize {
		panic(resourcelimit.New("program code size limit exceeded"))
	}

	// TODO: check resource policy (stack/memory/max-memory size etc.)

	funcIndex, err := prog.image.ResolveEntryFunc(launch.Function, false)
	_check(err)

	instImage, err := image.NewInstance(prog.image, policy.inst.MaxMemorySize, policy.inst.StackSize, funcIndex)
	_check(err)
	defer closeInstanceImage(&instImage)

	inst, prog, _ := s._registerProgramRefInstance(ctx, acc, prog, instImage, &policy.inst, know, launch, invoke)
	instImage = nil

	s.monitor(&event.ModuleUploadExist{
		Ctx:    ContextDetail(ctx),
		Module: progID,
	})

	s._runOrDeleteInstance(ctx, inst, prog, launch.Function)
	prog = nil

	s.monitor(&event.InstanceCreateKnown{
		Ctx:    ContextDetail(ctx),
		Create: newInstanceCreateEvent(inst.ID, progID, launch),
	})

	return inst
}

func (s *Server) _loadUnknownModuleInstance(ctx context.Context, acc *account, policy *instProgPolicy, upload *ModuleUpload, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) (string, *Instance) {
	prog, instImage := _buildProgram(s.ImageStorage, &policy.prog, &policy.inst, upload, launch.Function)
	defer closeInstanceImage(&instImage)
	defer s.unrefProgram(&prog)
	progID := prog.id

	inst, prog, redundantProg := s._registerProgramRefInstance(ctx, acc, prog, instImage, &policy.inst, know, launch, invoke)
	instImage = nil

	if upload.Hash != "" {
		if redundantProg {
			s.monitor(&event.ModuleUploadExist{
				Ctx:      ContextDetail(ctx),
				Module:   progID,
				Compiled: true,
			})
		} else {
			s.monitor(&event.ModuleUploadNew{
				Ctx:    ContextDetail(ctx),
				Module: progID,
			})
		}
	} else {
		if redundantProg {
			s.monitor(&event.ModuleSourceExist{
				Ctx:    ContextDetail(ctx),
				Module: progID,
				// TODO: source URI
				Compiled: true,
			})
		} else {
			s.monitor(&event.ModuleSourceNew{
				Ctx:    ContextDetail(ctx),
				Module: progID,
				// TODO: source URI
			})
		}
	}

	s._runOrDeleteInstance(ctx, inst, prog, launch.Function)
	prog = nil

	s.monitor(&event.InstanceCreateStream{
		Ctx:    ContextDetail(ctx),
		Create: newInstanceCreateEvent(inst.ID, progID, launch),
	})

	return progID, inst
}

func (s *Server) ModuleInfo(ctx context.Context, module string) (_ *api.ModuleInfo, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.programs == nil {
		panic(ErrServerClosed)
	}
	prog := s.programs[module]
	if prog == nil {
		panic(resourcenotfound.ErrModule)
	}

	acc := s.accounts[principal.Raw(pri)]
	if acc == nil {
		panic(resourcenotfound.ErrModule)
	}

	x, found := acc.programs[prog]
	if !found {
		panic(resourcenotfound.ErrModule)
	}

	info := &api.ModuleInfo{
		Id:   prog.id,
		Tags: append([]string(nil), x.tags...),
	}

	s.monitor(&event.ModuleInfo{
		Ctx:    ContextDetail(ctx),
		Module: prog.id,
	})

	return info, nil
}

func (s *Server) Modules(ctx context.Context) (_ *api.Modules, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	s.monitor(&event.ModuleList{
		Ctx: ContextDetail(ctx),
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	acc := s.accounts[principal.Raw(pri)]
	if acc == nil {
		return new(api.Modules), nil
	}

	infos := &api.Modules{
		Modules: make([]*api.ModuleInfo, 0, len(acc.programs)),
	}
	for prog, x := range acc.programs {
		infos.Modules = append(infos.Modules, &api.ModuleInfo{
			Id:   prog.id,
			Tags: append([]string(nil), x.tags...),
		})
	}
	return infos, nil
}

func (s *Server) ModuleContent(ctx context.Context, module string) (stream io.ReadCloser, length int64, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	prog := s.mu.GuardProgram(func(lock serverLock) *program {
		acc := s.accounts[principal.Raw(pri)]
		if acc == nil {
			return nil
		}

		prog := s.programs[module]
		if prog == nil {
			return nil
		}

		return acc.refProgram(lock, prog)
	})
	if prog == nil {
		panic(resourcenotfound.ErrModule)
	}

	length = prog.image.ModuleSize()
	stream = &moduleContent{
		ctx:   ContextDetail(ctx),
		r:     prog.image.NewModuleReader(),
		s:     s,
		prog:  prog,
		total: length,
	}
	return stream, length, nil
}

type moduleContent struct {
	ctx   *detail.Context
	r     io.Reader
	s     *Server
	prog  *program
	total int64
	read  int64
}

func (x *moduleContent) Read(b []byte) (int, error) {
	n, err := x.r.Read(b)
	x.read += int64(n)
	return n, err
}

func (x *moduleContent) Close() error {
	x.s.monitor(&event.ModuleDownload{
		Ctx:          x.ctx,
		Module:       x.prog.id,
		ModuleLength: uint64(x.total),
		LengthRead:   uint64(x.read),
	})

	x.s.unrefProgram(&x.prog)
	return nil
}

func (s *Server) PinModule(ctx context.Context, module string, know *api.ModuleOptions) (err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)
	if !know.GetPin() {
		panic("Server.PinModule called without ModuleOptions.Pin")
	}

	policy := new(progPolicy)
	ctx = _context(s.AccessPolicy.AuthorizeProgram(ctx, &policy.res, &policy.prog))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	modified := s.mu.GuardBool(func(lock serverLock) bool {
		if s.programs == nil {
			panic(ErrServerClosed)
		}
		prog := s.programs[module]
		if prog == nil {
			panic(resourcenotfound.ErrModule)
		}

		acc := s.accounts[principal.Raw(pri)]
		if acc == nil {
			panic(resourcenotfound.ErrModule)
		}

		if _, found := acc.programs[prog]; !found {
			for _, x := range acc.instances {
				if x.prog == prog {
					goto do
				}
			}
			panic(resourcenotfound.ErrModule)
		}

	do:
		// TODO: check resource limits
		return acc.ensureProgramRef(lock, prog, know.Tags)
	})

	if modified {
		s.monitor(&event.ModulePin{
			Ctx:      ContextDetail(ctx),
			Module:   module,
			TagCount: int32(len(know.Tags)),
		})
	}

	return
}

func (s *Server) UnpinModule(ctx context.Context, module string) (err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	found := s.mu.GuardBool(func(lock serverLock) bool {
		acc := s.accounts[principal.Raw(pri)]
		if acc == nil {
			return false
		}

		prog := s.programs[module]
		if prog == nil {
			return false
		}

		return acc.unrefProgram(lock, prog)
	})
	if !found {
		panic(resourcenotfound.ErrModule)
	}

	s.monitor(&event.ModuleUnpin{
		Ctx:    ContextDetail(ctx),
		Module: module,
	})

	return
}

type IOFunc func(context.Context, io.Reader, io.Writer) error

func (s *Server) InstanceConnection(ctx context.Context, instance string) (_ *Instance, _ IOFunc, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	inst := s._getInstance(ctx, instance)
	conn := inst.connect(ctx)
	if conn == nil {
		s.monitor(&event.FailRequest{
			Ctx:      ContextDetail(ctx),
			Failure:  event.FailInstanceNoConnect,
			Instance: inst.ID,
		})
		return inst, nil, nil
	}

	iofunc := func(ctx context.Context, r io.Reader, w io.Writer) error {
		s.monitor(&event.InstanceConnect{
			Ctx:      ContextDetail(ctx),
			Instance: inst.ID,
		})

		err := conn(ctx, r, w)

		s.Monitor(&event.InstanceDisconnect{
			Ctx:      ContextDetail(ctx),
			Instance: inst.ID,
		}, err)

		return err
	}

	return inst, iofunc, nil
}

func (s *Server) InstanceInfo(ctx context.Context, instance string) (_ *api.InstanceInfo, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	progID, inst := s._getInstanceProgramID(ctx, instance)
	info := inst.info(progID)
	if info == nil {
		panic(resourcenotfound.ErrInstance)
	}

	s.monitor(&event.InstanceInfo{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
	})

	return info, nil
}

func (s *Server) WaitInstance(ctx context.Context, instID string) (_ *api.Status, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	inst := s._getInstance(ctx, instID)
	status := inst.Wait(ctx)

	s.monitor(&event.InstanceWait{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
	})

	return status, err
}

func (s *Server) KillInstance(ctx context.Context, instance string) (_ *Instance, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	inst := s._getInstance(ctx, instance)
	inst.Kill()

	s.monitor(&event.InstanceKill{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
	})

	return inst, nil
}

func (s *Server) SuspendInstance(ctx context.Context, instance string) (_ *Instance, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	// Store the program in case the instance becomes non-transient.
	inst, prog := s._getInstanceRefProgram(ctx, instance)
	defer s.unrefProgram(&prog)

	prog._ensureStorage()
	inst.Suspend()

	s.monitor(&event.InstanceSuspend{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
	})

	return inst, nil
}

func (s *Server) ResumeInstance(ctx context.Context, instance string, resume *api.ResumeOptions, invoke *InvokeOptions) (_ *Instance, err error) {
	defer func() { err = asError(recover()) }()

	resume = prepareResumeOptions(resume)
	policy := new(instPolicy)

	ctx = _context(s.AccessPolicy.AuthorizeInstance(ctx, &policy.res, &policy.inst))

	inst, prog := s._getInstanceRefProgram(ctx, instance)
	defer s.unrefProgram(&prog)

	inst._checkResume(resume.Function)

	proc, services := s._allocateInstanceResources(ctx, &policy.inst)
	defer closeInstanceResources(&proc, &services)

	inst._doResume(resume.Function, proc, services, policy.inst.TimeResolution, invoke)
	proc = nil
	services = nil

	s._runOrDeleteInstance(ctx, inst, prog, resume.Function)
	prog = nil

	s.monitor(&event.InstanceResume{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
		Function: resume.Function,
	})

	return inst, nil
}

func (s *Server) DeleteInstance(ctx context.Context, instance string) (err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	inst := s._getInstance(ctx, instance)
	inst._annihilate()
	s.deleteNonexistentInstance(inst)

	s.monitor(&event.InstanceDelete{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
	})

	return
}

func (s *Server) Snapshot(ctx context.Context, instance string, know *api.ModuleOptions) (module string, err error) {
	defer func() { err = asError(recover()) }()
	know = _prepareModuleOptions(know)
	if !know.GetPin() {
		panic("Server.SnapshotInstance called without ModuleOptions.Pin")
	}

	inst := s._getInstance(ctx, instance)

	// TODO: implement suspend-snapshot-resume at a lower level

	if inst.Status().State == api.StateRunning {
		inst.suspend()
		if inst.Wait(context.Background()).State == api.StateSuspended {
			defer func() {
				_, e := s.ResumeInstance(ctx, instance, nil, nil)
				if module != "" {
					_check(e)
				}
			}()
		}
	}

	module = s._snapshot(ctx, instance, know)
	return
}

func (s *Server) _snapshot(ctx context.Context, instance string, know *api.ModuleOptions) string {
	ctx = _context(s.AccessPolicy.Authorize(ctx))

	// TODO: check module storage limits

	inst, oldProg := s._getInstanceRefProgram(ctx, instance)
	defer s.unrefProgram(&oldProg)

	newImage, buffers := inst._snapshot(oldProg)
	defer closeProgramImage(&newImage)

	h := api.KnownModuleHash.New()
	_, err := io.Copy(h, newImage.NewModuleReader())
	_check(err)
	progID := api.EncodeKnownModule(h.Sum(nil))

	_check(newImage.Store(progID))

	newProg := newProgram(progID, newImage, buffers, true)
	newImage = nil
	defer s.unrefProgram(&newProg)

	s._registerProgramRef(ctx, newProg, know)
	newProg = nil

	s.monitor(&event.InstanceSnapshot{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
		Module:   progID,
	})

	return progID
}

func (s *Server) UpdateInstance(ctx context.Context, instance string, update *api.InstanceUpdate) (_ *api.InstanceInfo, err error) {
	defer func() { err = asError(recover()) }()
	update = prepareInstanceUpdate(update)

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	progID, inst := s._getInstanceProgramID(ctx, instance)
	if inst.update(update) {
		s.monitor(&event.InstanceUpdate{
			Ctx:      ContextDetail(ctx),
			Instance: inst.ID,
			Persist:  update.Persist,
			TagCount: int32(len(update.Tags)),
		})
	}

	info := inst.info(progID)
	if info == nil {
		panic(resourcenotfound.ErrInstance)
	}

	return info, nil
}

func (s *Server) DebugInstance(ctx context.Context, instance string, req *api.DebugRequest) (_ *api.DebugResponse, err error) {
	defer func() { err = asError(recover()) }()
	policy := new(progPolicy)

	ctx = _context(s.AccessPolicy.AuthorizeProgram(ctx, &policy.res, &policy.prog))

	inst, defaultProg := s._getInstanceRefProgram(ctx, instance)
	defer s.unrefProgram(&defaultProg)

	rebuild, config, res := inst._debug(ctx, defaultProg, req)
	if rebuild != nil {
		var (
			progImage *image.Program
			textMap   stack.TextMap
			ok        bool
		)

		progImage, textMap = _rebuildProgramImage(s.ImageStorage, &policy.prog, defaultProg.image.NewModuleReader(), config.DebugInfo, config.Breakpoints)
		defer func() {
			if progImage != nil {
				progImage.Close()
			}
		}()

		res, ok = rebuild.apply(progImage, config, textMap)
		if !ok {
			panic(public.Err("conflict")) // TODO: http response code: conflict
		}
		progImage = nil
	}

	s.monitor(&event.InstanceDebug{
		Ctx:      ContextDetail(ctx),
		Instance: inst.ID,
		Compiled: rebuild != nil,
	})

	return res, nil
}

func (s *Server) Instances(ctx context.Context) (_ *api.Instances, err error) {
	defer func() { err = asError(recover()) }()

	ctx = _context(s.AccessPolicy.Authorize(ctx))

	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	s.monitor(&event.InstanceList{
		Ctx: ContextDetail(ctx),
	})

	type instProgID struct {
		inst   *Instance
		progID string
	}

	// Get instance references while holding server lock.
	var insts []instProgID
	s.mu.Guard(func(lock serverLock) {
		if acc := s.accounts[principal.Raw(pri)]; acc != nil {
			insts = make([]instProgID, 0, len(acc.instances))
			for _, x := range acc.instances {
				insts = append(insts, instProgID{x.inst, x.prog.id})
			}
		}
	})

	// Each instance has its own lock.
	infos := &api.Instances{
		Instances: make([]*api.InstanceInfo, 0, len(insts)),
	}
	for _, x := range insts {
		if info := x.inst.info(x.progID); info != nil {
			infos.Instances = append(infos.Instances, info)
		}
	}
	return infos, nil
}

// ensureAccount may return nil.  It must not be called while the server is
// shutting down.
func (s *Server) ensureAccount(_ serverLock, pri *principal.ID) *account {
	acc := s.accounts[principal.Raw(pri)]
	if acc == nil {
		acc = newAccount(pri)
		s.accounts[principal.Raw(pri)] = acc
	}
	return acc
}

func (s *Server) _refProgram(hash string, length int64) *program {
	lock := s.mu.Lock()
	defer s.mu.Unlock()

	prog := s.programs[hash]
	if prog == nil {
		return nil
	}

	if length != prog.image.ModuleSize() {
		panic(errModuleSizeMismatch)
	}

	return prog.ref(lock)
}

func (s *Server) unrefProgram(p **program) {
	prog := *p
	*p = nil
	if prog == nil {
		return
	}

	s.mu.Guard(prog.unref)
}

// registerProgramRef with the server and an account.  Caller's program
// reference is stolen (except on error).
func (s *Server) _registerProgramRef(ctx context.Context, prog *program, know *api.ModuleOptions) (redundant bool) {
	var pri *principal.ID

	if know.Pin {
		pri = principal.ContextID(ctx)
		if pri == nil {
			panic(errAnonymous)
		}

		prog._ensureStorage()
	}

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	prog, redundant = s._mergeProgramRef(lock, prog)

	if know.Pin {
		// mergeProgramRef checked for shutdown, so the ensure methods are safe
		// to call.
		if s.ensureAccount(lock, pri).ensureProgramRef(lock, prog, know.Tags) {
			// TODO: move outside of critical section
			s.monitor(&event.ModulePin{
				Ctx:      ContextDetail(ctx),
				Module:   prog.id,
				TagCount: int32(len(know.Tags)),
			})
		}
	}

	return
}

func (s *Server) _checkAccountInstanceID(ctx context.Context, instID string) *account {
	if instID != "" {
		_validateInstanceID(instID)
	}

	pri := principal.ContextID(ctx)
	if pri == nil {
		return nil
	}

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	if s.accounts == nil {
		panic(ErrServerClosed)
	}

	acc := s.ensureAccount(lock, pri)

	if instID != "" {
		acc._checkUniqueInstanceID(lock, instID)
	}

	return acc
}

// runOrDeleteInstance steals the program reference (except on error).
func (s *Server) _runOrDeleteInstance(ctx context.Context, inst *Instance, prog *program, function string) {
	defer s.unrefProgram(&prog)

	drive, err := inst.startOrAnnihilate(prog)
	if err != nil {
		s.deleteNonexistentInstance(inst)
		panic(err)
	}

	if drive {
		go s.driveInstance(detachedContext(ctx), inst, prog, function)
		prog = nil
	}
}

// driveInstance steals the program reference.
func (s *Server) driveInstance(ctx context.Context, inst *Instance, prog *program, function string) {
	defer s.unrefProgram(&prog)

	if nonexistent := inst.drive(ctx, prog, function, s.Monitor); nonexistent {
		s.deleteNonexistentInstance(inst)
	}
}

func (s *Server) _getInstance(ctx context.Context, instance string) *Instance {
	_, inst := s._getInstanceProgramID(ctx, instance)
	return inst
}

func (s *Server) _getInstanceProgramID(ctx context.Context, instance string) (string, *Instance) {
	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	inst, prog := s._getInstanceBorrowProgram(lock, pri, instance)
	return prog.id, inst
}

func (s *Server) _getInstanceRefProgram(ctx context.Context, instance string) (*Instance, *program) {
	pri := principal.ContextID(ctx)
	if pri == nil {
		panic(errAnonymous)
	}

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	inst, prog := s._getInstanceBorrowProgram(lock, pri, instance)
	return inst, prog.ref(lock)
}

func (s *Server) _getInstanceBorrowProgram(_ serverLock, pri *principal.ID, instance string) (*Instance, *program) {
	acc := s.accounts[principal.Raw(pri)]
	if acc == nil {
		panic(resourcenotfound.ErrInstance)
	}

	x, found := acc.instances[instance]
	if !found {
		panic(resourcenotfound.ErrInstance)
	}

	return x.inst, x.prog
}

func (s *Server) _allocateInstanceResources(ctx context.Context, policy *InstancePolicy) (*runtime.Process, InstanceServices) {
	if policy.Services == nil {
		panic(AccessForbidden("no service policy"))
	}

	services := policy.Services(ctx)
	defer func() {
		if services != nil {
			services.Close()
		}
	}()

	proc, err := s.ProcessFactory.NewProcess(ctx)
	_check(err)

	ss := services
	services = nil
	return proc, ss
}

// registerProgramRefInstance with server, and an account if ref is true.
// Caller's instance image is stolen (except on error).  Caller's program
// reference is replaced with a reference to the canonical program object.
func (s *Server) _registerProgramRefInstance(ctx context.Context, acc *account, prog *program, instImage *image.Instance, policy *InstancePolicy, know *api.ModuleOptions, launch *api.LaunchOptions, invoke *InvokeOptions) (inst *Instance, canonicalProg *program, redundantProg bool) {
	var (
		proc     *runtime.Process
		services InstanceServices
	)
	if !launch.Suspend && !instImage.Final() {
		proc, services = s._allocateInstanceResources(ctx, policy)
		defer closeInstanceResources(&proc, &services)
	}

	if know.Pin || !launch.Transient {
		if acc == nil {
			panic(errAnonymous)
		}
		prog._ensureStorage()
	}

	instance := launch.Instance
	if instance == "" {
		instance = makeInstanceID()
	}

	lock := s.mu.Lock()
	defer s.mu.Unlock()

	if acc != nil {
		if s.accounts == nil {
			panic(ErrServerClosed)
		}
		acc._checkUniqueInstanceID(lock, instance)
	}

	prog, redundantProg = s._mergeProgramRef(lock, prog)

	inst = newInstance(instance, acc, launch.Transient, instImage, prog.buffers, proc, services, policy.TimeResolution, launch.Tags, invoke)
	proc = nil
	services = nil

	if acc != nil {
		if know.Pin {
			// mergeProgramRef checked for shutdown, so ensureProgramRef is
			// safe to call.
			if acc.ensureProgramRef(lock, prog, know.Tags) {
				// TODO: move outside of critical section
				s.monitor(&event.ModulePin{
					Ctx:      ContextDetail(ctx),
					Module:   prog.id,
					TagCount: int32(len(know.Tags)),
				})
			}
		}
		acc.instances[instance] = accountInstance{inst, prog.ref(lock)}
	} else {
		s.anonymous[inst] = struct{}{}
	}

	canonicalProg = prog.ref(lock)
	return
}

func (s *Server) deleteNonexistentInstance(inst *Instance) {
	lock := s.mu.Lock()
	defer s.mu.Unlock()

	if inst.acc != nil {
		if x := inst.acc.instances[inst.ID]; x.inst == inst {
			delete(inst.acc.instances, inst.ID)
			x.prog.unref(lock)
		}
	} else {
		delete(s.anonymous, inst)
	}
}

// mergeProgramRef steals the program reference and returns a borrowed program
// reference which is valid until the server mutex is unlocked.
func (s *Server) _mergeProgramRef(lock serverLock, prog *program) (canonical *program, redundant bool) {
	switch existing := s.programs[prog.id]; existing {
	case nil:
		if s.programs == nil {
			panic(ErrServerClosed)
		}
		s.programs[prog.id] = prog // Pass reference to map.
		return prog, false

	case prog:
		if prog.refCount < 2 {
			panic("unexpected program reference count")
		}
		prog.unref(lock) // Map has reference; safe to drop temporary reference.
		return prog, false

	default:
		prog.unref(lock)
		return existing, true
	}
}

func _prepareModuleOptions(opt *api.ModuleOptions) *api.ModuleOptions {
	if opt == nil {
		return new(api.ModuleOptions)
	}
	return opt
}

func _prepareLaunchOptions(opt *api.LaunchOptions) *api.LaunchOptions {
	if opt == nil {
		return new(api.LaunchOptions)
	}
	if opt.Suspend && opt.Function != "" {
		panic(public.Err("function cannot be specified for suspended instance"))
	}
	return opt
}

func prepareResumeOptions(opt *api.ResumeOptions) *api.ResumeOptions {
	if opt == nil {
		return new(api.ResumeOptions)
	}
	return opt
}

func prepareInstanceUpdate(opt *api.InstanceUpdate) *api.InstanceUpdate {
	if opt == nil {
		return new(api.InstanceUpdate)
	}
	return opt
}

func closeProgramImage(p **image.Program) {
	if *p != nil {
		(*p).Close()
		*p = nil
	}
}

func closeInstanceImage(p **image.Instance) {
	if *p != nil {
		(*p).Close()
		*p = nil
	}
}

func closeInstanceResources(proc **runtime.Process, services *InstanceServices) {
	if *proc != nil {
		(*proc).Close()
		*proc = nil
	}
	if *services != nil {
		(*services).Close()
		*services = nil
	}
}

func newInstanceCreateEvent(instance, module string, launch *api.LaunchOptions) *event.InstanceCreate {
	return &event.InstanceCreate{
		Instance:  instance,
		Module:    module,
		Transient: launch.Transient,
		Suspended: launch.Suspend,
		TagCount:  int32(len(launch.Tags)),
	}
}

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

func _context(ctx context.Context, err error) context.Context {
	_check(err)
	return ctx
}

// asError panics on runtime error.
func asError(x interface{}) error {
	if x == nil {
		return nil
	}

	err, _ := x.(error)
	if err == nil {
		panic(x)
	}
	if _, ok := err.(goruntime.Error); ok {
		panic(x)
	}

	return err
}
