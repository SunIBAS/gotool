package compress

import (
	"archive/zip"
	"errors"
	"fmt"
	"gotool/FileDirUtils"
	"io"
	"os"
	"path/filepath"
)

type Zip struct {
	ZipFilePath string
	ZipItem     []ZipItem
	zipFile     *os.File
	zipWriter   *zip.Writer
}

// Close /////////////  Write zip file  ////////////////
func (z *Zip) Close() error {
	if err := z.zipWriter.Close(); err != nil {
		return errors.New(fmt.Sprintf("[close] zipWriter {%v}", err))
	}
	if err := z.zipFile.Close(); err != nil {
		return errors.New(fmt.Sprintf("[close] zipFile {%v}", err))
	}
	return nil
}
func (z *Zip) initForCompress() error {
	var err error
	if FileDirUtils.FileIsExist(z.ZipFilePath) {
		//z.zipFile, err = os.Open(z.ZipFilePath)
		err = z.copyForNewWriter(z.ZipFilePath)
	} else {
		z.zipFile, err = os.Create(z.ZipFilePath)
		z.zipWriter = zip.NewWriter(z.zipFile)
	}
	if err != nil {
		return err
	}
	return nil
}
func (z *Zip) copyForNewWriter(ZipFilePath string) error {
	f := filepath.Base(ZipFilePath)
	newF := ZipFilePath[:len(ZipFilePath)-len(f)] + "~" + f
	err := os.Rename(ZipFilePath, newF)
	if err != nil {
		return err
	}
	sourceZipFile, err := zip.OpenReader(newF)
	if err != nil {
		return err
	}
	destZipFile, err := os.Create(f)
	if err != nil {
		return err
	}
	// 创建一个zip writer，指向目标ZIP文件
	zipWriter := zip.NewWriter(destZipFile)
	err = copyFilesToZip(sourceZipFile, zipWriter)
	if err != nil {
		return err
	}
	z.zipWriter = zipWriter
	z.zipFile = destZipFile
	if err := sourceZipFile.Close(); err != nil {
		return err
	}
	if err := os.Remove(newF); err != nil {
		fmt.Println("warning:", err)
	}
	return nil
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
func (z *Zip) Compress() error {
	if e := z.initForCompress(); e != nil {
		return e
	}
	defer func() {
		if e := z.Close(); e != nil {
			panic(e)
		}
	}()
	for i := 0; i < len(z.ZipItem); i++ {
		if err := z.ZipItem[i].WriteFile(z.zipWriter); err != nil {
			return err
		}
	}
	return nil
}

// GetFiles /////////////  Read zip file  ////////////////
func (z *Zip) GetFiles() []*zip.File {
	r, err := zip.OpenReader(z.ZipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return nil
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	return r.File
}
