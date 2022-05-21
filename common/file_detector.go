package common

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Callback func(int, string) bool

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

func supportedFileExtension(extensions []string, callback Callback) (supported bool) {
	_, validFile := slice.Find(extensions, callback)
	return validFile
}

func handleFile(setting Setting, sourcePath string, destPath string) {
	if fileNames, err := fileutil.ListFileNames(sourcePath); err == nil {

		//查找最匹配且容量最大的文件
		var maxFileSize int64
		var curFileInfo os.FileInfo

		for _, name := range fileNames {
			//判断文件后缀是否支持
			validFile := supportedFileExtension(setting.fileExtension, func(i int, suffix string) bool {
				return strings.HasSuffix(name, suffix)
			})
			if !validFile {
				continue
			}

			if fileInfo, err := os.Stat(filepath.Join(sourcePath, name)); err == nil {
				fileSize := fileInfo.Size()
				if maxFileSize == 0 ||
					(fileSize > maxFileSize && fileSize >=
						setting.MinSize(setting.fileMinSize, setting.fileMinSizeUnit)) {
					maxFileSize, curFileInfo = fileSize, fileInfo
				}
			}
		}

		if curFileInfo != nil {
			picFileName := findPicture(setting, fileNames, curFileInfo)

			folderPath := strings.ReplaceAll(sourcePath, setting.From, "")
			linuxPath := filepath.ToSlash(folderPath)

			if picFileName == "" {
				copyFile(linuxPath, setting, sourcePath, curFileInfo)
			} else {
				destPath := filepath.Join(setting.To, linuxPath)
				err := ensureDestDir(destPath)
				if err != nil {
					log.WithFields(logrus.Fields{
						"destPath": destPath,
					}).WithError(err).Error("failed to create directory:")
					return
				}

				CopiedLog.Infoln("Copied：", filepath.Join(sourcePath, curFileInfo.Name()), "->",
					filepath.Join(destPath, curFileInfo.Name()))

				CopiedLog.Infoln("Copied：", filepath.Join(sourcePath, curFileInfo.Name()), "->",
					filepath.Join(destPath, curFileInfo.Name()))
			}
		}
	}
}

func copyFile(linuxPath string, setting Setting, sourcePath string, curFileInfo os.FileInfo) {
	//当没有图片时，拷贝到根目录
	//source: F:\new_down\, 第一级目录：\冲1
	//F:\new_down\冲1\xx.mp4  -> 不会拷贝上一层目录
	//F:\new_down\冲1\fold1\mm.mp4  -> 拷贝到d:\dest\冲1\下

	var catalogFolder string = linuxPath
	if strings.HasPrefix(linuxPath, "/") {
		catalogFolder = strings.TrimLeft(catalogFolder, "/")
	}
	if strings.Contains(catalogFolder, "/") {
		catalogFolder = strings.Split(catalogFolder, "/")[0]
	}

	destFilePath := filepath.Join(setting.To, catalogFolder, curFileInfo.Name())
	err := ensureDestDir(destFilePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"destPath": destFilePath,
		}).WithError(err).Error("failed to create directory:")
		return
	}

	CopiedLog.Infoln("Copied：", filepath.Join(sourcePath, curFileInfo.Name()), "->", destFilePath)
}

func findPicture(setting Setting, fileNames []string, curFileInfo os.FileInfo) string {
	for _, name := range fileNames {
		validPicFile := supportedFileExtension(setting.pictureExtension, func(i int, suffix string) bool {
			return strings.HasSuffix(name, suffix)
		})
		if !validPicFile {
			continue
		}
		realName := strings.Split(name, ".")[0]
		if strings.Contains(curFileInfo.Name(), realName) {
			//找到最符合的图片
			return name
		}
	}
	return ""
}

func Detect(setting Setting) {
	//遍历所有的目录并插入到chan中
	go detectDirs(setting, setting.From, setting.To)

	//处理各个目录下的文件
	go func(s Setting) {
		for folderInfo := range FolderChan {
			handleFile(s, folderInfo.sourceDir, folderInfo.ToDestDir())
		}
	}(setting)
}

func detectDirs(setting Setting, sourcePath string, parentPath string) {
	//分析源目录下的文件
	dirs, err := os.ReadDir(sourcePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"sourcePath": parentPath,
		}).WithError(err).Warnln("failed To check this directory:")
	}

	//将当前目录添加到channel中
	FolderChan <- FolderInfo{
		sourceDir: sourcePath,
		setting:   setting,
	}

	for _, dirEntry := range dirs {
		if dirEntry.IsDir() {
			detectDirs(setting, filepath.Join(sourcePath, dirEntry.Name()),
				filepath.Join(parentPath, dirEntry.Name()))
		}
	}
}
