package Datas

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Struct2Bytes(iStruct interface{}, bigOrLittle binary.ByteOrder) ([]byte, string) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, bigOrLittle, iStruct)
	if err != nil {
		return nil, ""
	}
	bOLEndianBytes := buf.Bytes()
	return bOLEndianBytes, fmt.Sprintf("%x", bOLEndianBytes)
}
