package common

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"os"
	"path/filepath"
)

func ensureDestDir(setting Setting) error {
	destPath := setting.to
	exist := fileutil.IsExist(destPath)
	if !exist {
		err := os.MkdirAll(destPath, os.ModePerm.Perm())
		if err != nil {
			return err
		}
	}
	return nil
}

func Detect(setting Setting) error {
	dirs, err := os.ReadDir(setting.from)

	for _, dir := range dirs {
		if dir.IsDir() {

			continue
		}
		validFile := slice.Contain(setting.fileExtension, filepath.Ext(dir.Name()))
		if validFile {
			println("valid file:", dir.Name())
		}
	}

	//for _, sourceDir := range dirs {
	//	ext := filepath.Ext(sourceDir.Name())
	//	//find the files with corresponding extension
	//	if slice.Contain(setting.fileExtension, ext) {
	//		println(ext)
	//	}
	//}
	return err
}
