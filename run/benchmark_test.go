// Copyright (c) 2017 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package run_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/tsavola/gate/internal/runtest"
	"github.com/tsavola/gate/run"
	"github.com/tsavola/wag"
	"github.com/tsavola/wag/traps"
)

func readProgram(testName string) []byte {
	f := openProgram(testName)
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return data
}

const (
	benchStackSize = 65536
)

var (
	benchRT = runtest.NewRuntime()

	benchProgNop   = readProgram("nop")
	benchProgHello = readProgram("hello")
	benchProgPeer  = readProgram("peer")
)

func compileBenchmark(prog []byte) (m *wag.Module) {
	m = new(wag.Module)

	err := run.Load(m, bytes.NewReader(prog), benchRT.Runtime, new(bytes.Buffer), nil, nil)
	if err != nil {
		panic(err)
	}

	return
}

func prepareBenchmark(m *wag.Module) (image *run.Image) {
	image = new(run.Image)

	if err := image.Init(); err != nil {
		panic(err)
	}

	_, memorySize := m.MemoryLimits()

	if err := image.Populate(m, memorySize, benchStackSize); err != nil {
		panic(err)
	}

	return
}

func executeBenchmark(image *run.Image, output io.Writer) (exit int, trap traps.Id, err error) {
	var proc run.Process

	err = proc.Init(context.Background(), benchRT.Runtime, image, nil)
	if err != nil {
		return
	}
	defer proc.Close()

	exit, trap, err = run.Run(context.Background(), benchRT.Runtime, &proc, image, &testServiceRegistry{output})
	return
}

func BenchmarkCompileNop(b *testing.B)   { benchmarkCompile(b, benchProgNop) }
func BenchmarkCompileHello(b *testing.B) { benchmarkCompile(b, benchProgHello) }
func BenchmarkCompilePeer(b *testing.B)  { benchmarkCompile(b, benchProgPeer) }

func benchmarkCompile(b *testing.B, prog []byte) {
	for i := 0; i < b.N; i++ {
		compileBenchmark(prog)
	}
}

func BenchmarkPrepareNop(b *testing.B)   { benchmarkPrepare(b, benchProgNop) }
func BenchmarkPrepareHello(b *testing.B) { benchmarkPrepare(b, benchProgHello) }
func BenchmarkPreparePeer(b *testing.B)  { benchmarkPrepare(b, benchProgPeer) }

func benchmarkPrepare(b *testing.B, prog []byte) {
	m := compileBenchmark(prog)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p := prepareBenchmark(m)
		p.Close()
	}
}

func BenchmarkExecuteNop(b *testing.B)   { benchmarkExecute(b, benchProgNop, "") }
func BenchmarkExecuteHello(b *testing.B) { benchmarkExecute(b, benchProgHello, "hello world\n") }

func benchmarkExecute(b *testing.B, prog []byte, expectOutput string) {
	m := compileBenchmark(prog)
	p := prepareBenchmark(m)
	defer p.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var output bytes.Buffer

		exit, trap, err := executeBenchmark(p, &output)
		if err != nil {
			panic(err)
		}
		if trap != 0 {
			panic(trap)
		}
		if exit != 0 {
			panic(exit)
		}

		if output.String() != expectOutput {
			panic(fmt.Sprint(output.Bytes()))
		}
	}
}
