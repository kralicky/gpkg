//go:build race

package sync

import (
	"bytes"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

func TestUnsafeExpunged(t *testing.T) {
	var ptr1 atomic.Pointer[time.Time]
	tp := time.Now()
	ptr1.Store(&tp)
	ptr1.CompareAndSwap(&tp, (*time.Time)(expunged))

	var ptr2 atomic.Pointer[io.Reader]
	var iface io.Reader = &bytes.Buffer{}
	ptr2.Store(&iface)
	ptr2.CompareAndSwap(&iface, (*io.Reader)(expunged))
}
