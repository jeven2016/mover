package common

import (
	"path/filepath"
	"strings"
)

var FolderChan = make(chan FolderInfo, 1000)

type FolderInfo struct {
	sourceDir    string
	setting      Setting
	hasPicture   bool
	pictureName  string
	fileName     string
	copyAllFiles bool
}

// ToDestDir 拼接需要创建的目录
func (f *FolderInfo) ToDestDir() string {
	suffix := strings.ReplaceAll(f.sourceDir, f.setting.From, "")
	dest := filepath.Join(f.setting.To, suffix)
	println("source=", f.sourceDir, "dest=", dest)
	return dest
}
