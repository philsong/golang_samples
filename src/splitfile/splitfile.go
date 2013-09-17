// splitfile.go
package main

import (
	//	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

import "strconv"

var infile *string = flag.String("f", "Null", "please input a file name or dir.")
var size *string = flag.String("s", "0", "please input a dst file size.")

//func dataCopy(file *os.File, file *os.File)

func SplitFile(file *os.File, size int) {
	finfo, err := file.Stat()
	if err != nil {
		fmt.Println("get file info failed:", file, size)
	}

	fmt.Println(finfo, size)

	//每次最多拷贝1m
	bufsize := 1024 * 1024
	if size < bufsize {
		bufsize = size
	}

	buf := make([]byte, bufsize)

	num := (int(finfo.Size()) + size - 1) / size
	fmt.Println(num, len(buf))

	for i := 0; i < num; i++ {
		copylen := 0
		newfilename := finfo.Name() + strconv.Itoa(i)
		newfile, err1 := os.Create(newfilename)
		if err1 != nil {
			fmt.Println("failed to create file", newfilename)
		} else {
			fmt.Println("create file:", newfilename)
		}

		for copylen < size {
			n, err2 := file.Read(buf)
			if err2 != nil && err2 != io.EOF {
				fmt.Println(err2, "failed to read from:", file)
				break
			}

			if n <= 0 {
				break
			}

			//fmt.Println(n, len(buf))

			//写文件
			w_buf := buf[:n]
			newfile.Write(w_buf)
			copylen += n
		}
	}

	return
}

func main() {
	flag.Parse()

	if *infile == "Null" {
		fmt.Println("no file to input")
		return
	}

	file, err := os.Open(*infile)
	if err != nil {
		fmt.Println("failed to open:", *infile)
	}

	defer file.Close()

	size, _ := strconv.Atoi(*size)

	SplitFile(file, size*1024)

}
