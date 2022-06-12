package unsafelib

import (
	"reflect"
	"syscall"
	"unsafe"
)

var (
	PROT_FLAG_WRITE = syscall.PROT_WRITE
	PROT_FLAG_READ  = syscall.PROT_READ

	PROT_FLAG_RW = PROT_FLAG_READ | PROT_FLAG_WRITE
	PROT_FLAG_RO = PROT_FLAG_READ
)

// mprotect
func mprotect(addr uintptr, len int, prot int) (err error) {
	page := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: addr,
		Len:  len,
	}))
	err = syscall.Mprotect(page, prot)
	return
}

type mem_profile struct {
	addr    uintptr // page begin addr
	len     int
	oldFlag int
	newFlag int
	state   bool // default: false
}

func (m *mem_profile) toggle() (err error, setNewFlag bool) {
	setNewFlag = !m.state
	var flag = m.oldFlag
	if setNewFlag {
		flag = m.newFlag
	}
	err = mprotect(m.addr, m.len, flag)
	return
}

// Testing android mprotect: https://stackoverflow.com/questions/9565056/android-mprotect-not-changing-protections
