package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

//程序运行前就会初始化
func init() {
	// Log as JSON instead of the default ASCII formatter.
	//Log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	Log.SetOutput(os.Stdout)

	// Only Log the warning severity or above.
	Log.SetLevel(logrus.InfoLevel)
	Log.SetReportCaller(true) //打印代码信息
	Log.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	Log.Infoln("Log component is initialized")
}
