package FileDirUtils

import (
	"os"
	"path/filepath"
)

func FileIsExist(filePath string) bool {
	// 使用 os.Stat 函数检查文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		//fmt.Println("文件存在")
		return true
	} else if os.IsNotExist(err) {
		//fmt.Println("文件不存在")
		return false
	} else {
		//fmt.Println("发生了其他错误:", err)
		panic(err)
	}
}

func GetAllFilesInFolder(folderPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
