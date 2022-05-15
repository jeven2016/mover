package main

import (
	"mover/src/common"
)

var log = common.Log

func main() {
	log.Infoln("Mover program started")
	setting, err := common.GetParams()
	if err != nil {
		log.WithError(err).Error("Exited:")
	} else {
		log.Infoln("the arguments parsed:", setting.String())
		common.Detect(*setting)
	}

}
