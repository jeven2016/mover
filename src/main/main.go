package main

import (
	"github.com/sirupsen/logrus"
	"mover/src/common"
)

var log = common.Log

func main() {
	log.Infoln("main ready")
	setting, err := common.GetParams()
	if err != nil {
		println(err)
	}
	log.WithFields(logrus.Fields{
		"from": setting.From,
		"to":   setting.To,
	}).Infoln("the arguments retrieved")
}
