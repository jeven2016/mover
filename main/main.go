package main

import (
	"github.com/duke-git/lancet/v2/convertor"
	common2 "mover/common"
)

var log = common2.Log

func main() {
	log.Infoln("Mover program started")
	if setting, err := common2.GetParams(); err != nil {
		log.WithError(err).Error("Program Exited:")
	} else {
		str := convertor.ToString(*setting)
		log.Infoln("the arguments parsed:", str)
		common2.Detect(*setting)
		log.Infoln("Program started")
		select {}
	}

}
