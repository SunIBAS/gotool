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
type Zip struct {
	ZipFilePath string
	ZipItem     []ZipItem
	zipFile     *os.File
	zipWriter   *zip.Writer
}

func (z *Zip) Close() error {
	if err := z.zipWriter.Close(); err != nil {
		return errors.New(fmt.Sprintf("[close] zipWriter {%v}", err))
	}
	if err := z.zipFile.Close(); err != nil {
		return errors.New(fmt.Sprintf("[close] zipFile {%v}", err))
	}
	return nil
}
func (z *Zip) Init() error {
	var err error
	z.zipFile, err = os.Create(z.ZipFilePath)
	if err != nil {
		return err
	}
	z.zipWriter = zip.NewWriter(z.zipFile)
	return nil
}
func (z *Zip) Compress() error {
	for i := 0; i < len(z.ZipItem); i++ {
		if err := z.ZipItem[i].WriteFile(z.zipWriter); err != nil {
			return err
		}
	}
	return nil
}

type zipStringItem struct {
	filePath string
	content  string
}

func NewZipStringItem(filepath, content string) zipStringItem {
	return zipStringItem{
		filePath: filepath,
		content:  content,
	}
}
func (zsi zipStringItem) FilePath() string {
	return zsi.filePath
}
func (zsi zipStringItem) WriteFile(writer *zip.Writer) error {
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

type zipFileItem struct {
	absoluteFilePath string
	filePath         string
}

func NewZipFileItem(absoluteFilePath, filepath string) zipFileItem {
	fileInfo, err := os.Stat(absoluteFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
		} else {
			fmt.Println("Error:", err)
		}
		return zipFileItem{}
	}

	if fileInfo.IsDir() {
		return zipFileItem{}
	} else {
		//fmt.Println(absoluteFilePath, "is not a directory")
		return zipFileItem{
			absoluteFilePath: absoluteFilePath,
			filePath:         filepath,
		}
	}
}
func (zsi zipFileItem) FilePath() string {
	return zsi.filePath
}
func (zsi zipFileItem) WriteFile(writer *zip.Writer) error {
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
