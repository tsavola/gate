// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <gate.h>

int debug(void)
{
	gate_debug("hello, ");
	gate_debug("world\n");
	return 0;
}
