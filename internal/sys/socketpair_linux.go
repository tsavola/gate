// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sys

import (
	"fmt"
	"os"
	"syscall"
)

func SocketFilePair(flags int) (f1, f2 *os.File, err error) {
	p, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC|flags, 0)
	if err != nil {
		err = fmt.Errorf("socketpair: %v", err)
		return
	}

	err = syscall.SetNonblock(p[1], true)
	if err != nil {
		err = fmt.Errorf("set nonblock: %v", err)
		return
	}

	f1 = os.NewFile(uintptr(p[0]), "unix")
	f2 = os.NewFile(uintptr(p[1]), "unix")
	return
}
