package common

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var nameReg = regexp.MustCompile("(\\[.*?])|(.*?原版首发_)|(.*?@)|(_uncensored)")

var wg = sync.WaitGroup{}

type Callback func(int, string) bool

func ensureDestDir(destPath string) error {
	exist := fileutil.IsExist(destPath)
	if !exist {
		log.Infoln("Will create directory:", destPath)
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

func getPureName(name string) string {
	realName := name
	return nameReg.ReplaceAllString(realName, "")
}

func doCopy(source string, dest string) error {
	if err := ensureDestDir(filepath.Dir(dest)); err == nil {
		err := fileutil.CopyFile(source, dest)
		if err != nil {
			log.WithFields(logrus.Fields{
				"source": source,
				"dest":   dest,
			}).WithError(err).Error("Failed to copy this file")

			if err = fileutil.RemoveFile(dest); err != nil {
				log.WithFields(logrus.Fields{
					"source": source,
					"dest":   dest,
				}).WithError(err).Error("Failed to remove this file in destination directory")
			}
			return err
		}
		if err = fileutil.RemoveFile(source); err != nil {
			log.WithFields(logrus.Fields{
				"source": source,
				"dest":   dest,
			}).WithError(err).Error("Failed to remove source file after the file is copied")
		}
		return err
	}
	return nil
}

func handleFile(setting *Setting, folderInfo FolderInfo) {
	sourcePath := folderInfo.sourceDir
	if fileNames, err := fileutil.ListFileNames(sourcePath); err == nil {

		if folderInfo.copyAllFiles {
			for _, name := range fileNames {
				//判断文件后缀是否支持
				validFile := supportedFileExtension(setting.FileExtension, func(i int, suffix string) bool {
					return strings.HasSuffix(name, suffix)
				})
				if !validFile {
					println(filepath.Join(sourcePath, name), "ignored")
					continue
				}

				destPath := folderInfo.ToDestDir()

				newFileName := getPureName(name)
				realSource := filepath.Join(sourcePath, name)
				realDest := filepath.Join(destPath, newFileName)

				if err = doCopy(realSource, realDest); err == nil {
					CopiedLog.Infoln("Root File Moved：", realSource, "->", realDest)
				}
			}
			return
		}

		//查找最匹配且容量最大的文件
		var maxFileSize int64
		var curFileInfo os.FileInfo

		for _, name := range fileNames {
			//判断文件后缀是否支持
			validFile := supportedFileExtension(setting.FileExtension, func(i int, suffix string) bool {
				return strings.HasSuffix(name, suffix)
			})
			if !validFile {
				continue
			}

			if fileInfo, err := os.Stat(filepath.Join(sourcePath, name)); err == nil {
				fileSize := fileInfo.Size()
				if maxFileSize == 0 && fileSize >
					setting.MinSize(setting.FileMinSize, setting.FileMinSizeUnit) {
					maxFileSize, curFileInfo = fileSize, fileInfo
				}

				if maxFileSize != 0 && fileSize > maxFileSize {
					maxFileSize, curFileInfo = fileSize, fileInfo
				}
			}
		}

		if curFileInfo != nil {
			picFileName := findPicture(setting, fileNames, curFileInfo)

			folderPath := strings.ReplaceAll(sourcePath, folderInfo.baseSourceDir, "")
			linuxPath := filepath.ToSlash(folderPath)

			if picFileName == "" {
				copyFile(linuxPath, setting, sourcePath, curFileInfo)
			} else {
				destPath := filepath.Join(setting.To, linuxPath)
				newFileName := getPureName(curFileInfo.Name())
				realSource := filepath.Join(sourcePath, curFileInfo.Name())
				realDest := filepath.Join(destPath, newFileName)
				if err = doCopy(realSource, realDest); err == nil {
					CopiedLog.Infoln("Root File Moved：", realSource, "->", realDest)
				}

				realSource = filepath.Join(sourcePath, picFileName)
				realDest = filepath.Join(destPath, picFileName)
				if err = doCopy(realSource, realDest); err == nil {
					CopiedLog.Infoln("Pic Moved：", realSource, "->", realDest)
				}
			}

			if err == nil {
				//删除源目录
				if err := os.RemoveAll(sourcePath); err != nil {
					log.WithFields(logrus.Fields{
						"sourcePath": sourcePath,
					}).WithError(err).Errorln("failed to delete the source directory")
				}
			}
		}
	}
}

func copyFile(linuxPath string, setting *Setting, sourcePath string, curFileInfo os.FileInfo) {
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

	//确保目录已经创建
	destFilePath := filepath.Join(filepath.Join(setting.To, catalogFolder), getPureName(curFileInfo.Name()))
	realSource := filepath.Join(sourcePath, curFileInfo.Name())
	realDest := destFilePath
	if err := doCopy(realSource, realDest); err == nil {
		CopiedLog.Infoln("File Moved：", realSource, "->", realDest)
	}
}

func findPicture(setting *Setting, fileNames []string, curFileInfo os.FileInfo) string {
	for _, name := range fileNames {
		validPicFile := supportedFileExtension(setting.PictureExtension, func(i int, suffix string) bool {
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

func Detect(setting *Setting) {
	wg.Add(1)

	//遍历所有的目录并插入到chan中
	go func() {
		defer func() {
			log.Infoln("Finished to detect the source directory")
			close(FolderChan)
			wg.Done()
		}()
		//第一级目录下的文件全部拷贝，子目录只拷贝特定的文件
		for _, from := range setting.From {
			detectDirs(from, setting, from, setting.To, 0)
		}
	}()

	//4 processors
	//处理各个目录下的文件
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for folderInfo := range FolderChan {
				handleFile(setting, folderInfo)
			}
		}()
	}

	wg.Wait()
}

func detectDirs(originSourcePath string, setting *Setting, sourcePath string, parentPath string, depth int32) {
	//分析源目录下的文件
	dirs, err := os.ReadDir(sourcePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"sourcePath": parentPath,
		}).WithError(err).Warnln("failed To check this directory:")
	}

	//将当前目录添加到channel中
	FolderChan <- FolderInfo{
		sourceDir:     sourcePath,
		baseSourceDir: originSourcePath,
		setting:       setting,
		copyAllFiles:  depth == 1, //第一层级拷贝所有合适的文件
	}

	for _, dirEntry := range dirs {
		if dirEntry.IsDir() {
			detectDirs(originSourcePath, setting, filepath.Join(sourcePath, dirEntry.Name()),
				filepath.Join(parentPath, dirEntry.Name()), depth+1)
		}
	}
}
