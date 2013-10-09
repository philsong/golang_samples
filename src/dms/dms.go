package main

import (
	"fmt"
	"mif"
)

const (
	TRUE  = 1
	FALSE = 0
)

func main() {
	defer mif.FreeLib()
	hSerial := mif.HVD_OpenPort(2)

	defer mif.HVD_ClosePort(hSerial)

	fmt.Printf("hSerial=%06x\n", hSerial)
	var ret int
	if ret := mif.MIF_REQ_ANTI_SEL(hSerial); ret != TRUE {
		mif.MessageBox("Please place the card in card reader!")
		return
	}
	fmt.Println(ret)
	fmt.Println(mif.KEY_A, mif.KEY_B)
	var KEY_A mif.MIF_KEY
	var keyA string = "123456"
	//copy a string to array, need convert array to slice first
	copy(KEY_A.M_Value[:], []byte(keyA))
	if mif.MIF_AUTH_KEY(hSerial, mif.KEY_A, 0, &KEY_A) != TRUE {
		mif.MessageBox("Auth failed!")
		return
	}

	var BlockData mif.MIF_DATA_BLOCK
	if mif.MIF_Read(hSerial, 0, &BlockData) != TRUE {
		mif.MessageBox("MIF_Read failed!")
		return
	}

	fmt.Println(BlockData)

	copy(BlockData.M_Value[:], []byte("1111111111111111"))
	if mif.MIF_Write(hSerial, 0, &BlockData) != TRUE {
		mif.MessageBox("MIF_Write failed!")
		return
	}

	if mif.MIF_Read(hSerial, 0, &BlockData) != TRUE {
		mif.MessageBox("MIF_Read failed!")
		return
	}

	fmt.Println(BlockData)
	fmt.Printf("End hSerial=%06x\n", hSerial)
}
