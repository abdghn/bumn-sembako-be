package helper

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func CommonLogger() *logrus.Logger {
	log := logrus.New()
	path := "./logs/common.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(30*24*3600)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(24*3600)*time.Second),
	)

	log.SetReportCaller(true)
	//print to multiple medium
	log.SetOutput(io.MultiWriter(writer, os.Stdout))
	return log
}
