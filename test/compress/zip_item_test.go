package compress

import (
	"github.com/SunIBAS/gotool/compress"
	"testing"
)

func TestZipItem(f *testing.T) {
	cp := compress.Zip{
		ZipFilePath: "D:\\all_code\\gotool\\test\\compress\\zip_item_test_TestZipItem.zip",
		ZipItem:     []compress.ZipItem{},
	}
	defer func() {
		if err := cp.Close(); err != nil {
			f.Error(err)
		}
	}()
	if err := cp.Init(); err != nil {
		f.Error(err)
	}

	cp.ZipItem = append(cp.ZipItem, compress.NewZipStringItem("name.txt", "name"))

	zipItem := compress.NewZipFileItem("D:\\all_code\\gotool\\test\\compress\\zip_item_test.go", "root/a/b/test.go")
	cp.ZipItem = append(cp.ZipItem, zipItem)

	if err := cp.Compress(); err != nil {
		f.Error(err)
	}
}
