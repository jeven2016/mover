package common

type Setting struct {
	from                string
	to                  string
	fileExtension       []string
	fileMinSize         int64
	fileMinSizeUnit     string
	checkPicture        bool
	createRootDirectory bool
}

func (s *Setting) Defaults() *Setting {
	s.fileExtension = []string{".mp4", ".mkv", ".avi"}
	s.fileMinSize = 500 * 1024 * 1024 //500MB
	s.checkPicture = true
	s.createRootDirectory = true
	return s
}

func (s *Setting) MinSize() int64 {
	var initialVal int64
	switch s.fileMinSizeUnit {
	case "B":
		initialVal = 1
	case "KB":
		initialVal = 1024
	case "MB":
		initialVal = 1024 * 1024
	case "GB":
		initialVal = 1024 * 1024 * 1024
	default:
		initialVal = 1
	}

	return initialVal
}

var FileChan = make(chan string)
