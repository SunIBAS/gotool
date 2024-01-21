package FileDirUtils

import (
	"encoding/json"
	"os"
)

func WriteToJson(obj interface{}, saveFile string) error {
	if bytes, err := json.Marshal(obj); err != nil {
		err = os.WriteFile(saveFile, bytes, os.ModeAppend)
		return err
	} else {
		return err
	}
}
