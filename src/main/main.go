package main

import (
	"github.com/duke-git/lancet/v2/convertor"
	"mover/src/common"
)

var log = common.Log

func main() {
	log.Infoln("Mover program started")
	if setting, err := common.GetParams(); err != nil {
		log.WithError(err).Error("Program Exited:")
	} else {
		str := convertor.ToString(*setting)
		log.Infoln("the arguments parsed:", str)
		common.Detect(*setting)
		log.Infoln("Program started")
		select {}
	}

}
