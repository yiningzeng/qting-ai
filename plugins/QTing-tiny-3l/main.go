package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/sirupsen/logrus"
	"qting-ai/models"
	"strings"
	"time"
)

var (
	pluginName    = "QTing-tiny-3l"
	pluginVersion = 10
)

func Version(args ...interface{}) (ret interface{}, err error) {
	t := models.QtPlugins{
		Module:      pluginName,
		VersionCode: pluginVersion,
		Symbol:      "Run",
		Args:        "taskId string, modelPath string",
		PluginName:  args[0].(string),
		CreateTime:  time.Now(),
	}
	return t, nil
}

// 插件通用的执行方法
func Run(taskId string, modelPath string) (err error) {
	temp := strings.Split(strings.ReplaceAll(modelPath, beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"), ""), "/")
	projectName := temp[0]
	aiFrameworkName := temp[2]
	if qtProject, err := models.GetQtProjectsByProjectName(projectName); err == nil {
		qtAiFramework, _ := models.GetQtAiFrameworkByFrameworkName(aiFrameworkName)
		var tempName = ""
		// 先查询是否存在相关的训练记录
		taskIdTemp := taskId
		if strings.Contains(taskId, "-") {
			taskIdTemp = strings.Split(taskId, "-")[0]
		}
		if qtTrainRecord, err := models.GetQtTrainRecordByTaskId(taskIdTemp); err == nil {
			tempName = "(" + qtTrainRecord.TaskName + ")"
		} else {
			tempName = "(轻蜓训练)"
			// 接下去插入记录
			_, err = models.AddQtTrainRecord(&models.QtTrainRecord{
				TaskId:        taskIdTemp,
				ContainerId:	"",
				TaskName:      "自动新增-轻蜓训练",
				ProjectId:     qtProject,
				Status:        4,
				AiFrameworkId: qtAiFramework,
				AssetsType:    "",
				IsJump:        0,
				DrawUrl:       "",
				CreateTime:    time.Now(),
			})
			if err !=nil{
				logrus.Error(err.Error())
			}
		}
		// 插入模型数据
		if qtModels, err := models.GetQtModelsByTaskId(taskId); err == nil {
			qtModels.TaskId = taskId
			qtModels.ProjectId = qtProject
			qtModels.AiFrameworkId = qtAiFramework
			qtModels.ModelName = time.Now().Format("2006-01-02 03:04") + tempName
			qtModels.ModelPath = modelPath
			qtModels.Status = 0
			qtModels.CreateTime = time.Now()
			err = models.UpdateQtModelsById(qtModels)
			if err != nil {
				logrus.WithField("err", "qtingsl更新模型失败").Error(err.Error())
			}
		} else {
			_, err = models.AddQtModels(&models.QtModels{
				TaskId:        taskId,
				ProjectId:     qtProject,
				AiFrameworkId: qtAiFramework,
				ModelName:     time.Now().Format("2006-01-02 03:04") + tempName,
				ModelPath:     modelPath,
				Status:        0,
				CreateTime:    time.Now(),
			})
			if err != nil {
				logrus.WithField("err", "新增模型失败").Error(err.Error())
			}
		}
	} else {
		logrus.Error(err)
	}
	return nil
}
