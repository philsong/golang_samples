package dmslogic

import (
	"encoding/hex"
	"fmt"
	"mif"
)

const (
	TRUE  = 1
	FALSE = 0
)

func mif_write(handle uintptr, blocknumber byte, blockdata string) {
	var BlockData mif.MIF_DATA_BLOCK
	byteblockdata, _ := hex.DecodeString(blockdata)
	copy(BlockData.M_Value[:], byteblockdata)
	if mif.MIF_Write(handle, blocknumber, &BlockData) != TRUE {
		mif.MessageBox("mif_write failed!")
		return
	}
}
func GenDepKeyCard(block1, block2 string) {
	defer mif.FreeLib()
	hSerial := mif.HVD_OpenPort(3)

	defer mif.HVD_ClosePort(hSerial)

	if ret := mif.MIF_REQ_ANTI_SEL(hSerial); ret != TRUE {
		mif.MessageBox("Please place the card in card reader!")
		return
	}

	fmt.Println(mif.KEY_A, mif.KEY_B)
	var KEY_B mif.MIF_KEY
	keyB, _ := hex.DecodeString("85A438F72A8A")
	copy(KEY_B.M_Value[:], keyB)
	if mif.MIF_AUTH_KEY(hSerial, mif.KEY_B, 0, &KEY_B) != TRUE {
		mif.MessageBox("Auth failed!")
		return
	}

	// write block 1 data
	mif_write(hSerial, 1, block1)
	// write block 2 data
	mif_write(hSerial, 2, block2)
	fmt.Println("The End")
}
