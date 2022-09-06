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
	err, fileSize, fileSizeUnit := parseFileSize(params.MinSize)
	if err != nil {
		return nil, err
	}
	setting.FileMinSize = fileSize
	setting.FileMinSizeUnit = fileSizeUnit

	//parse the pic_min_size
	err, fileSize, fileSizeUnit = parseFileSize(params.MinSize)
	if err != nil {
		return nil, err
	}
	setting.PicMinSize = fileSize
	setting.PicMinSizeUnit = fileSizeUnit

	setting.From = sourcePaths
	setting.To = to
	setting.FileExtension = strings.Split(params.FileExtension, ",")
	setting.PictureExtension = strings.Split(params.PicExtension, ",")

	setting.CheckPicture = params.CheckPicture
	setting.CreateRootDirectory = params.CreateRootDirectory

	return setting, nil
}

func parseFileSize(minSize string) (error, int64, string) {
	subStrings := fileSizeRegex.FindStringSubmatch(minSize)
	subLen := len(subStrings)
	errorMsg := "file_min_size or pic_min_size is invalid, the value should be in this format: 300MB"
	if subLen == 3 {
		fileSize, err := convertor.ToInt(subStrings[1])
		if err != nil || !checkFileSizeUnitType(subStrings[2]) {
			return errors.New(errorMsg), 0, ""
		}
		return nil, fileSize, subStrings[2]
	}
	return errors.New(errorMsg), 0, ""
}

func checkFileSizeUnitType(fileSizeUnitType string) bool {
	for _, fs := range fileSizeUnitTypes {
		if strings.EqualFold(fileSizeUnitType, fs) {
			return true
		}
	}
	return false
}
