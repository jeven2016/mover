package main

import (
	"github.com/duke-git/lancet/v2/convertor"
	"mover/src/common"
)

var log = common.Log

func main() {
	log.Infoln("Mover program started")
	setting, err := common.GetParams()
	if err != nil {
		log.WithError(err).Error("Exited:")
	} else {
		str := convertor.ToString(*setting)
		log.Infoln("the arguments parsed:", str)
		common.Detect(*setting)
		//select {}
	}

}
