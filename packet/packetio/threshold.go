// Copyright (c) 2019 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packetio

import (
	"sync/atomic"

	"github.com/tsavola/gate/internal/error/badprogram"
)

const errNegativeIncrement = badprogram.Err("stream flow increment is negative")

// Threshold is an observable scalar value (uint32) which gets incremented.
// The producer calls Increase and Finish, and the consumer calls Changed and
// Current.
type Threshold struct {
	c chan struct{}
	n uint32 // Atomic.
}

// MakeThreshold is for initializing a field.  The value must not be copied.
func MakeThreshold() Threshold {
	return Threshold{
		c: make(chan struct{}, 1),
	}
}

func NewThreshold() *Threshold {
	t := MakeThreshold()
	return &t
}

// Increase the threshold value.  The value may wrap around.
func (t *Threshold) Increase(increment int32) (err error) {
	if increment < 0 {
		err = errNegativeIncrement
		return
	}

	atomic.AddUint32(&t.n, uint32(increment))

	select {
	case t.c <- struct{}{}:
		return

	default:
		return
	}
}

// Finish closes the Changed channel.  Increase must not be called after this.
func (t *Threshold) Finish() {
	close(t.c)
}

// nonatomic is producer-side Current.
func (t *Threshold) nonatomic() uint32 {
	return t.n
}

// Changed channel is unblocked after the threshold has been raised.  It is
// closed by Finish.
func (t *Threshold) Changed() <-chan struct{} {
	return t.c
}

// Current value.  The value may have wrapped around.
func (t *Threshold) Current() uint32 {
	return atomic.LoadUint32(&t.n)
}
