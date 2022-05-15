package common

import (
	"fmt"
	"strings"
)

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

func (s Setting) String() string {
	return fmt.Sprintf(`
	{
		from                : %v,
		to                  : %v,
		fileExtension       : %v
		fileMinSize         : %v,
		fileMinSizeUnit     : %v,
		checkPicture        : %v,
		createRootDirectory : %v,
    }
`, s.from, s.to, strings.Join(s.fileExtension, ","), s.fileMinSize, s.fileMinSizeUnit, s.checkPicture,
		s.createRootDirectory)
}
