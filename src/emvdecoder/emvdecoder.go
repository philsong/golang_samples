/*
File      : emvdecoder.go
Author    : Philsong
E-Mail    : 78623269@qq.com
*/
package main

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	//"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type Prompt struct {
	Index      int
	Tip1st     string
	Tip2nd     string
	tlvdatalen int
	bits       [][8]string
}

type Emvresult struct {
	Index  string
	Data   string
	Output string
}

var emvdecoder = [10]Prompt{
	{1, ("9F07 - AUC"), "Please enter a AUC value: ", 2, make([][8]string, 2)},
	{2, ("95 - TVR"), "Please enter a TVR value: ", 5, make([][8]string, 5)},
	{3, ("9B - TSI"), "Please enter a TSI value: ", 2, make([][8]string, 2)},
	{4, ("9F34 - CVR(TBD)"), "Please enter a CVR value: ", 2, make([][8]string, 2)},
	{5, ("82 - AIP"), "Please enter a AIP value: ", 2, make([][8]string, 2)},

	{6, ("9F33 - TC"), "Please enter a TC value: ", 3, make([][8]string, 3)},
	{7, ("9F40 - ATC"), "Please enter a ATC value: ", 5, make([][8]string, 5)},
	{8, ("CVM(TBD)"), "Please enter a CVM value: ", 3, make([][8]string, 3)},
	{9, ("Issuer Script Results(TBD)"), "Please enter a Issuer Script Results value: ", 5, make([][8]string, 5)},
	{10, ("Authorisation Response Code(TBD)"), "Please enter a Authorisation Response Code value: ", 2, make([][8]string, 2)},
}

func checkError(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		debug.PrintStack()
		//os.Exit(1)
	}
}

func parseHandle(w http.ResponseWriter, r *http.Request) {
	println("Request ", r.URL.Path, " from ", r.RemoteAddr)
	r.ParseForm()

	fmt.Println("idx:", r.Form["idx"])
	fmt.Println("data:", r.Form["data"])

	var emvresult Emvresult
	var output string
	var tlvdata string
	trimmed := strings.TrimSpace(r.FormValue("idx"))
	item, _ := strconv.Atoi(trimmed)

	if item < 1 || item > len(emvdecoder) {
		output = "Please select a valid item!"

		emvresult = Emvresult{"invalid select", "0x" + tlvdata, output}
	} else {
		tlvdata = strings.TrimSpace(r.FormValue("data"))
		tlvdata = strings.Replace(tlvdata, " ", "", -1)

		if len(tlvdata) != emvdecoder[item-1].tlvdatalen*2 {
			output = fmt.Sprintf("%s 's length should be %d bytes\n%s\n\n",
				emvdecoder[item-1].Tip1st, emvdecoder[item-1].tlvdatalen, emvdecoder[item-1].Tip2nd)
		} else {
			output = parseEMV(item, tlvdata)
		}

		emvresult = Emvresult{emvdecoder[item-1].Tip1st, "0x" + tlvdata, output}
	}

	t, err := template.ParseFiles("parse.html")
	checkError(w, err)

	//fmt.Print(emvresult)
	err = t.Execute(w, emvresult)
	checkError(w, err)
}

var indexTemplate = template.Must(template.ParseFiles("index.html"))
var uploadTemplate = template.Must(template.ParseFiles("upload.html"))

func indexHandle(w http.ResponseWriter, r *http.Request) {
	println("Request ", r.URL.Path, " from ", r.RemoteAddr)
	//   path := r.URL.Path[1:]
	path := "." + r.URL.Path

	if path == "./favicon.ico" {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()

	if err := indexTemplate.Execute(w, emvdecoder); err != nil {
		//	log.Fatal("Execute: ", err.Error())
		checkError(w, err)
		return
	}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		if err := uploadTemplate.Execute(w, nil); err != nil {
			//		log.Fatal("Execute: ", err.Error())
			checkError(w, err)
			return
		}
	} else {
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Fatal("FormFile: ", err.Error())
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Fatal("Close: ", err.Error())
				return
			}
		}()

		const layout = "2013-Feb-03"
		fileremote := "uploaddir/" + header.Filename + "-" + time.Now().Format(layout)

		for {
			f, err := os.Open(fileremote)
			if err != nil && os.IsNotExist(err) {
				fmt.Printf("%s file does not exist!\n", fileremote)

				f, err = os.Create(fileremote)
				defer f.Close()

				io.Copy(f, file)
				/*
					bytes, err := ioutil.ReadAll(file)
					if err != nil {
						log.Fatal("ReadAll: ", err.Error())
						return
					}

					w.Write(bytes)
				*/
				break
			} else {
				fmt.Printf("%s file exist!\n", fileremote)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				fmt.Println(r.Intn(10000))
				fileremote += strconv.Itoa(r.Intn(10000))
			}
		}
	}
}

type TraceHandler struct {
	h http.Handler
	n int
}

func (r *TraceHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.n++
	fmt.Printf("counter = %d\n", r.n) //why counter always zero
	fmt.Println("get", req.URL.Path, " from ", req.RemoteAddr)
	r.h.ServeHTTP(w, req)
}

func main() {
	port := "9090" //Default port
	if len(os.Args) > 1 {
		port = strings.Join(os.Args[1:2], "")
	}

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	h := http.StripPrefix("/icclogs/", http.FileServer(http.Dir("./logs/")))
	http.Handle("/icclogs/", &TraceHandler{h: h, n: 0})

	http.HandleFunc("/", indexHandle)      //设置访问的路由
	http.HandleFunc("/parse", parseHandle) //设置访问的路由
	http.HandleFunc("/upload", uploadHandle)

	println("Listening on port ", port, "...")
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parseEMV(item int, tlvdata string) string {
	//readTerminalLog()
	var output string

	tlvbytes, _ := hex.DecodeString(tlvdata)

	fmt.Printf("-------------------------start--------------------------\n")

	switch item {
	case 1:
		emvdecoder[item-1].bits = initAUC()
	case 2:
		emvdecoder[item-1].bits = initTVR()
	case 3:
		emvdecoder[item-1].bits = initTSI()
	case 4:
		emvdecoder[item-1].bits = initCVR()
	case 5:
		emvdecoder[item-1].bits = initAIP()
	case 6:
		emvdecoder[item-1].bits = initTC()
	case 7:
		emvdecoder[item-1].bits = initATC()
	case 8:
		emvdecoder[item-1].bits = initCVR()
	case 9:
		emvdecoder[item-1].bits = initCVR()
	case 10:
		emvdecoder[item-1].bits = initCVR()
	default:
		return ""
	}

	for i, v := range tlvbytes {
		//fmt.Printf("BYTE[%d] base16 is 0x%02x\n", i, v)
		output += fmt.Sprintf("BYTE[%d] %08b\n", i+1, v)

		output += printElement(emvdecoder[item-1].bits[i], v)
	}

	fmt.Printf("-------------------------the end--------------------------\n")

	return output
}

func printElement(tvr_elements [8]string, v uint8) string {
	var output string
	for j := 0; j < 8; j++ {
		var shiftNum uint32 = uint32(7 - j)
		//fmt.Printf("shift[%d]\n", shiftNum)
		var mask uint8 = 0x01 << shiftNum
		//fmt.Fprintf(w, "mask", mask)
		output += fmt.Sprintf("%d: ", 8-j)
		if v&mask == mask {
			output += "[***]"
		} else {
			output += "[   ]"
		}
		output += tvr_elements[j]
		output += "\n"
	}

	return output
}

func readTerminalLog() {
	file, err := os.Open("test 01.txt")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return
	}
	str := string(bs)
	writeTerminalLog(str)
	fmt.Printf(str)
}

func writeTerminalLog(str string) {
	file, err := os.Create("test.txt")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()
	file.WriteString(str)
}

func initAUC() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "Valid for domestic cash transactions "
	elements[index][1] = "Valid for international cash transactions "
	elements[index][2] = "Valid for domestic goods "
	elements[index][3] = "Valid for international goods "
	elements[index][4] = "Valid for domestic services "
	elements[index][5] = "Valid for international services "
	elements[index][6] = "Valid at ATMs "
	elements[index][7] = "Valid at terminals other than ATMs "

	index = 1
	elements[index][0] = "Domestic cashback allowed "
	elements[index][1] = "International cashback allowed "
	elements[index][2] = "Application not yet effective"
	elements[index][3] = "Requested service not allowed for card product "
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

//TVR data
func initTVR() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 5)

	index := 0
	tvr_elements[index][0] = "Offline data authentication was not performed"
	tvr_elements[index][1] = "SDA failed"
	tvr_elements[index][2] = "ICC data missing"
	tvr_elements[index][3] = "Card appears on terminal exception file"
	tvr_elements[index][4] = "DDA failed"
	tvr_elements[index][5] = "CDA failed"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 1
	tvr_elements[index][0] = "ICC and terminal have different application versions"
	tvr_elements[index][1] = "Expired application "
	tvr_elements[index][2] = "Application not yet effective"
	tvr_elements[index][3] = "Requested service not allowed for card product "
	tvr_elements[index][4] = "New card "
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "Cardholder verification was not "
	tvr_elements[index][1] = "Unrecognised CVM "
	tvr_elements[index][2] = "PIN Try Limit exceeded "
	tvr_elements[index][3] = "PIN entry required and PIN pad "
	tvr_elements[index][4] = "not present or not working "
	tvr_elements[index][5] = "Online PIN entered "
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 3
	tvr_elements[index][0] = "Transaction exceeds floor limit "
	tvr_elements[index][1] = "Lower consecutive offline limit"
	tvr_elements[index][2] = "Upper consecutive offline limit "
	tvr_elements[index][3] = "Transaction selected randomly for online processing "
	tvr_elements[index][4] = "Merchant forced transaction online"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 4
	tvr_elements[index][0] = "Default TDOL used "
	tvr_elements[index][1] = "Issuer authentication failed "
	tvr_elements[index][2] = "Script processing failed before final GENERATE AC"
	tvr_elements[index][3] = "Script processing failed after final GENERATE AC"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	return tvr_elements
}

func initTSI() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "Offline data authentication was performed  "
	elements[index][1] = "Cardholder verification was performed "
	elements[index][2] = "Card risk management was performed  "
	elements[index][3] = "Issuer authentication was performed "
	elements[index][4] = "Terminal risk management was performed "
	elements[index][5] = "Script processing was performed "
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

func initCVR() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

/*
Fail cardholder verification if this CVM is unsuccessful
Apply succeeding CV Rule if this CVM is unsuccessful
Fail CVM processing
Plaintext PIN verification performed by ICC
Enciphered PIN verified online
Plaintext PIN verification
performed by ICC and signature (paper)
Enciphered PIN verification performed by ICC
Enciphered PIN verification performed by ICC and signature (paper)
Values in the range 000110-011101 reserved for future use by this specification
Signature (paper)
No CVM required
Values in the range 100000-101111 reserved for use by the individual payment systems
Values in the range 110000-111110 reserved for use by the issuer
This value is not available for use
*/
func initAIP() [][8]string {
	var elements [][8]string = make([][8]string, 2)

	index := 0
	elements[index][0] = "RFU"
	elements[index][1] = "SDA supported "
	elements[index][2] = "DDA supported"
	elements[index][3] = "Cardholder verification is supported "
	elements[index][4] = "Terminal risk management is to be performed "
	elements[index][5] = "Issuer authentication is supported"
	elements[index][6] = "RFU"
	elements[index][7] = "CDA supported "

	index = 1
	elements[index][0] = "RFU"
	elements[index][1] = "RFU"
	elements[index][2] = "RFU"
	elements[index][3] = "RFU"
	elements[index][4] = "RFU"
	elements[index][5] = "RFU"
	elements[index][6] = "RFU"
	elements[index][7] = "RFU"

	return elements
}

func initTC() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 3)

	index := 0
	tvr_elements[index][0] = "Manual key entry "
	tvr_elements[index][1] = "Magnetic stripe "
	tvr_elements[index][2] = "IC with contacts"
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 1
	tvr_elements[index][0] = "Plaintext PIN for ICC verification"
	tvr_elements[index][1] = "Enciphered PIN for online verification   "
	tvr_elements[index][2] = "Signature (paper)"
	tvr_elements[index][3] = "Enciphered PIN for offline verification "
	tvr_elements[index][4] = "No CVM Required  "
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "SDA"
	tvr_elements[index][1] = "DDA"
	tvr_elements[index][2] = "Card capture  "
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "CDA"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	return tvr_elements
}

func initATC() [][8]string {
	var tvr_elements [][8]string = make([][8]string, 5)

	index := 0
	tvr_elements[index][0] = "Cash"
	tvr_elements[index][1] = "Goods"
	tvr_elements[index][2] = "Services"
	tvr_elements[index][3] = "Cashback"
	tvr_elements[index][4] = "Inquiry"
	tvr_elements[index][5] = "Transfer"
	tvr_elements[index][6] = "Payment"
	tvr_elements[index][7] = "Administrative"

	index = 1
	tvr_elements[index][0] = "Cash Deposit"
	tvr_elements[index][1] = "RFU"
	tvr_elements[index][2] = "RFU"
	tvr_elements[index][3] = "RFU"
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 2
	tvr_elements[index][0] = "Numeric keys "
	tvr_elements[index][1] = "Alphabetic and special characters keys  "
	tvr_elements[index][2] = "Command keys  "
	tvr_elements[index][3] = "Function keys  "
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "RFU"
	tvr_elements[index][7] = "RFU"

	index = 3
	tvr_elements[index][0] = "Print, attendant "
	tvr_elements[index][1] = "Print, cardholder "
	tvr_elements[index][2] = "Display, attendant "
	tvr_elements[index][3] = "Display, cardholder "
	tvr_elements[index][4] = "RFU"
	tvr_elements[index][5] = "RFU"
	tvr_elements[index][6] = "Code table 10"
	tvr_elements[index][7] = "Code table 9"

	index = 4
	tvr_elements[index][0] = "Code table 8"
	tvr_elements[index][1] = "Code table 7"
	tvr_elements[index][2] = "Code table 6"
	tvr_elements[index][3] = "Code table 5"
	tvr_elements[index][4] = "Code table 4"
	tvr_elements[index][5] = "Code table 3"
	tvr_elements[index][6] = "Code table 2"
	tvr_elements[index][7] = "Code table 1"

	return tvr_elements
}
