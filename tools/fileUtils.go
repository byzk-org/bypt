package tools

import (
	"os"
	"path/filepath"
)



func CreateFile(filename string) {
	path, _ := filepath.Abs(filename)
	if !FileIsExist(path) {
		if file, err := os.Create(path); err != nil {
			ErrOut("创建配置文件失败，请查看当前用户对当前目录是否具有操作权!")
			os.Exit(-1)
		} else {
			_ = file.Close()
		}

	}

}
