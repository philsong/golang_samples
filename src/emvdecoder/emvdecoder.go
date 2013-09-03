package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	defer func() {
		str := recover()
		fmt.Println(str)
	}()

	//var input string
	//fmt.Scanln(&input)
	//godir()
	//gowalk()
	//readTerminalLog()

	fmt.Println("emv decoder V0.02 by Philsong@techtrex.com", "\nPls send suggestion to me, thanks")
	fmt.Println("\n------------------------\nemv TAG support list in below:")

	fmt.Println("Card side:")
	fmt.Println("1:AUC(support bits analysis)")
	fmt.Println("2:TVR(support bits analysis)")

	fmt.Println("3:TSI")
	fmt.Println("4:CVR")
	fmt.Println("5:AIP")

	fmt.Println("Terminal side:")
	fmt.Println("6:TC")
	fmt.Println("7:ATC")
	fmt.Println("8:CVM")
	fmt.Println("9:Issuer Script Results")
	fmt.Println("10:Authorisation Response Code")

	value := Input("\nPlease enter a EMV index in above: ")

	trimmed := strings.TrimSpace(value)
	//fmt.Printf("item %s!\n", trimmed)
	item, _ := strconv.Atoi(trimmed)

	var tlvvalue string
	switch {
	case item == 1:
		tlvvalue = Input("Please enter a AUC value: ")
	case item == 2:
		tlvvalue = Input("Please enter a TVR value: ")
	case item < 11:
		tlvvalue = Input("Please enter a emv value: ")
	default:
		fmt.Println("Please enter a valid index!")
		return
	}

	tlvdata := strings.TrimSpace(tlvvalue)
	tlvdata = strings.Replace(tlvdata, " ", "", -1)
	fmt.Printf("tlvdata str value[%s]\n", tlvdata)

	tlvbytes, _ := hex.DecodeString(tlvdata)
	fmt.Println("tlvbytes mem value", tlvbytes)
	fmt.Printf("tlvbytes hex value 0x%08x\n", tlvbytes)

	printTLV(item, tlvbytes)

	//Input("press return to exit\n")
}

func godir() {
	dir, err := os.Open(".")
	if err != nil {
		return
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return
	}
	for _, fi := range fileInfos {
		fmt.Println(fi.Name())
	}
}

func gowalk() {
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
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
	fmt.Println(str)
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

func Input(str string) string {
	print(str)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input
}

func initAUC() [2][8]string {
	var elements [2][8]string

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
func initTVR() [5][8]string {
	var tvr_elements [5][8]string

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

func printElement(tvr_elements [8]string, v uint8) {
	for j := 0; j < 8; j++ {
		var shiftNum uint32 = uint32(7 - j)
		//fmt.Printf("shift[%d]\n", shiftNum)
		var mask uint8 = 0x01 << shiftNum
		//fmt.Println("mask", mask)
		fmt.Print(j, ":  ")
		if v&mask == mask {
			fmt.Println(tvr_elements[j])
		} else {
			fmt.Println("--------------")
		}
	}
}

func printTLV(item int, tlvbytes []byte) {
	fmt.Println("-------------------------我是分割线--------------------------")

	tlvNotes := make([][8]string, 5)
	switch item {
	case 1:
		if len(tlvbytes) != 2 {
			//fmt.Println("wrong data, must be 2 bytes\n\n")
			panic("wrong data, must be 2 bytes\n\n")
			//return
		}
		var auc [2][8]string = initAUC()
		//tlvNotes
		tlvNotes = auc[:]
	case 2:
		if len(tlvbytes) != 5 {
			fmt.Println("wrong data, must be 5 bytes\n\n")
			return
		}
		var tvr [5][8]string = initTVR()
		tlvNotes = tvr[:]
	}

	for i, v := range tlvbytes {
		//fmt.Printf("BYTE[%d] base16 is 0x%02x\n", i, v)
		fmt.Printf("BYTE[%d] %08b\n", i+1, v)

		printElement(tlvNotes[i], v)
	}

	fmt.Println("-------------------------the end---------------------------")
}
