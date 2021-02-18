package version

import (
	"github.com/sirupsen/logrus"
)

var ID = "undefined"
var BuildDate = "undefined"

func PrintVersion()  {
	logrus.WithFields(logrus.Fields{
		"VersionId": ID,
		"BuildDate": BuildDate,
	}).Info("当前版本信息")
}

