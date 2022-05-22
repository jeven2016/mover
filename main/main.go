package main

import (
	"github.com/duke-git/lancet/v2/convertor"
	"mover/common"
)

var log = common.Log

func main() {
	defer log.Infoln("All jobs finished, exit now")

	log.Infoln("Mover program started")
	setting, err := common.GetParams()
	if err != nil {
		log.WithError(err).Error("Program Exited:")
		return
	}

	str := convertor.ToString(*setting)
	log.Infoln("the arguments parsed:", str)
	log.Infoln("Program started")
	common.Detect(*setting)

}
