package common

import (
	"github.com/duke-git/lancet/v2/slice"
	"os"
	"path/filepath"
)

func Detect(setting Setting) error {
	dirs, err := os.ReadDir(setting.from)
	if err != nil {
		return err
	}

	for _, sourceDir := range dirs {
		ext := filepath.Ext(sourceDir.Name())
		//find the files with corresponding extension
		if slice.Contain(setting.fileExtension, ext) {
			println(ext)
		}
	}
	return nil
}
