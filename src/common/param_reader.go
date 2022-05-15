package common

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
)

var log = Log

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
	var setting = &Setting{From: "", To: ""}

	if err != nil {
		Log.Printf("failed to read config.ini: %v", err)
		return setting, err
	}
	from = cfg.Section("setting").Key("from").String()
	to = cfg.Section("setting").Key("to").String()

	if from == "" || to == "" {
		return setting, errors.New(fmt.Sprintf("invalid arguments: from and to are required ,  from=%v, to=%v", from, to))
	}

	setting.From = from
	setting.To = to
	return setting, nil
}

func params() (bool, string, string) {
	from := flag.String("from", "", "the source directory to check")
	to := flag.String("to", "", "the destination directory to move")
	flag.Parse()
	if *from == "" || *to == "" {
		return false, "", ""
	}
	return true, *from, *to
}
