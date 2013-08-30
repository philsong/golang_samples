package main

import (
	"github.com/lxn/win"
	"strconv"
	"syscall"
)

func _TEXT(_str string) *uint16 {
	return syscall.StringToUTF16Ptr(_str)
}
func _toString(_n int32) string {
	return strconv.Itoa(int(_n))
}
func main() {
	var hwnd win.HWND
	cxScreen := win.GetSystemMetrics(win.SM_CXSCREEN)
	cyScreen := win.GetSystemMetrics(win.SM_CYSCREEN)
	win.MessageBox(hwnd, _TEXT("屏幕长-:"+_toString(cxScreen)+"宽:"+_toString(cyScreen)), _TEXT(" 消息"), win.MB_OK)
}
