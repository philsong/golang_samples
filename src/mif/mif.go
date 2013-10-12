/**
 * Mifare reader/writer
 *
 * An open source Mifare reader/writer development library for Golang
 *
 * @package		mif
 * @author		Philsong
 * @e-Mail    : 78623269@qq.com
 * @since		Version 1.0
 */
package mif

import (
	"fmt"
	"syscall"
	"unsafe"
)

type MIF_KEY_TYPE int

const (
	KEY_A = iota
	KEY_B
)

type MIF_KEY struct {
	M_Value [6]byte
}

type MIF_DATA_BLOCK struct {
	M_Value [16]byte
}

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func init() {
	fmt.Print("Starting Up\n")
}

func FreeLib() {
	fmt.Print("Ending\n")
	syscall.FreeLibrary(miflib)
	syscall.FreeLibrary(user32)
}

var (
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

func MessageBoxEx(caption, text string, style uintptr) (result int) {
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

func MessageBox(text string) (result int) {
	return MessageBoxEx("alert", text, MB_YESNOCANCEL)
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

func HVD_ClosePort(handle uintptr) {
	var nargs uintptr = 1
	if _, _, callErr := syscall.Syscall(uintptr(cHVD_ClosePort),
		nargs,
		handle,
		0,
		0); callErr != 0 {
		abort("Call HVD_ClosePort", callErr)
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
		abort("Call MIF_REQ_ANTI_SEL", callErr)
	} else {
		result = ret
	}
	return
}

func MIF_AUTH_KEY(handle uintptr,
	KeyType MIF_KEY_TYPE,
	BlockNumber byte,
	pKey *MIF_KEY) (result uintptr) {
	var nargs uintptr = 4
	if ret, _, callErr := syscall.Syscall6(uintptr(cMIF_AUTH_KEY),
		nargs,
		handle,
		uintptr(KeyType),
		uintptr(BlockNumber),
		uintptr(unsafe.Pointer(pKey)),
		0,
		0); callErr != 0 {
		abort("Call MIF_AUTH_KEY", callErr)
	} else {
		result = ret
	}
	return
}

func MIF_Read(handle uintptr,
	BlockNumber byte,
	pBlockData *MIF_DATA_BLOCK) (result uintptr) {
	var nargs uintptr = 3
	if ret, _, callErr := syscall.Syscall6(uintptr(cMIF_Read),
		nargs,
		handle,
		uintptr(BlockNumber),
		uintptr(unsafe.Pointer(pBlockData)),
		0,
		0,
		0); callErr != 0 {
		abort("Call MIF_Read", callErr)
	} else {
		result = ret
	}
	return
}

func MIF_Write(handle uintptr,
	BlockNumber byte,
	pBlockData *MIF_DATA_BLOCK) (result uintptr) {
	var nargs uintptr = 3
	if ret, _, callErr := syscall.Syscall6(uintptr(cMIF_Write),
		nargs,
		handle,
		uintptr(BlockNumber),
		uintptr(unsafe.Pointer(pBlockData)),
		0,
		0,
		0); callErr != 0 {
		abort("Call MIF_Write", callErr)
	} else {
		result = ret
	}
	return
}
