package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getremote() {
	//resp, err := http.Get("http://pay.pingliwang.com:9090/")
	resp, err := http.Get("http://192.168.7.100:8123/CRT%20log")

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%s", resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("%s", body)

}

func main() {
	getremote()
}
