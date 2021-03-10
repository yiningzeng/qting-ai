package plugins

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"qting-ai/models"
	"qting-ai/tools"
)
var Watcher *fsnotify.Watcher
func AddPluginInfo(fileName string) {
	if path.Ext(fileName) == beego.AppConfig.DefaultString("pluginSuffix", ".yn") {
		if res, err := tools.PluginRun(fileName, "Version", fileName); err == nil {
			qtPlugin := res.(models.QtPlugins)
			qtPlugin.PluginName = fileName
			_, _ = models.AddQtPlugins(&qtPlugin)
		} else {
			logrus.Info(err)
		}
	}
}
func StartWatcher() {
	dir := beego.AppConfig.DefaultString("pluginDir", "./plugins")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Error(fmt.Sprintf("%+v", errors.Wrap(err, "监听文件夹出错")))
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

func StopWatcher() {
	Watcher.Close()
}
