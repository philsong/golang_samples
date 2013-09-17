package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	file, _ := os.Getwd()

	http.Handle("/", http.FileServer(http.Dir("./")))

	conn, err := net.Dial("udp", "baidu.com:80")
	//conn, err := net.Dial("tcp", "google.com:80")
	var ipaddr string
	if err != nil {
		// handle error
		ipaddr = "ip"
	} else {
		ipaddr = strings.TrimSpace(strings.Split(conn.LocalAddr().String(), ":")[0])
	}
	fmt.Printf("Directory [%s] can be accessed via http://%s:8123\n in internet browser, report bugs to philsong@techtrex.com", file, ipaddr)
	defer conn.Close()

	http.ListenAndServe(":8123", nil)
}
