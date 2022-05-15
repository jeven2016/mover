package common

import (
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
	return nil
}
