package tools

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"plugin"
)

func PluginRun(fileName string, symbol string, args ...interface{}) (ret interface{}, err error) {
	p, err := plugin.Open(fmt.Sprintf("%s/%s", beego.AppConfig.DefaultString("pluginDir", "./plugins"), fileName))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"fileName": fileName,
			"symbol": symbol,
		}).Error(err)
		return nil, err
	}
	f, err := p.Lookup(symbol)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"fileName": fileName,
			"symbol": symbol,
		}).Error(err)
		return nil, err
	}
	fn, ok := f.(func(...interface{}) (interface{}, error))
	if ok == false {
		logrus.Error(fmt.Sprintf("%+v", errors.Wrap(err, "func type error")))
		logrus.WithFields(logrus.Fields{
			"fileName": fileName,
			"symbol": symbol,
		}).Error(err)
		return nil, err
	}
	//New 返回空接口
	return fn(args...)
}

//func PluginRunDebug(fileName string, symbol string, args ...interface{}) (ret interface{}, err error) {
//	return qt.Run(args...)
//}
