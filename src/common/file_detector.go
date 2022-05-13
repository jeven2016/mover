package common

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"os"
)

func Detect(setting Setting) error {
	dirs, err := os.ReadDir(setting.from)
	if err != nil {
		return err
	}

	for _, sourceDir := range dirs {
		log.Infoln("sourceDir is", sourceDir.Name())
	}
	fileutil.IsDir("")
	return nil
}
