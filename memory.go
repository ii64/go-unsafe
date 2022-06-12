package unsafelib

import (
	"os"
)

var (
	PAGE_SIZE = uintptr(os.Getpagesize())
)

func getpagebase(addr uintptr) uintptr {
	return (addr / PAGE_SIZE) * PAGE_SIZE
}
func getpageend(addr uintptr, i int) uintptr {
	return (addr/PAGE_SIZE + uintptr(i)) * PAGE_SIZE
}

func new_mem_profile(addr uintptr, pagec int, oldFlag, newFlag int) *mem_profile {
	return &mem_profile{
		addr:    getpagebase(addr),
		len:     pagec * int(PAGE_SIZE),
		oldFlag: oldFlag,
		newFlag: newFlag,
		state:   false,
	}
}
