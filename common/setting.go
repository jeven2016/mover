package common

import "math"

type Setting struct {
	From                []string
	To                  string
	FileExtension       []string
	FileMinSize         int64
	FileMinSizeUnit     string
	CheckPicture        bool
	PictureExtension    []string
	PicMinSize          int64
	PicMinSizeUnit      string
	CreateRootDirectory bool
}

func (s *Setting) Defaults() *Setting {
	s.FileExtension = []string{".mp4", ".mkv", ".avi"}
	s.FileMinSize = 500 * 1024 * 1024 //500MB
	s.CheckPicture = true
	s.CreateRootDirectory = true
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
