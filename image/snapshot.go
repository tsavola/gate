// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"encoding/binary"
	"syscall"

	"github.com/tsavola/gate/internal/file"
	"github.com/tsavola/gate/internal/manifest"
	"github.com/tsavola/gate/internal/varint"
	"github.com/tsavola/gate/snapshot"
	"github.com/tsavola/gate/snapshot/wasm"
	"github.com/tsavola/wag/object/abi"
	"github.com/tsavola/wag/section"
	"github.com/tsavola/wag/wa"
	"github.com/tsavola/wag/wa/opcode"
)

const wasmModuleHeaderSize = 8

func Snapshot(oldProg *Program, inst *Instance, buffers snapshot.Buffers, suspended bool,
) (newProg *Program, err error) {
	// Old program.
	var (
		man       = oldProg.Manifest()
		oldRanges = man.Sections
		oldFD     = oldProg.file.Fd()
	)

	// Instance file.
	var (
		instStackOffset   = int64(0)
		instGlobalsOffset = instStackOffset + int64(inst.man.StackSize)
		instMemoryOffset  = instGlobalsOffset + alignPageOffset32(inst.man.GlobalsSize)
	)

	instMap, err := mmap(inst.file.Fd(), instStackOffset, int(instMemoryOffset-instStackOffset), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return
	}
	defer mustMunmap(instMap)

	var (
		instStackMapping   = instMap[:inst.man.StackSize]
		instGlobalsMapping = instMap[inst.man.StackSize:]
	)

	instStackUnused, memorySize, ok := checkStack(instStackMapping, len(instStackMapping))
	if !ok {
		err = ErrInvalidState
		return
	}
	if memorySize > uint32(inst.man.MaxMemorySize) {
		err = ErrInvalidState
		return
	}

	// Stack, globals and memory contents without unused regions or padding.
	var (
		newInitRoutine uint32
		newTextAddr    uint64
		instStackData  []byte
		stackUsage     int
		globalsData    = instGlobalsMapping[len(instGlobalsMapping)-len(man.GlobalTypes)*8:]
	)

	if suspended {
		if instStackUnused != 0 {
			newInitRoutine = abi.TextAddrResume
			newTextAddr = inst.man.TextAddr
			instStackData = instStackMapping[instStackUnused:]
			stackUsage = len(instStackData)
		} else {
			// Starting is equivalent to resuming at virtual call site at
			// beginning of start routine.
			newInitRoutine = abi.TextAddrStart
			stackUsage = 16
		}
	} else {
		// New program reuses old text segment which may invoke start function.
		// Don't invoke it again.
		newInitRoutine = abi.TextAddrEnter
	}

	// New module sections.
	// TODO: align section contents to facilitate reflinking?
	// TODO: stitch module together during download?
	newRanges := make([]manifest.ByteRange, section.Data+1)

	off := int64(wasmModuleHeaderSize)
	off = mapOldSection(off, newRanges, oldRanges, section.Type)
	off = mapOldSection(off, newRanges, oldRanges, section.Import)
	off = mapOldSection(off, newRanges, oldRanges, section.Function)
	off = mapOldSection(off, newRanges, oldRanges, section.Table)

	memorySection := makeMemorySection(memorySize) // TODO: maximum value
	off = mapNewSection(off, newRanges, len(memorySection), section.Memory)

	globalSection := makeGlobalSection(man.GlobalTypes, globalsData)
	off = mapNewSection(off, newRanges, len(globalSection), section.Global)

	off = mapOldSection(off, newRanges, oldRanges, section.Export)
	off = mapOldSection(off, newRanges, oldRanges, section.Element)
	off = mapOldSection(off, newRanges, oldRanges, section.Code)

	bufferHeader, bufferSectionSize := makeBufferSectionHeader(buffers)
	bufferSectionOffset := off
	off += bufferSectionSize

	var (
		stackHeader        []byte
		stackSectionSize   int
		stackSectionOffset int64
	)
	if suspended {
		stackHeader = makeStackSectionHeader(stackUsage)
		stackSectionSize = len(stackHeader) + stackUsage
		stackSectionOffset = off
		off += int64(stackSectionSize)
	}

	dataHeader := makeDataSectionHeader(int(memorySize))
	dataSectionSize := len(dataHeader) + int(memorySize)
	off = mapNewSection(off, newRanges, dataSectionSize, section.Data)

	// New module size.
	newModuleSize := man.ModuleSize
	newModuleSize -= man.Sections[section.Memory].Length
	newModuleSize -= man.Sections[section.Global].Length
	newModuleSize -= man.Sections[section.Start].Length
	newModuleSize -= man.BufferSection.Length
	newModuleSize -= man.StackSection.Length
	newModuleSize -= man.Sections[section.Data].Length
	newModuleSize += int64(len(memorySection))
	newModuleSize += int64(len(globalSection))
	newModuleSize += bufferSectionSize
	newModuleSize += int64(stackSectionSize)
	newModuleSize += int64(dataSectionSize)

	// New program file.
	newFile, err := oldProg.storage.newProgramFile()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			newFile.Close()
		}
	}()

	// Copy module header and sections up to and including table section.
	copyLen := int(wasmModuleHeaderSize)
	for i := section.Table; i >= section.Type; i-- {
		if newRanges[i].Offset != 0 {
			copyLen = int(newRanges[i].Offset + newRanges[i].Length)
			break
		}
	}

	newOff := progModuleOffset
	oldOff := progModuleOffset
	err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, copyLen)
	if err != nil {
		return
	}

	// Write new memory and global section, and skip old ones.
	n, err := newFile.WriteAt(append(memorySection, globalSection...), newOff)
	if err != nil {
		return
	}
	newOff += int64(n)
	oldOff += oldRanges[section.Memory].Length + oldRanges[section.Global].Length

	// If there is a start section, copy export section separately, and skip
	// start section.
	nextSection := section.Export

	if oldRanges[section.Start].Length > 0 {
		err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, int(oldRanges[section.Export].Length))
		if err != nil {
			return
		}

		oldOff += oldRanges[section.Start].Length
		nextSection = section.Element
	}

	// Copy sections up to and including code section.
	copyLen = 0
	for _, s := range oldRanges[nextSection : section.Code+1] {
		if s.Length > 0 {
			copyLen += int(s.Length)
		}
	}

	err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, copyLen)
	if err != nil {
		return
	}

	// Write new buffer section, and skip old one.
	if bufferSectionSize > 0 {
		n, err = newFile.WriteAt(bufferHeader, newOff)
		if err != nil {
			return
		}
		newOff += int64(n)

		n, err = writeBufferSectionDataAt(newFile, buffers, newOff)
		if err != nil {
			return
		}
		newOff += int64(n)
	}
	oldOff += man.BufferSection.Length

	// Write new stack section, and skip old one.
	if suspended {
		newStackSection := make([]byte, len(stackHeader)+stackUsage)
		copy(newStackSection, stackHeader)

		newStack := newStackSection[len(stackHeader):]

		if instStackData != nil {
			err = exportStack(newStack, instStackData, inst.man.TextAddr, oldProg.Map) // TODO: in-place?
			if err != nil {
				return
			}
		} else {
			// Synthesize portable stack, suspended at virtual call site at
			// index 0 (beginning of start routine).
			binary.LittleEndian.PutUint64(newStack[0:], 0)
			binary.LittleEndian.PutUint32(newStack[8:], inst.man.EntryIndex)
		}

		n, err = newFile.WriteAt(newStackSection, newOff)
		if err != nil {
			return
		}
		newOff += int64(n)
	}
	oldOff += man.StackSection.Length

	// Copy new data section from instance, and skip old one.
	n, err = newFile.WriteAt(dataHeader, newOff)
	if err != nil {
		return
	}
	newOff += int64(n)

	if oldProg.storage.singleBackend() {
		instOff := instMemoryOffset
		err = copyFileRange(inst.file.Fd(), &instOff, newFile.Fd(), &newOff, int(memorySize))
		if err != nil {
			return
		}
	} else {
		panic("TODO")
	}
	oldOff += oldRanges[section.Data].Length

	// Copy remaining (custom) sections.
	copyLen = int(man.ModuleSize - (oldOff - progModuleOffset))
	err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, copyLen)
	if err != nil {
		return
	}

	// Copy object map from program.
	newOff = align8(newOff)
	oldOff = align8(oldOff)
	err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, int(man.CallSitesSize)+int(man.FuncAddrsSize))
	if err != nil {
		return
	}

	// Copy text from program.
	newOff = progTextOffset
	oldOff = progTextOffset
	err = copyFileRange(oldFD, &oldOff, newFile.Fd(), &newOff, alignPageSize32(man.TextSize))
	if err != nil {
		return
	}

	var newProgMem []byte

	if oldProg.storage.singleBackend() {
		// Copy stack from instance (again).
		if instStackData != nil {
			copyLen := alignPageSize(stackUsage)
			newOff = progGlobalsOffset - int64(copyLen)
			instOff := instGlobalsOffset - int64(copyLen)
			err = copyFileRange(inst.file.Fd(), &instOff, newFile.Fd(), &newOff, copyLen)
			if err != nil {
				return
			}
		}

		// Copy globals and memory from instance (again).
		newOff = progGlobalsOffset
		instOff := instGlobalsOffset
		err = copyFileRange(inst.file.Fd(), &instOff, newFile.Fd(), &newOff, alignPageSize32(inst.man.GlobalsSize)+int(memorySize))
		if err != nil {
			return
		}
	} else {
		panic("TODO")
	}

	// New program manifest.
	man.TextAddr = newTextAddr
	man.StackUsage = uint32(stackUsage)
	man.MemoryDataSize = memorySize
	man.MemorySize = memorySize
	man.InitRoutine = newInitRoutine
	man.ModuleSize = newModuleSize
	man.Sections = newRanges
	man.BufferSection = manifest.ByteRange{
		Offset: bufferSectionOffset,
		Length: bufferSectionSize,
	}
	man.BufferSectionHeaderLength = int64(len(bufferHeader))
	man.StackSection = manifest.ByteRange{
		Offset: stackSectionOffset,
		Length: int64(stackSectionSize),
	}

	newProg = &Program{
		Map:     oldProg.Map,
		storage: oldProg.storage,
		man:     man,
		file:    newFile,
		mem:     newProgMem,
	}
	return
}

func makeMemorySection(currentMemorySize uint32) []byte {
	buf := make([]byte, 4+binary.MaxVarintLen32)

	n := binary.PutUvarint(buf[4:], uint64(currentMemorySize>>wa.PageBits)) // Initial value
	buf[3] = 0                                                              // Maximum flag
	buf[2] = 1                                                              // Item count
	buf[1] = byte(2 + n)                                                    // Payload length
	buf[0] = byte(section.Memory)                                           // Section id

	return buf[:4+n]
}

func makeGlobalSection(globalTypes []byte, segment []byte) []byte {
	const (
		// Section id, payload length, item count.
		maxHeaderSize = 1 + binary.MaxVarintLen32 + binary.MaxVarintLen32

		// Type, mutable flag, const op, const value, end op.
		maxItemSize = 1 + 1 + 1 + binary.MaxVarintLen64 + 1
	)

	buf := make([]byte, maxHeaderSize+len(globalTypes)*maxItemSize)

	// Items:
	itemsSize := putGlobals(buf[maxHeaderSize:], globalTypes, segment)

	// Header:
	countSize := putVaruint32Before(buf, maxHeaderSize, uint32(len(globalTypes)))
	payloadLen := countSize + itemsSize
	payloadLenSize := putVaruint32Before(buf, maxHeaderSize-countSize, uint32(payloadLen))
	buf[maxHeaderSize-countSize-payloadLenSize-1] = byte(section.Global)

	return buf[maxHeaderSize-countSize-payloadLenSize-1 : maxHeaderSize+itemsSize]
}

func putGlobals(target []byte, globalTypes []byte, segment []byte) (totalSize int) {
	for _, b := range globalTypes {
		t := wa.GlobalType(b)

		value := binary.LittleEndian.Uint64(segment)
		segment = segment[8:]

		encoded := t.Encode()
		n := copy(target, encoded[:])
		target = target[n:]
		totalSize += n

		switch t.Type() {
		case wa.I32:
			target[0] = byte(opcode.I32Const)
			n = 1 + binary.PutVarint(target[1:], int64(int32(uint32(value))))

		case wa.I64:
			target[0] = byte(opcode.I64Const)
			n = 1 + binary.PutVarint(target[1:], int64(value))

		case wa.F32:
			target[0] = byte(opcode.F32Const)
			binary.LittleEndian.PutUint32(target[1:], uint32(value))
			n = 1 + 4

		case wa.F64:
			target[0] = byte(opcode.F64Const)
			binary.LittleEndian.PutUint64(target[1:], value)
			n = 1 + 8

		default:
			panic(t)
		}
		target = target[n:]
		totalSize += n

		target[0] = byte(opcode.End)
		target = target[1:]
		totalSize++
	}

	return
}

func makeBufferSectionHeader(buffers snapshot.Buffers) (header []byte, sectionSize int64) {
	if len(buffers.Services) == 0 && len(buffers.Input) == 0 && len(buffers.Output) == 0 && buffers.Flags == 0 {
		return
	}

	// Section id, payload length.
	const maxSectionFrameSize = 1 + varint.MaxLen

	maxHeaderSize := maxSectionFrameSize
	maxHeaderSize += 1                       // Section name length
	maxHeaderSize += len(wasm.SectionBuffer) // Section name
	maxHeaderSize += varint.MaxLen           // Flags
	maxHeaderSize += varint.MaxLen           // Input data size
	maxHeaderSize += varint.MaxLen           // Output data size
	maxHeaderSize += varint.MaxLen           // Service count

	for _, s := range buffers.Services {
		maxHeaderSize += 1             // Service name length
		maxHeaderSize += len(s.Name)   // Service name
		maxHeaderSize += varint.MaxLen // Service data size
	}

	buf := make([]byte, maxHeaderSize)

	tail := buf[maxSectionFrameSize:]
	tail = putByte(tail, byte(len(wasm.SectionBuffer)))
	tail = putString(tail, wasm.SectionBuffer)
	tail = varint.Put(tail, int32(buffers.Flags))
	tail = varint.Put(tail, int32(len(buffers.Input)))
	tail = varint.Put(tail, int32(len(buffers.Output)))
	tail = varint.Put(tail, int32(len(buffers.Services)))

	dataSize := int64(len(buffers.Input)) + int64(len(buffers.Output))

	for _, s := range buffers.Services {
		tail = putByte(tail, byte(len(s.Name)))
		tail = putString(tail, s.Name)
		tail = varint.Put(tail, int32(len(s.Buffer)))

		dataSize += int64(len(s.Buffer))
	}

	payloadLen := uint32(len(buf)-len(tail)-maxSectionFrameSize) + uint32(dataSize)
	payloadLenSize := putVaruint32Before(buf, maxSectionFrameSize, payloadLen)
	buf[maxSectionFrameSize-payloadLenSize-1] = byte(section.Custom)
	buf = buf[maxSectionFrameSize-payloadLenSize-1:]

	header = buf[:len(buf)-len(tail)]
	sectionSize = int64(len(header)) + dataSize
	return
}

func writeBufferSectionDataAt(f *file.File, bs snapshot.Buffers, off int64) (total int, err error) {
	n, err := f.WriteAt(bs.Input, off)
	if err != nil {
		return
	}
	total += n
	off += int64(n)

	n, err = f.WriteAt(bs.Output, off)
	if err != nil {
		return
	}
	total += n
	off += int64(n)

	for _, s := range bs.Services {
		n, err = f.WriteAt(s.Buffer, off)
		if err != nil {
			return
		}
		total += n
		off += int64(n)
	}

	return
}

func makeStackSectionHeader(stackSize int) []byte {
	// Section id, payload length.
	const maxSectionFrameSize = 1 + binary.MaxVarintLen32

	// Name length, name string.
	var customHeaderSize = 1 + len(wasm.SectionStack)

	buf := make([]byte, maxSectionFrameSize+customHeaderSize)
	buf[0] = byte(section.Custom)
	payloadLenSize := binary.PutUvarint(buf[1:], uint64(customHeaderSize+stackSize))
	buf[1+payloadLenSize] = byte(len(wasm.SectionStack))
	copy(buf[1+payloadLenSize+1:], wasm.SectionStack)

	return buf[:1+payloadLenSize+1+len(wasm.SectionStack)]
}

func makeDataSectionHeader(memorySize int) []byte {
	const (
		// Section id, payload length.
		maxSectionFrameSize = 1 + binary.MaxVarintLen32

		// Count, memory index, init expression, size.
		maxSegmentHeaderSize = 1 + 1 + 3 + binary.MaxVarintLen32
	)

	buf := make([]byte, maxSectionFrameSize+maxSegmentHeaderSize)

	segment := buf[maxSectionFrameSize:]
	segment[0] = 1 // Count
	segment[1] = 0 // Memory index
	segment[2] = byte(opcode.I32Const)
	segment[3] = 0 // Offset
	segment[4] = byte(opcode.End)
	segmentHeaderSize := 5 + binary.PutUvarint(segment[5:], uint64(memorySize))

	payloadLen := segmentHeaderSize + memorySize
	payloadLenSize := putVaruint32Before(buf, maxSectionFrameSize, uint32(payloadLen))
	buf[maxSectionFrameSize-payloadLenSize-1] = byte(section.Data)

	return buf[maxSectionFrameSize-payloadLenSize-1 : maxSectionFrameSize+segmentHeaderSize]
}

func putByte(dest []byte, x byte) (tail []byte) {
	dest[0] = x
	return dest[1:]
}

func putString(dest []byte, s string) (tail []byte) {
	copy(dest, s)
	return dest[len(s):]
}

func putVaruint32Before(dest []byte, offset int, x uint32) (n int) {
	var temp [binary.MaxVarintLen32]byte
	n = binary.PutUvarint(temp[:], uint64(x))
	copy(dest[offset-n:], temp[:n])
	return
}

func mapOldSection(offset int64, dest, src []manifest.ByteRange, i section.ID) int64 {
	if src[i].Length > 0 {
		dest[i] = manifest.ByteRange{
			Offset: offset,
			Length: src[i].Length,
		}
		offset += src[i].Length
	}
	return offset
}

func mapNewSection(offset int64, dest []manifest.ByteRange, length int, i section.ID) int64 {
	dest[i] = manifest.ByteRange{
		Offset: offset,
		Length: int64(length),
	}
	return offset + int64(length)
}
