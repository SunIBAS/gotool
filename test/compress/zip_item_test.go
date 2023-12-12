package compress

import (
	"archive/zip"
	"fmt"
	"github.com/google/uuid"
	"gotool/compress"
	"io"
	"os"
	"testing"
)

func TestZipItem(f *testing.T) {
	cp := compress.Zip{
		ZipFilePath: "D:\\all_code\\gotool\\test\\compress\\zip_item_test_TestZipItem.zip",
		ZipItem:     []compress.ZipItem{},
	}

	cp.ZipItem = append(cp.ZipItem, compress.NewZipStringItem("name.txt", "name"))

	zipItem := compress.NewZipFileItem(".\\zip_item_test.go", "root/a/b/test.go")
	cp.ZipItem = append(cp.ZipItem, zipItem)

	if err := cp.Compress(); err != nil {
		f.Error(err)
	}
}
func getId() string {
	id := uuid.New()
	ids := id.String()
	return ids
}
func TestZipAddFile(f *testing.T) {
	cp := compress.Zip{
		ZipFilePath: "D:\\all_code\\gotool\\test\\compress\\zip_item_test_TestZipAddFile.zip",
		ZipItem:     []compress.ZipItem{},
	}

	id := getId()
	fmt.Println(id)
	cp.ZipItem = append(cp.ZipItem, compress.NewZipStringItem("name.txt", id))

	zipItem := compress.NewZipFileItem("D:\\all_code\\gotool\\test\\compress\\zip_item_test.go", fmt.Sprintf("root/%s.go", id))
	cp.ZipItem = append(cp.ZipItem, zipItem)

	if err := cp.Compress(); err != nil {
		f.Error(err)
	}
}
func copyFilesToZip(source *zip.ReadCloser, dest *zip.Writer) error {
	for _, file := range source.File {
		sourceFile, err := file.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destFile, err := dest.Create(file.Name)
		if err != nil {
			return err
		}

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			return err
		}
	}
	return nil
}
func TestAdd(t *testing.T) {
	// 打开要读取的源ZIP文件
	sourceZipFilename := "D:\\all_code\\gotool\\test\\compress\\zip_item_test_TestZipAddFile.zip"
	sourceZipFile, err := zip.OpenReader(sourceZipFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sourceZipFile.Close()

	// 创建要写入的目标ZIP文件
	destZipFilename := "D:\\all_code\\gotool\\test\\compress\\destination.zip"
	destZipFile, err := os.Create(destZipFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer destZipFile.Close()

	// 创建一个zip writer，指向目标ZIP文件
	zipWriter := zip.NewWriter(destZipFile)
	defer zipWriter.Close()

	// 复制源ZIP文件中的内容到目标ZIP文件
	err = copyFilesToZip(sourceZipFile, zipWriter)
	if err != nil {
		fmt.Println(err)
		return
	}
	zipItem := compress.NewZipFileItem("D:\\all_code\\gotool\\test\\compress\\zip_item_test.go", fmt.Sprintf("root/haha.go"))
	zipItem.WriteFile(zipWriter)

	fmt.Println("内容已成功从一个ZIP文件复制到另一个ZIP文件")
}

func TestReadZip(t *testing.T) {
	sourceZipFilename := "Z:\\jiayu\\tmp\\zips\\1034-7bd695ae75d6aa186cec4cd87d8bd9eb.zip"
	r, err := zip.OpenReader(sourceZipFilename)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		fmt.Println(f.Name)
	}
}

func TestReadZipFiles(t *testing.T) {
	sourceZipFilename := "Z:\\jiayu\\tmp\\zips\\1034-7bd695ae75d6aa186cec4cd87d8bd9eb.zip"
	cp := compress.Zip{
		ZipFilePath: sourceZipFilename,
		ZipItem:     nil,
	}
	files := cp.GetFiles()
	for i := 0; i < len(files); i++ {
		if !files[i].FileInfo().IsDir() {
			fmt.Println(files[i].Name)
		}
	}

}
