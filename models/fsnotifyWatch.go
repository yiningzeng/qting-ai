package models

import (
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
	"time"
)
var watcher *fsnotify.Watcher
//监控目录
func WatchDir(dir string) (err error) {
	if watcher == nil {
		watcher, err = fsnotify.NewWatcher()
	}
	if err != nil {
		logrus.Error(err)
		return err
	}
	//defer watcher.Close()
	//通过Walk来遍历目录下的所有子目录
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() && strings.Contains(path, "backup") {
			path, err := filepath.Abs(path)
			if err != nil {
				logrus.Error(err)
			}
			err = watcher.Add(path)
			if err != nil {
				logrus.Error(err)
			}
			//fmt.Println("监控 : ", path)
		}
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
					logrus.WithField("filename", ev.Name).Debug("Create || Write")
					//获取新创建文件的信息，如果是目录，则加入监控中
					if file, err := os.Stat(ev.Name); err == nil {
						if file.IsDir() {
							_ = watcher.Add(ev.Name)
							logrus.WithField("filename", ev.Name).Debug("Add Watch")
						} else if strings.Contains(file.Name(), ".weights") {
							//好了，这里就需要开始
							taskId := strings.ReplaceAll(file.Name(), ".weights", "") // 这是任务标识
							QTingTiny3l(taskId, ev.Name)
						} else if strings.Contains(file.Name(), ".names") {
							taskId := strings.ReplaceAll(file.Name(), ".names", "") // 这是任务标识
							err = watchAddLabels(taskId, ev.Name)
						}  else if strings.Contains(file.Name(), ".suggest") {
							taskId := strings.ReplaceAll(file.Name(), ".suggest", "") // 这是任务标识
							err = watchAddSuggest(taskId, ev.Name)
						}
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					logrus.WithField("filename", ev.Name).Debug("Remove Watch")
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() && strings.Contains(ev.Name, "backup") {
						_ = watcher.Remove(ev.Name)
						logrus.WithField("filename", ev.Name).Debug("Remove Watch")
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					logrus.WithField("filename", ev.Name).Debug("Rename")
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() && strings.Contains(ev.Name, "backup") {
						_ = watcher.Remove(ev.Name)
						logrus.WithField("filename", ev.Name).Debug("Remove Watch")
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
	return err
}

func AddWatchDir(dir string) (err error) {
	if watcher !=nil{
		err = watcher.Add(dir)
	} else {
		watcher, err = fsnotify.NewWatcher()
		if err == nil {
			err = watcher.Add(dir)
		}
		//err = errors.New("watcher is nil")
	}
	return err
}

func RemoveWatchDir(dir string) (err error) {
	if watcher !=nil{
		err = watcher.Remove(dir)
	} else {
		err = errors.New("watcher is nil")
	}
	return err
}

func watchAddSuggest(taskId string, path string) (err error) {
	temp := strings.Split(strings.ReplaceAll(path, beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"), ""), "/")
	if qtProject, err := GetQtProjectsByProjectName(temp[0]); err == nil {
		//读取labels.names 如果不存在就直接不执行
		if by, err := afero.ReadFile(afero.NewOsFs(), path); err == nil {
			//temp := strings.Trim(string(by))
			jsonStr := strings.ReplaceAll(string(by), "\n", "")
			jsonStr = strings.ReplaceAll(jsonStr, " ", "")
			logrus.Info(jsonStr)
			if qtModels, err := GetQtModelsByTaskId(taskId); err == nil {
				qtModels.TaskId = taskId
				qtModels.SuggestScore = jsonStr
				err = UpdateQtModelsById(qtModels)
				if err !=nil {
					logrus.WithField("err", "watchAddSuggest更新推荐置信度").Error(err.Error())
				}
			} else {
				_, err = AddQtModels(&QtModels{
					TaskId:       taskId,
					ProjectId:    qtProject,
					SuggestScore: jsonStr,
					CreateTime:   time.Now(),
				})
				if err !=nil {
					logrus.WithField("err", "watchAddSuggest在这里新增推荐置信度").Error(err.Error())
				}
			}
		}
	}
	return err
}

func watchAddLabels(taskId string, labelPath string) (err error) {
	temp := strings.Split(strings.ReplaceAll(labelPath, beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"), ""), "/")
	isMultilabel := 0
	var labelStr string
	if qtProject, err := GetQtProjectsByProjectName(temp[0]); err == nil {
		//读取labels.names 如果不存在就直接不执行
		if by, err := afero.ReadFile(afero.NewOsFs(), labelPath); err == nil {
			temp := strings.Split(string(by), "\n")
			labelStr = "," + strings.ReplaceAll(string(by), "\n", ",") // 此处首位加了个,为了到时候筛选的时候过滤掉个别标签中包含搜索的关键词，到时候搜索直接,关键词,
			if !strings.HasSuffix(labelStr, ",") { // 需要确保以,结尾， 因为有些情况不存在最后一个空行
				labelStr += ","
			}
			for _, v := range temp {
				if !strings.EqualFold(v, "") {
					_, _ = AddQtLabels(&QtLabels{
						ProjectId:  qtProject,
						LabelName:  v,
						Remarks:    "自动新增",
						CreateTime: time.Now(),
					})
				}
			}
			if strings.Count(labelStr, ",") > 2 { // 因为上面行首新增了个,所以这里要大于2
				isMultilabel = 1
			}
			if qtModels, err := GetQtModelsByTaskId(taskId); err == nil {
				qtModels.TaskId = taskId
				qtModels.CreateTime = time.Now()
				qtModels.IsMultilabel = isMultilabel
				qtModels.LabelStr = labelStr
				err = UpdateQtModelsById(qtModels)
				if err !=nil {
					logrus.WithField("err", "watchlabels更新模型失败").Error(err.Error())
				}
			} else {
				_, _ = AddQtModels(&QtModels{
					TaskId:       taskId,
					IsMultilabel: isMultilabel, // 这里可能存在.names文件后面复制进来的情况，所以这些是空的，没事watchAddLabels会补齐这个
					LabelStr:     labelStr,     // 这里可能存在.names文件后面复制进来的情况，所以这些是空的，没事watchAddLabels会补齐这个
					CreateTime:   time.Now(),
				})
			}
		}
	}
	return err
}

func QTingTiny3l(taskId string, modelPath string) {
	temp := strings.Split(strings.ReplaceAll(modelPath, beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"), ""), "/")
	projectName := temp[0]
	aiFrameworkName := temp[2]
	if qtProject, err := GetQtProjectsByProjectName(projectName); err == nil {
		qtAiFramework, _ := GetQtAiFrameworkByFrameworkName(aiFrameworkName)
		var tempName = ""
		// 先查询是否存在相关的训练记录
		//if qtTrainRecord, err := GetQtTrainRecordByTaskId(taskId); err == nil {
		//	tempName = "(" + qtTrainRecord.TaskName + ")"
		//} else {
		//	tempName = "(轻蜓训练)"
		//	// 接下去插入记录
		//	_, err = AddQtTrainRecord(&QtTrainRecord{
		//		TaskId:        taskId,
		//		ContainerId:	"",
		//		TaskName:      "自动新增(轻蜓训练)",
		//		ProjectId:     qtProject,
		//		Status:        4,
		//		AiFrameworkId: qtAiFramework,
		//		AssetsType:    "",
		//		IsJump:        0,
		//		DrawUrl:       "",
		//		CreateTime:    time.Now(),
		//	})
		//	if err !=nil{
		//		logrus.Error(err.Error())
		//	}
		//}
		// 插入模型数据
		if qtModels, err := GetQtModelsByTaskId(taskId); err == nil {
			qtModels.TaskId = taskId
			qtModels.ProjectId = qtProject
			qtModels.AiFrameworkId = qtAiFramework
			qtModels.ModelName = time.Now().Format("2006-01-02 03:04") + tempName
			qtModels.ModelPath = modelPath
			qtModels.Status = 0
			qtModels.CreateTime = time.Now()
			err = UpdateQtModelsById(qtModels)
			if err != nil {
				logrus.WithField("err", "qtingsl更新模型失败").Error(err.Error())
			}
		} else {
			_, _ = AddQtModels(&QtModels{
				TaskId:        taskId,
				ProjectId:     qtProject,
				AiFrameworkId: qtAiFramework,
				ModelName:     time.Now().Format("2006-01-02 03:04") + tempName,
				ModelPath:     modelPath,
				Status:        0,
				CreateTime:    time.Now(),
			})
		}
	} else {
		logrus.Error(err)
	}
}