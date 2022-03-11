// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import (
	"sync/atomic"
	"unsafe"
)

// A Value provides an atomic load and store of a consistently typed value.
// The zero value for a Value returns nil from Load.
// Once Store has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value[T any] struct {
	data unsafe.Pointer
}

// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value[T]) Load() (val T) {
	p := atomic.LoadPointer(&v.data)
	if p == nil {
		return
	}
	return *(*T)(p)
}

// Store sets the value of the Value to x.
func (v *Value[T]) Store(val T) {
	atomic.StorePointer(&v.data, unsafe.Pointer(&val))
}

// Swap stores new into Value and returns the previous value. It returns nil if
// the Value is empty.
func (v *Value[T]) Swap(new T) (old T) {
	p := atomic.SwapPointer(&v.data, unsafe.Pointer(&new))
	if p == nil {
		return
	}
	return *(*T)(p)
}

type ComparableValue[T comparable] struct {
	Value[T]
}

// CompareAndSwap executes the compare-and-swap operation for the Value.
func (v *ComparableValue[T]) CompareAndSwap(old, new T) (swapped bool) {
	// Compare old and current via runtime equality check.
	// This allows value types to be compared, something
	// not offered by the package functions.
	// CompareAndSwapPointer below only ensures vp.data
	// has not changed since LoadPointer.
	data := atomic.LoadPointer(&v.data)
	if data == nil {
		var zero T
		if old != zero {
			return false
		}
	}
	if old != (*(*T)(data)) {
		return false
	}
	return atomic.CompareAndSwapPointer(&v.data, data, unsafe.Pointer(&new))
}
