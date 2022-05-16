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

var FileChan = make(chan string)
