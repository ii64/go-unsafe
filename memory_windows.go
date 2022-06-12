package unsafelib

import (
	"syscall"
	"unsafe"
)

var (
	// https://docs.microsoft.com/en-us/windows/win32/memory/memory-protection-constants
	PROT_FLAG_READ  = syscall.PAGE_READONLY
	PROT_FLAG_WRITE = 0 // unk

	PROT_FLAG_RW = syscall.PAGE_READWRITE
	PROT_FLAG_RO = PROT_FLAG_READ
)

var (
	hKernel32Dll                = syscall.NewLazyDLL("KERNEL32.DLL")
	hKernel32Dll_VirtualProtect = hKernel32Dll.NewProc("VirtualProtect")
)

func virtualProtect(addr uintptr, len int, prot int, oldProt *int) (err error) {
	var r1 uintptr
	r1, _, err = hKernel32Dll_VirtualProtect.Call(addr, uintptr(len),
		uintptr(prot), uintptr(unsafe.Pointer(oldProt)))
	if r1 == 0 { // BOOL
		return
	}
	err = nil
	return
}

type mem_profile struct {
	addr    uintptr
	len     int
	oldFlag int
	newFlag int
	state   bool
}

func (m *mem_profile) toggle() (err error, setNewFlag bool) {
	err = virtualProtect(m.addr, m.len, m.newFlag, &m.oldFlag)
	if err == nil {
		m.state = !m.state
		setNewFlag = m.state
		m.newFlag = m.oldFlag
	}
	return
}
