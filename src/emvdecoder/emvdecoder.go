/*
package main

import (
	"fmt"
	//"strings"
)

func main() {
	//get input
	var emvdata string
	fmt.Scanf("%X", emvdata)
	//parser
	fmt.Println(emvdata)

	//put putput
}

*/
package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("emv decoder V0.01 by Philsong@techtrex.com", "\nPls send suggestion to me, thanks")
	fmt.Println("\n------------------------\nemv TAG support list in below:")

	fmt.Println("Card side:")
	fmt.Println("1:AUC")
	fmt.Println("2:TVR")
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

	for i, v := range tlvbytes {
		//fmt.Printf("BYTE[%d] base16 is 0x%02x\n", i, v)
		fmt.Printf("BYTE[%d] base2 is %08b\n", i, v)

		if item == 2 {
			if i == 0 {
				tvr_elements := make(map[int]string)

				tvr_elements[0] = "Offline data authentication was not performed"
				tvr_elements[1] = "SDA failed"
				tvr_elements[2] = "ICC data missing"
				tvr_elements[3] = "Card appears on terminal exception file"
				tvr_elements[4] = "DDA failed"
				tvr_elements[5] = "CDA failed"
				tvr_elements[6] = "RFU"
				tvr_elements[7] = "RFU"

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
		}
	}

	Input("press return to exit\n")
}

func Input(str string) string {
	print(str)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input
}
