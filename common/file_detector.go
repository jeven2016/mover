package common

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/sirupsen/logrus"
)

// TODO: 1. 文件夹下有多个视频文件，只有最大的文件被拷贝，如何支持
//      2. 文件夹要最后删除，如果提前删除会导致子目录下的文件无法拷贝

var nameReg = regexp.MustCompile("(\\[.*?])|(.*?原版首发_)|(.*?@)|(_uncensored)")

var wg = sync.WaitGroup{}

// 结果统计类
var stats = &Stats{}

type Callback func(int, string) bool

func handleFile(setting *Setting, folderInfo FolderInfo) {
	sourcePath := folderInfo.sourceDir

	hasDownloadingFile := false
	hasError := false
	if fileNames, err := fileutil.ListFileNames(sourcePath); err == nil {
		for _, name := range fileNames {
			// 判断文件后缀是否是视频
			if !validateFiles(name, sourcePath, setting) {
				atomic.AddInt32(&stats.IgnoreFiles, 1)

				if strings.HasSuffix(name, ".bc!") {
					hasDownloadingFile = true
				}
				continue
			}
			folderPath := strings.ReplaceAll(sourcePath, folderInfo.baseSourceDir, "")
			linuxPath := filepath.ToSlash(folderPath)
			destPath := filepath.Join(setting.To, linuxPath)

			realSource := filepath.Join(sourcePath, name)
			realDest := filepath.Join(destPath, getPureName(name))

			// 拷贝文件
			success := copyTo(realSource, realDest, &stats.FileCount, &stats.FileFailure)

			if !success {
				hasError = true
				continue
			}

			picFileName := findPicture(fileNames, sourcePath, name, setting)
			if len(picFileName) > 0 {
				// 拷贝图片
				success = copyTo(filepath.Join(sourcePath, picFileName), filepath.Join(destPath, picFileName), &stats.PictureCount, &stats.PictureFailure)
				if !success {
					hasError = true
				}
			}

		}

		if !hasError && !hasDownloadingFile {
			// 删除源目录, 最后删除才行
			// if err := os.RemoveAll(sourcePath); err != nil {
			// 	log.WithFields(logrus.Fields{
			// 		"sourcePath": sourcePath,
			// 	}).WithError(err).Errorln("failed to delete the source directory")
			// 	atomic.AddInt32(&stats.RemoveOldFoldersFailure, 1)
			// } else {
			// 	atomic.AddInt32(&stats.RemoveOldFolders, 1)
			// }
		}
	}
}

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

func supportedFileExtension(extensions []string, callback Callback) bool {
	_, validFile := slice.Find(extensions, callback)
	return validFile
}

func getPureName(name string) string {
	realName := name
	return nameReg.ReplaceAllString(realName, "")
}

func copyTo(source string, dest string, successCount *int32, failureCount *int32) bool {
	if err := doCopy(source, dest); err == nil {
		atomic.AddInt32(successCount, 1)
		return true
	} else {
		atomic.AddInt32(failureCount, 1)
		log.WithError(err).Printf("Failed to move: %v\n", source)
		return false
	}
}

func doCopy(source string, dest string) error {
	if err := ensureDestDir(filepath.Dir(dest)); err == nil {
		err := fileutil.CopyFile(source, dest)
		if err != nil {
			log.WithFields(logrus.Fields{
				"source": source,
				"dest":   dest,
			}).WithError(err).Error("Failed to copyTo this file")

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

func findPicture(fileNames []string, directory string, fileName string, setting *Setting) string {
	for _, name := range fileNames {
		if !validPictures(name, directory, setting) {
			continue
		}
		realName := strings.Split(name, ".")[0]
		if strings.Contains(fileName, realName) {
			// 找到最符合的图片
			return name
		}
	}
	return ""
}

func Detect(setting *Setting) {

	// 遍历所有的目录并插入到chan中
	// 第一级目录下的文件全部拷贝，子目录只拷贝特定的文件
	for _, from := range setting.From {
		detectDirs(from, setting, from, setting.To, 0)
	}
	close(FolderChan)

	stats.FolderCount = int32(len(FolderChan))

	progressSetting := &ProgressSetting{
		Total: int64(stats.FolderCount),
		Units: &progress.UnitsDefault,
	}

	// init progress
	InitProgress()
	AddTracker(progressSetting, &wg)

	var current int32 = 0

	// 4 processors
	// 处理各个目录下的文件
	for i := 0; i < 4; i++ {
		wg.Add(1)
		// goroutine
		go func() {
			defer wg.Done()
			for folderInfo := range FolderChan {
				handleFile(setting, folderInfo)
				GetTracker().Increment(1)

				atomic.AddInt32(&current, 1)
				if atomic.LoadInt32(&current) == stats.FolderCount {
					GetTracker().MarkAsDone()
				}
			}
		}()
	}

	// time.Sleep(30 * time.Second)
	wg.Wait()

	ShowTable(stats)
}

func detectDirs(originSourcePath string, setting *Setting, sourcePath string, parentPath string, depth int32) {
	// 分析源目录下的文件
	dirs, err := os.ReadDir(sourcePath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"sourcePath": parentPath,
		}).WithError(err).Warnln("failed To check this directory:")
	}

	// 将当前目录添加到channel中
	FolderChan <- FolderInfo{
		sourceDir:     sourcePath,
		baseSourceDir: originSourcePath,
		setting:       setting,
		copyAllFiles:  depth == 0, // 第一层级拷贝所有合适的文件
	}

	for _, dirEntry := range dirs {
		if dirEntry.IsDir() {
			detectDirs(originSourcePath, setting, filepath.Join(sourcePath, dirEntry.Name()),
				filepath.Join(parentPath, dirEntry.Name()), depth+1)
		}
	}
}

func validPictures(name string, directory string, setting *Setting) bool {
	validPics := supportedFileExtension(setting.PictureExtension, func(i int, suffix string) bool {
		return strings.HasSuffix(name, suffix)
	})
	return validPics && checkFileSize(filepath.Join(directory, name), setting, PicType)
}

func validateFiles(name string, directory string, setting *Setting) bool {
	valid := supportedFileExtension(setting.FileExtension, func(i int, suffix string) bool {
		return strings.HasSuffix(name, suffix)
	})
	return valid && checkFileSize(filepath.Join(directory, name), setting, FileType)
}

func checkFileSize(filePath string, setting *Setting, resType ResourceType) bool {
	if fileInfo, err := os.Stat(filePath); err == nil {
		fileSize := fileInfo.Size()
		var fileLimit int64
		switch resType {
		case FileType:
			fileLimit = setting.MinSize(setting.FileMinSize, setting.FileMinSizeUnit)
		case PicType:
			fileLimit = setting.MinSize(setting.PicMinSize, setting.PicMinSizeUnit)
		}
		return fileSize >= fileLimit
	}
	return false
}
