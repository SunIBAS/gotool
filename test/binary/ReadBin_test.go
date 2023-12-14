package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/SunIBAS/gotool/Datas"
	"io"
	"os"
	"testing"
)

var binFile = "D:\\all_code\\gotool\\test\\binary\\test.bin"
var tifFile = "C:\\Users\\11340\\Documents\\paper\\中期\\zq\\附件2：培养环节考核要求\\附件2-1.研究生开题、中期登记表和报告（模板）\\中文版\\图\\数据\\ndvi.2020.20.tif"

type fieldType uint16
type Tag uint16
type iFDEntry struct {
	// Bytes 0-1
	//
	// The Tag that identifies the field.
	Tag Tag

	// Bytes 2-3
	//
	//The field FType.
	FType fieldType

	//  Bytes 4-7
	//
	// The number of values, Count of the indicated Type.
	Count uint32

	// Bytes 8-11
	//
	// The Value Offset, the file offset (in bytes) of the Value for
	// the field. The Value is expected to begin on a word boundary; the
	// corresponding Value Offset will thus be an even number. This file offset
	// may point anywhere in the file, even after the image data.
	ValueOffset uint32
}

func TestReadBin(t *testing.T) {
	f, _ := os.Open(binFile)
	bo := binary.LittleEndian
	var ifd1 iFDEntry
	//f.Seek(2, io.SeekStart)
	var numDirectoryEntries uint16
	binary.Read(f, bo, &numDirectoryEntries)
	fmt.Println(numDirectoryEntries)
	if err := binary.Read(f, bo, &ifd1); err != nil {
		panic(err)
	}
	fmt.Println(ifd1)

	var ifd = iFDEntry{
		Tag:         256,
		FType:       3,
		Count:       1,
		ValueOffset: 2745,
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ifd)
	littleEndianBytes := buf.Bytes()
	fmt.Printf("Little-endian bytes: % x\n", littleEndianBytes)
}

func TestReadBin2(t *testing.T) {
	f, _ := os.Open(tifFile)
	iFDOffset := 31812236
	if _, err := f.Seek(int64(iFDOffset), io.SeekStart); err != nil {
		t.Error("error: unable to read IFD Start")
	}
	var numDirectoryEntries uint16
	bo := binary.LittleEndian
	if err := binary.Read(f, bo, &numDirectoryEntries); err != nil {
		t.Error("error: unable to read directory entry")
	}
	fmt.Println(numDirectoryEntries)
	for i := 0; i < int(numDirectoryEntries); i++ {
		var ifd1 iFDEntry
		//f.Seek(2, io.SeekStart)
		if err := binary.Read(f, bo, &ifd1); err != nil {
			panic(err)
		}
		fmt.Println(ifd1)

		_, littleEndianBytes := Datas.Struct2Bytes(ifd1, binary.LittleEndian)
		fmt.Printf("Little-endian bytes: %s\n", littleEndianBytes)
	}
}
