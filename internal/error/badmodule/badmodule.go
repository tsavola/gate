// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package badmodule

import (
	werrors "gate.computer/wag/errors"
)

type Error = werrors.ModuleError

type Dual struct {
	Private string
	Public  string
}

func (x *Dual) Error() string       { return x.Private }
func (x *Dual) PublicError() string { return x.Public }
func (x *Dual) ModuleError()        {}
