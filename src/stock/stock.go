/*
 *@author 菠菜
 *@Version 0.4
 *@Update time 2013-11-14
 *@golang 微信公众平台API引擎开发 Plugin for Stock
 *@演示 微信订阅号	gostock
 *@开源 https://github.com/philsong/
 */
package stock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	//	"strings"
)

func Deal(keyword string) (resp string) {
	//url := "http://qt.gtimg.cn/q=sh000001,sz399001"
	url := "http://qt.gtimg.cn/q=sh000001"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if response.Status != "200 OK" {
		fmt.Println(response.Status)
		return
	}

	b, _ := httputil.DumpResponse(response, false)
	fmt.Println("dump")
	fmt.Print(string(b))

	data, _ := ioutil.ReadAll(response.Body)

	s := string(data)
	fmt.Println("ReadAll")
	fmt.Printf("%s", s)
	fmt.Println("ReadAll end")

	fmt.Println("Decode start")
	//todo, who can help in here?
	fmt.Println("Decode end")

	return s
}
