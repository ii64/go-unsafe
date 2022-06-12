package unsafelib

import (
	"syscall"
	"testing"
	"unsafe"
)

func TestVirtualProtect(t *testing.T) {
	orig := "hello"
	s := (*String)(unsafe.Pointer(&orig))
	println(s.Data)
	newProtectFlag := PROT_FLAG_RW
	oldProtectFlag := 0
	err := virtualProtect(getpagebase(uintptr(s.Data)), int(PAGE_SIZE), newProtectFlag, &oldProtectFlag)
	if err != nil {
		t.Fail()
	}
	println(oldProtectFlag == syscall.PAGE_READONLY)
	(*[5]byte)(s.Data)[1] = 'a'
	println(orig)
}

func TestMemProfile(t *testing.T) {
	orig := "hello"
	s := (*String)(unsafe.Pointer(&orig))
	println(s.Data)
	prf := new_mem_profile(uintptr(s.Data), 1, PROT_FLAG_RO, PROT_FLAG_RW)
	prf.toggle()
	(*[5]byte)(s.Data)[1] = 'a'
	prf.toggle()
	// prf.toggle()
	// (*[5]byte)(s.Data)[1] = 'a'

	println(orig)
}
