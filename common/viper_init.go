package common

import (
	"github.com/spf13/viper"
)

var config = new(Parameters)

func SetupViper() (*Parameters, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load conf.ini: %v", err)
		return nil, err
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Failed to convert the data in conf.ini: %v", err)
		return nil, err
	}

	// 监控配置文件变化
	//viper.WatchConfig()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("Config file is changed, reload it now")
	//	if err := viper.Unmarshal(config); err != nil {
	//		panic(fmt.Errorf("The config cann't be updated:%s \n", err))
	//	}
	//})

	return config, nil
}

func getConfig() *Parameters {
	return config
}
