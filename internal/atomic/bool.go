package atomic

import (
	"sync/atomic"
)

// Atomic wrapper around bool.
type Bool uint32

func (a *Bool) CompareAndSwap(o, n bool) bool {
	return atomic.CompareAndSwapUint32((*uint32)(a), boolToUint(o), boolToUint(n))
}

func boolToUint(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
