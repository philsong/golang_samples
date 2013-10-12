/* example 1 start process
package main

import "syscall"
import "unsafe"

func main() {
	var hand uintptr = uintptr(0)
	var operator uintptr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("open")))
	var fpath uintptr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:\\WINDOWS\\NOTEPAD.EXE")))
	var param uintptr = uintptr(0)
	var dirpath uintptr = uintptr(0)
	var ncmd uintptr = uintptr(1)
	shell32 := syscall.NewLazyDLL("shell32.dll")
	ShellExecuteW := shell32.NewProc("ShellExecuteW")
	_, _, _ = ShellExecuteW.Call(hand, operator, fpath, param, dirpath, ncmd)
}
*/

//example 2:守护进程
package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	lf, err := os.OpenFile("angel.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		os.Exit(1)
	}
	defer lf.Close()

	// 日志
	l := log.New(lf, "", os.O_APPEND)

	for {
		cmd := exec.Command("C:\\WINDOWS\\NOTEPAD.EXE", "C:\\WINDOWS\\ODBC.INI")
		err := cmd.Start()
		if err != nil {
			l.Printf("%s 启动命令失败", time.Now().Format("2006-01-02 15:04:05"), err)

			time.Sleep(time.Second * 5)
			continue
		}
		l.Printf("%s 进程启动", time.Now().Format("2006-01-02 15:04:05"), err)
		err = cmd.Wait()
		l.Printf("%s 进程退出", time.Now().Format("2006-01-02 15:04:05"), err)

		time.Sleep(time.Second * 1)
	}
}

/*
这里还有一个shell实现的. 记得给予执行权限哦,chmod +x you_command
#! /bin/bash
while true; do
    ./you_command
done
*/
