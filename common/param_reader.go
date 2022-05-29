package common

import (
	"errors"
	"flag"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"gopkg.in/ini.v1"
	"regexp"
	"strconv"
	"strings"
)

var log = Log
var fileSizeRegex = regexp.MustCompile("(\\d*)(.*)")
var fileSizeUnitTypes = []string{"MB", "KB", "GB", "TB", "B"}

func GetParams() (*Setting, error) {
	defer func() {
		//转换成error使用，而不是使用any类型
		if e := recover(); e != nil {
			log.WithError(e.(error)).Info("error occurs")
		}
	}()

	var (
		err     error
		from    string
		to      string
		exists  bool
		setting *Setting
	)
	exists, from, to = params()
	if !exists {
		setting, err = ReadIni()
		if err != nil {
			println(err)
			return nil, err
		}
	} else {
		setting = &Setting{
			From: from,
			To:   to,
		}
	}
	return setting, nil
}

func ReadIni() (*Setting, error) {
	var (
		from string
		to   string
	)
	cfg, err := ini.Load("./config.ini")

	if err != nil {
		Log.Printf("failed To read config.ini: %v", err)
		return nil, err
	}
	var setting = &Setting{}
	setting.Defaults()

	from = cfg.Section("setting").Key("from").String()
	to = cfg.Section("setting").Key("to").String()
	fileExtension := cfg.Section("setting").Key("file_extension").String()
	fileMinSize := cfg.Section("setting").Key("file_min_size").String()
	checkPicture, chkErr := cfg.Section("setting").Key("check_picture").Bool()
	pictureExtension := cfg.Section("setting").Key("picture_extension").String()
	picMinSize := cfg.Section("setting").Key("pic_min_size").String()
	createRootDirectory, cerr := cfg.Section("setting").Key("create_root_directory").Bool()

	if from == "" || to == "" {
		return setting, errors.New(fmt.Sprintf("invalid arguments: From and To are required ,  From=%v, To=%v", from, to))
	}

	if !fileutil.IsExist(strings.Trim(from, " ")) {
		return setting, errors.New("argument From isn't a valid directory")
	}

	if chkErr != nil {
		return nil, errors.New("check_picture isn't a boolean value")
	}

	if cerr != nil {
		return nil, errors.New("create_root_directory isn't a boolean value")
	}

	//parse the file_min_size
	setting.FileMinSize, setting.FileMinSizeUnit, err = parseNumberAndUnit(fileMinSize)
	if err != nil {
		return nil, err
	}

	//parse the pic_min_size
	setting.PicMinSize, setting.PicMinSizeUnit, err = parseNumberAndUnit(picMinSize)
	if err != nil {
		return nil, err
	}

	setting.From = from
	setting.To = to
	setting.FileExtension = strings.Split(fileExtension, ",")

	setting.CheckPicture = checkPicture
	setting.PictureExtension = strings.Split(pictureExtension, ",")
	setting.CreateRootDirectory = createRootDirectory

	return setting, nil
}

func parseNumberAndUnit(minSizeString string) (int64, string, error) {
	//parse the file_min_size
	subStrings := fileSizeRegex.FindStringSubmatch(minSizeString)
	subLen := len(subStrings)
	if subLen == 3 {
		fileSize, err := strconv.ParseInt(subStrings[1], 10, 64)
		if err != nil || !checkFileSizeUnitType(subStrings[2]) {
			return 0, "", errors.New("file_min_size is invalid, the value should be in this format '200KB' or '300MB'")
		}
		return fileSize, subStrings[2], nil
	}
	return 0, "", nil
}

func params() (bool, string, string) {
	from := flag.String("From", "", "the source directory To check")
	to := flag.String("To", "", "the destination directory To move")
	flag.Parse()
	if *from == "" || *to == "" {
		return false, "", ""
	}
	return true, *from, *to
}

func checkFileSizeUnitType(fileSizeUnitType string) bool {
	for _, fs := range fileSizeUnitTypes {
		if strings.EqualFold(fileSizeUnitType, fs) {
			return true
		}
	}
	return false
}
