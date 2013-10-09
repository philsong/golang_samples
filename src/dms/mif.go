package mif

import "syscall"
import "unsafe"
import (
	"fmt"
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func init() {
	fmt.Print("Starting Up\n")
}

func exit() {
	fmt.Print("Ending\n")
	defer syscall.FreeLibrary(kernel32)
	defer syscall.FreeLibrary(user32)
	defer syscall.FreeLibrary(miflib)
}

var (
	kernel32, _        = syscall.LoadLibrary("kernel32.dll")
	getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")

	user32, _     = syscall.LoadLibrary("user32.dll")
	messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")

	miflib, _            = syscall.LoadLibrary("miflib.dll")
	cHVD_OpenPort, _     = syscall.GetProcAddress(miflib, "HVD_OpenPort")
	cHVD_ClosePort, _    = syscall.GetProcAddress(miflib, "HVD_ClosePort")
	cMIF_REQ_ANTI_SEL, _ = syscall.GetProcAddress(miflib, "MIF_REQ_ANTI_SEL")
	cMIF_HALT, _         = syscall.GetProcAddress(miflib, "MIF_HALT")
	cMIF_AUTH_KEY, _     = syscall.GetProcAddress(miflib, "MIF_AUTH_KEY")
	cMIF_Read, _         = syscall.GetProcAddress(miflib, "MIF_Read")
	cMIF_Write, _        = syscall.GetProcAddress(miflib, "MIF_Write")
)

const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND

	MB_DEFBUTTON1 = 0x00000000
	MB_DEFBUTTON2 = 0x00000100
	MB_DEFBUTTON3 = 0x00000200
	MB_DEFBUTTON4 = 0x00000300
)

func MessageBox(caption, text string, style uintptr) (result int) {
	var nargs uintptr = 4
	ret, _, callErr := syscall.Syscall9(uintptr(messageBox),
		nargs,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style,
		0,
		0,
		0,
		0,
		0)
	if callErr != 0 {
		abort("Call MessageBox", callErr)
	}
	result = int(ret)
	return
}

func HVD_OpenPort(CommNumber byte) (handle uintptr) {
	var nargs uintptr = 1
	if ret, _, callErr := syscall.Syscall(uintptr(cHVD_OpenPort),
		nargs,
		uintptr(CommNumber),
		0,
		0); callErr != 0 {
		abort("Call HVD_OpenPort", callErr)
	} else {
		handle = ret
	}

	return
}

func HVD_ClosePort(handle uintptr) (result uintptr) {
	var nargs uintptr = 1
	if ret, _, callErr := syscall.Syscall(uintptr(cHVD_ClosePort),
		nargs,
		handle,
		0,
		0); callErr != 0 {
		abort("Call HVD_OpenPort", callErr)
	} else {
		result = ret
	}

	return
}

func MIF_REQ_ANTI_SEL(handle uintptr) (result uintptr) {
	var nargs uintptr = 1
	if ret, _, callErr := syscall.Syscall(uintptr(cMIF_REQ_ANTI_SEL),
		nargs,
		handle,
		0,
		0); callErr != 0 {
		abort("Call HVD_OpenPort", callErr)
	} else {
		result = ret
	}
	return
}

func GetModuleHandle() (handle uintptr) {
	var nargs uintptr = 0
	if ret, _, callErr := syscall.Syscall(uintptr(getModuleHandle), nargs, 0, 0, 0); callErr != 0 {
		abort("Call GetModuleHandle", callErr)
	} else {
		handle = ret
	}
	return
}
