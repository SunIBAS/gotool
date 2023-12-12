package FileDirUtils

import "path/filepath"

func BaseAction(s *string) {
	*s = filepath.Base(*s)
}

type FileWalk struct {
	folder       string
	fileAction   func(*string)
	folderAction func(*string)
}

// NewFileWalk todo 有必要吗
func NewFileWalk(folder string) *FileWalk {
	return &FileWalk{
		folder: folder,
	}
}

//func FileWalk(sourceFolder string) {
//	err := filepath.Walk(sourceFolder, func(_path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		relativePath, err := filepath.Rel(sourceFolder, _path)
//		if err != nil {
//			return err
//		}
//		//zipEntryPath := filepath.ToSlash(filepath.Join(filepath.Base(sourceFolder), relativePath))
//		if !info.IsDir() {
//			//fmt.Println(zipEntryPath)
//			basePath := filepath.Base(sourceFolder)
//		}
//		return nil
//	})
//}
