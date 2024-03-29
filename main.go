package main

import (
	"github.com/duke-git/lancet/v2/convertor"

	"mover/common"
)

func main() {
	var log = common.Log

	config, err := common.SetupViper()
	if err != nil {
		panic(err)
	}

	// common.ShowProgress()
	if setting, err := common.Validate(config.Parameters); err != nil {
		log.WithError(err).Error("An error occurs")
	} else {
		str := convertor.ToString(*setting)
		log.Infoln("the parameters parsed:", str)
		common.Detect(setting)
	}

}
