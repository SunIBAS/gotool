package compress

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
)

type ZipItem interface {
	FilePath() string
	WriteFile(writer *zip.Writer) error
}

type ZipStringItem struct {
	filePath string
	content  string
}

func NewZipStringItem(filepath, content string) ZipStringItem {
	return ZipStringItem{
		filePath: filepath,
		content:  content,
	}
}
func (zsi ZipStringItem) FilePath() string {
	return zsi.filePath
}
func (zsi ZipStringItem) WriteFile(writer *zip.Writer) error {
	// 在压缩包根目录添加 _filelist_.txt 文件
	filelist, err := writer.Create(zsi.filePath)
	if err != nil {
		return err
	}
	_, err = filelist.Write([]byte(zsi.content))
	if err != nil {
		return err
	}
	return nil
}

type ZipFileItem struct {
	absoluteFilePath string
	filePath         string
}

func NewZipFileItem(absoluteFilePath, filepath string) ZipFileItem {
	fileInfo, err := os.Stat(absoluteFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
		} else {
			fmt.Println("Error:", err)
		}
		return ZipFileItem{}
	}

	if fileInfo.IsDir() {
		return ZipFileItem{}
	} else {
		//fmt.Println(absoluteFilePath, "is not a directory")
		return ZipFileItem{
			absoluteFilePath: absoluteFilePath,
			filePath:         filepath,
		}
	}
}
func (zsi ZipFileItem) FilePath() string {
	return zsi.filePath
}
func (zsi ZipFileItem) WriteFile(writer *zip.Writer) error {
	if zsi.absoluteFilePath == "" {
		return errors.New(fmt.Sprintf("[%s] will be not file", zsi.absoluteFilePath))
	}
	// 在压缩包根目录添加 _filelist_.txt 文件
	f, err := os.Open(zsi.absoluteFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zsi.filePath
	header.Method = zip.Deflate

	w, err := writer.CreateHeader(header)
	//file, err := os.Open(zsi.absoluteFilePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	return nil
}
