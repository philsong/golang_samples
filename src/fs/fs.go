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

	conn, _ := net.Dial("udp", "baidu.com:80")
	defer conn.Close()
	fmt.Printf("others can access your directly %s via open http://%s:8123\n in internet browser, report bugs to philsong@techtrex.com", file, strings.TrimSpace(strings.Split(conn.LocalAddr().String(), ":")[0]))
	http.ListenAndServe(":8123", nil)
}
