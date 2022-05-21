package common

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Log = logrus.New()
var CopiedLog = logrus.New()

func initLog(log *logrus.Logger, out io.Writer, displayCaller bool) {
	// Log as JSON instead of the default ASCII formatter.
	//Log.SetFormatter(&logrus.JSONFormatter{})
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		// FullTimestamp:true,
		// DisableLevelTruncation:true,
	})

	// Output To stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(out)

	// Only Log the warning severity or above.
	log.SetLevel(logrus.InfoLevel)
	if displayCaller {
		log.SetReportCaller(true) //打印代码信息
	}
}

//程序运行前就会初始化
func init() {
	initLog(Log, os.Stdout, true)

	logfile, _ := os.OpenFile("./file-copied.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	initLog(CopiedLog, logfile, false)
}
