package FileDirUtils

import (
	"encoding/binary"
	"fmt"
	"os"
)

type FileType = int

const (
	BigEndFile    FileType = 1
	LittleEndFile FileType = 2
)

func createFile(fileName string, content []uint16, filetype FileType) error {
	//fileName := "big_endian_file.bin"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// 创建大端字节序的数据
	//data := uint16(0xABCD)
	bytes := make([]byte, len(content)*2, len(content)*2)
	for i := 0; i < len(content); i++ {
		bs := make([]byte, 2)
		if filetype == BigEndFile {
			binary.BigEndian.PutUint16(bs, content[i])
		} else {
			binary.LittleEndian.PutUint16(bs, content[i])
		}
		bytes = append(bytes, bs...)
	}

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	//fmt.Println("Big-endian file created:", fileName)
	return nil
}

func CreateBigEndFile(fileName string, content []uint16) error {
	return createFile(fileName, content, BigEndFile)
}

func CreateLittleEndFile(fileName string, content []uint16) error {
	return createFile(fileName, content, LittleEndFile)
}
