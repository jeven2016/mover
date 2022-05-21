package common

import "math"

type Setting struct {
	From                string
	To                  string
	fileExtension       []string
	fileMinSize         int64
	fileMinSizeUnit     string
	checkPicture        bool
	pictureExtension    []string
	picMinSize          int64
	picMinSizeUnit      string
	createRootDirectory bool
}

func (s *Setting) Defaults() *Setting {
	s.fileExtension = []string{".mp4", ".mkv", ".avi"}
	s.fileMinSize = 500 * 1024 * 1024 //500MB
	s.checkPicture = true
	s.createRootDirectory = true
	return s
}

func (s *Setting) MinSize(baseValue int64, unit string) int64 {
	var initialVal int64
	switch unit {
	case "B":
		initialVal = 1
	case "KB":
		initialVal = 1024
	case "MB":
		initialVal = int64(math.Pow(1024, 2))
	case "GB":
		initialVal = int64(math.Pow(1024, 3))
	default:
		initialVal = 1
	}

	return baseValue * initialVal
}
