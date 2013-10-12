// list smart card driver in windows
package main

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	WinSCard, _                  = syscall.LoadLibrary(`C:\windows\system32\WinSCard.dll`)
	procSCardEstablishContext, _ = syscall.GetProcAddress(WinSCard, "SCardEstablishContext")
	procSCardReleaseContext, _   = syscall.GetProcAddress(WinSCard, "SCardReleaseContext")
	procSCardListReaders, _      = syscall.GetProcAddress(WinSCard, "SCardListReadersW")
)

const (
	SCARD_SCOPE_USER   = 0
	SCARD_SCOPE_SYSTEM = 2

	SCARD_ALL_READERS     = "SCard$AllReaders"
	SCARD_DEFAULT_READERS = "SCard$DefaultReaders"
)

func SCardListReaders(hContext syscall.Handle, mszGroups *uint16, mszReaders *uint16, pcchReaders *uint32) (retval error) {
	r0, _, _ := syscall.Syscall6(
		uintptr(procSCardListReaders),
		4,
		uintptr(hContext),
		uintptr(unsafe.Pointer(mszGroups)),
		uintptr(unsafe.Pointer(mszReaders)),
		uintptr(unsafe.Pointer(pcchReaders)),
		0,
		0,
	)
	if r0 != 0 {
		retval = syscall.Errno(r0)
	}
	return
}

func SCardReleaseContext(hContext syscall.Handle) (retval error) {
	r0, _, _ := syscall.Syscall(
		uintptr(procSCardReleaseContext),
		1,
		uintptr(hContext),
		0,
		0,
	)
	if r0 != 0 {
		retval = syscall.Errno(r0)
	}
	return
}

func SCardEstablishContext(dwScope uint32, pvReserved1 uintptr, pvReserved2 uintptr, phContext *syscall.Handle) (retval error) {
	r0, _, _ := syscall.Syscall6(
		uintptr(procSCardEstablishContext),
		4,
		uintptr(dwScope),
		uintptr(pvReserved1),
		uintptr(pvReserved2),
		uintptr(unsafe.Pointer(phContext)),
		0,
		0,
	)
	if r0 != 0 {
		retval = syscall.Errno(r0)
	}
	return
}

func ReturnValue(err error) uint32 {
	rv, ok := err.(syscall.Errno)
	if !ok {
		rv = 0
	}
	return uint32(rv)
}

func UTF16ToStrings(ls []uint16) []string {
	var ss []string
	if len(ls) == 0 {
		return ss
	}
	if ls[len(ls)-1] != 0 {
		ls = append(ls, 0)
	}
	i := 0
	for j, cu := range ls {
		if cu == 0 {
			if j >= 1 && ls[j-1] == 0 {
				break
			}
			if j-i > 0 {
				ss = append(ss, string(utf16.Decode(ls[i:j])))
			}
			i = j + 1
			continue
		}
	}
	return ss
}

func main() {
	var (
		context  syscall.Handle
		scope    uint32
		groups   *uint16
		cReaders uint32
	)

	context = 0
	groups, err := syscall.UTF16PtrFromString(SCARD_ALL_READERS)
	if err != nil {
		fmt.Println("Reader Group: ", err)
		return
	}
	err = SCardListReaders(context, groups, nil, &cReaders)
	if err != nil {
		fmt.Printf("SCardListReaders: 0x%X %s\n", ReturnValue(err), err)
		return
	}
	r := make([]uint16, cReaders)
	err = SCardListReaders(context, groups, &r[0], &cReaders)
	if err != nil {
		fmt.Printf("SCardListReaders: 0x%X %s\n", ReturnValue(err), err)
		return
	}
	readers := UTF16ToStrings(r[:cReaders])
	fmt.Println("Readers:", len(readers), readers)

	scope = SCARD_SCOPE_SYSTEM
	err = SCardEstablishContext(scope, 0, 0, &context)
	if err != nil {
		fmt.Printf("SCardEstablishContext: 0x%X %s\n", ReturnValue(err), err)
		return
	}
	defer SCardReleaseContext(context)
	fmt.Printf("Context: %X\n", context)
}
