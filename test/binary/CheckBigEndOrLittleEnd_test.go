package binary

import (
	"encoding/binary"
	"fmt"
	"github.com/SunIBAS/gotool/Datas"
	"github.com/SunIBAS/gotool/FileDirUtils"
	"os"
	"testing"
)

func TestToCreateBigLittleEndFile(t *testing.T) {
	bigFile := "D:\\all_code\\gotool\\test\\binary\\bigFile.txt"
	littleFile := "D:\\all_code\\gotool\\test\\binary\\littleFile.txt"
	content := "ABCD"
	FileDirUtils.CreateBigEndFile(bigFile, Datas.StringToUint16(content))
	FileDirUtils.CreateLittleEndFile(littleFile, Datas.StringToUint16(content))
}

func TestBigLittle(t *testing.T) {
	//bigFile := "D:\\all_code\\gotool\\test\\binary\\bigFile.txt"
	littleFile := "D:\\all_code\\gotool\\test\\binary\\littleFile.txt"
	bifF, _ := os.Open(littleFile)
	defer bifF.Close()
	var byteOrder uint16
	err := binary.Read(bifF, binary.BigEndian, &byteOrder)
	const (
		littleEndian = 0x4949
		bigEndian    = 0x4D4D
	)
	switch byteOrder {
	case littleEndian:
		fmt.Println("little end")
	case bigEndian:
		fmt.Println("big end")
	default:
		fmt.Println("undefined")
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("byteOrder:[%v]\n", byteOrder))
}
