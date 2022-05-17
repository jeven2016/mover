package common

import (
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func ensureDestDir(destPath string) error {
	exist := fileutil.IsExist(destPath)
	if !exist {
		err := os.MkdirAll(destPath, os.ModePerm.Perm())
		if err != nil {
			return err
		}
	}
	return nil
}

func handleFile(setting Setting, dir os.DirEntry, parentPath string, checkParent bool) {
	if !checkParent {
		return
	}
	validFile := slice.Contain(setting.fileExtension, filepath.Ext(dir.Name()))
	if validFile {
		fileInfo, err := dir.Info()
		if err != nil {
			log.WithError(err).Warnf(fmt.Sprintf("the file is failed to handle: %v", dir.Name()))
			return
		}
		//file size should be greater than the defined minimum size
		if fileInfo.Size() >= setting.MinSize() {
			if err := ensureDestDir(parentPath); err != nil {
				log.WithFields(logrus.Fields{
					"parentPath": parentPath,
					"file":       dir.Name(),
				}).WithError(err).Warnln("failed to create parent directory for this file")
			}
			log.Infoln(fmt.Sprintf("file %v is moved to %v", fileInfo.Name(), filepath.Join(parentPath, dir.Name())))
		} else {
			log.WithFields(logrus.Fields{
				"parentPath": parentPath,
				"file":       dir.Name(),
			}).Infoln("the file is ignored since its size(%v) less than %v", fileInfo.Size())
		}
	}
}

func Detect(setting Setting) error {
	return detectSubDir(setting, setting.from, setting.to)
}

func detectSubDir(setting Setting, sourcePath string, parentPath string) error {
	dirs, err := os.ReadDir(sourcePath)

	for _, dir := range dirs {
		go func(dirEntry os.DirEntry) {
			if dirEntry.IsDir() {
				entryErr := detectSubDir(setting, filepath.Join(sourcePath, dirEntry.Name()),
					filepath.Join(parentPath, dirEntry.Name()))
				if entryErr != nil {
					log.WithFields(logrus.Fields{
						"sourcePath": parentPath,
						"parentPath": parentPath,
						"dir":        dirEntry.Name(),
					}).WithError(err).Warnln("failed to handle directory:")
				}
			} else {
				handleFile(setting, dirEntry, parentPath, true)
			}
		}(dir)
	}
	return err
}
