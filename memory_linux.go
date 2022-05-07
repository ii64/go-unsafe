package unsafelib

import (
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

var (
	PAGE_SIZE = uintptr(os.Getpagesize())
)

var (
	memPageMap = map[uintptr]int{}
	memPageMu  sync.Mutex
)

func getpagebase(addr uintptr) uintptr {
	return (addr / PAGE_SIZE) * PAGE_SIZE
}
func getpageend(addr uintptr, i int) uintptr {
	return (addr/PAGE_SIZE + uintptr(i)) * PAGE_SIZE
}

// mprotect
func mprotect(addr uintptr, len int, prot int) (err error) {
	page := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: addr,
		Len:  len,
	}))
	err = syscall.Mprotect(page, prot)
	return
}
