package unsafelib

import (
	"syscall"
)

// mem_mkrw page read-write,
func mem_mkrw(addr uintptr, pagec int) error {
	addr = getpagebase(addr)
	return mprotect(addr, pagec*int(PAGE_SIZE), syscall.PROT_READ|syscall.PROT_WRITE)
}

// mem_mkro page read-only
func mem_mkro(addr uintptr, pagec int) error {
	addr = getpagebase(addr)
	return mprotect(addr, pagec*int(PAGE_SIZE), syscall.PROT_READ)
}
