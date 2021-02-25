package plugins

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"qting-ai/models"
	"time"
)
func AddPluginInfo(f string) {
	if path.Ext(f) == beego.AppConfig.DefaultString("pluginSuffix", ".yn") {
		if res, err := Run(f, "Version", f); err == nil {
			qtPlugin := res.(models.QtPlugins)
			qtPlugin.PluginName = f
			_, _ = models.AddQtPlugins(&qtPlugin)
		} else {
			logrus.Info(err)
		}
	}
}
func Watcher() {
	dir := beego.AppConfig.DefaultString("pluginDir", "./plugins")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Error(err)
		return
	}
	err = watcher.Add(dir)

	//defer watcher.Close()
	//通过Walk来遍历目录下的所有子目录
	err = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		AddPluginInfo(info.Name())
		return nil
	})

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create{
					AddPluginInfo(path.Base(ev.Name))
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					f := path.Base(ev.Name)
					logrus.WithField("filename", f).Info("remove")
					if path.Ext(f) == beego.AppConfig.DefaultString("pluginSuffix", ".yn") {
						_ = models.DeleteQtPluginsByFileName(f)
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					f := path.Base(ev.Name)
					logrus.WithField("filename", f).Info("Rename")
					if path.Ext(f) == beego.AppConfig.DefaultString("pluginSuffix", ".yn") {
						_ = models.UpdateQtPluginsByMV("QTing-tiny-3l", 10, &models.QtPlugins{
							Module:      "QTing-tiny-3l",
							VersionCode: 10,
							Symbol:      "asdsd",
							Args:        "asds",
							PluginName:  "asds",
							CreateTime:  time.Now(),
						})
						//if res, err := Run(f, "Version"); err == nil {
						//	qtPlugin := res.(models.QtPlugins)
						//	qtPlugin.PluginName = f
						//	_ = models.UpdateQtPluginsByMV(qtPlugin.Module, qtPlugin.VersionCode, &qtPlugin)
						//} else {
						//	logrus.Info(err)
						//}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()
	<-done
}

func Run(fileName string, symbol string, args ...interface{}) (ret interface{}, err error) {
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
		logrus.Error("func type error: New.")
		logrus.WithFields(logrus.Fields{
			"fileName": fileName,
			"symbol": symbol,
		}).Error(err)
		return nil, err
	}
	//New 返回空接口
	return fn(args...)
}

