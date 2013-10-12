package mif

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	TRUE  = 1
	FALSE = 0
)

func TestMif(t *testing.T) {
	defer FreeLib()
	hSerial := HVD_OpenPort(3)

	defer HVD_ClosePort(hSerial)

	if ret := MIF_REQ_ANTI_SEL(hSerial); ret != TRUE {
		MessageBox("Please place the card in card reader!")
		return
	}

	fmt.Println(KEY_A, KEY_B)
	var stKEY_B MIF_KEY
	bytekeyB, _ := hex.DecodeString("85A438F72A8A")
	copy(stKEY_B.M_Value[:], bytekeyB)
	if MIF_AUTH_KEY(hSerial, KEY_B, 0, &stKEY_B) != TRUE {
		MessageBox("Auth failed!")
		return
	}

	var BlockData MIF_DATA_BLOCK
	if MIF_Read(hSerial, 1, &BlockData) != TRUE {
		MessageBox("MIF_Read failed!")
		return
	}

	fmt.Println(BlockData)

	byteBlockData, _ := hex.DecodeString("11111111111111111111111111111111")
	fmt.Println(byteBlockData)
	copy(BlockData.M_Value[:], byteBlockData)
	if MIF_Write(hSerial, 1, &BlockData) != TRUE {
		MessageBox("MIF_Write failed!")
		return
	}

	if MIF_Read(hSerial, 1, &BlockData) != TRUE {
		MessageBox("MIF_Read failed!")
		return
	}

	fmt.Println(BlockData)
	fmt.Println("The End")
}
