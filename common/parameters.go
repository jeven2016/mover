package common

import (
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"regexp"
	"strings"
)

var log = Log
var fileSizeRegex = regexp.MustCompile("(\\d*)(.*)")
var fileSizeUnitTypes = []string{"MB", "KB", "GB", "TB", "B"}

func Validate(params *Parameters) (*Setting, error) {
	var (
		from = params.From
		to   = params.To
	)
	var setting = &Setting{}

	if from == "" || to == "" {
		return setting, errors.New(fmt.Sprintf("invalid arguments: from and to are required ,  from=%v, to=%v", from, to))
	}

	//check from variable
	var sourcePaths = strings.Split(from, ",")
	sourcePaths = slice.Filter(sourcePaths, func(index int, value string) bool {
		return len(value) > 0
	})

	for _, path := range sourcePaths {
		if !fileutil.IsExist(strings.Trim(path, " ")) {
			return setting, errors.New("invalid directories defined for variable 'from'")
		}
	}

	//parse the file_min_size
	subStrings := fileSizeRegex.FindStringSubmatch(params.MinSize)
	subLen := len(subStrings)
	if subLen == 3 {
		fileSize, err := convertor.ToInt(subStrings[1])
		if err != nil || !checkFileSizeUnitType(subStrings[2]) {
			return nil, errors.New("file_min_size is invalid, the value should be in this format: 300MB")
		}
		setting.FileMinSize = fileSize
		setting.FileMinSizeUnit = subStrings[2]
	}

	setting.From = sourcePaths
	setting.To = to
	setting.FileExtension = strings.Split(params.FileExtension, ",")

	setting.CheckPicture = params.CheckPicture
	setting.CreateRootDirectory = params.CreateRootDirectory

	return setting, nil
}

func checkFileSizeUnitType(fileSizeUnitType string) bool {
	for _, fs := range fileSizeUnitTypes {
		if strings.EqualFold(fileSizeUnitType, fs) {
			return true
		}
	}
	return false
}
