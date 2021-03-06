package main

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"qting-ai/plugins"
	"time"
)

var (
	pluginName    = "QTing-tiny-3l-single"
	pluginVersion = 20
)

func Version(args ...interface{}) (ret interface{}, err error) {
	t := plugins.QtPlugins {
		Module:      pluginName,
		VersionCode: pluginVersion,
		Symbol:      "Run",
		Args:        "taskId string, projectName string, aiFrameworkName string",
		PluginName:  args[0].(string),
		CreateTime:  time.Now(),
	}
	res, err := json.Marshal(&t)
	return string(res), err
}
// 框架开始训练执行方法
// ret 默认是返回队列的数据字符串
/*
args[0] 项目结构体Json字符串
args[1] 框架结构体Json字符串
args[2] 训练参数字符串
*/
func Train(args ...interface{}) (ret interface{}, err error) {
	logrus.WithFields(logrus.Fields{
		"pluginName":    pluginName,
		"pluginVersion": pluginVersion,
	}).Info("插件调用成功...")
	var qtProject plugins.QtProjects; err = json.Unmarshal([]byte(args[0].(string)), &qtProject); if err != nil {return nil, errors.Wrap(err, "Unmarshal failed QtProjects")}
	var qtAiFramework plugins.QtAiFramework; err = json.Unmarshal([]byte(args[1].(string)), &qtAiFramework); if err != nil {return nil, errors.Wrap(err, "Unmarshal failed QtAiFramework")}
	trainStr := args[2].(string)
	var v plugins.AiTrain
	if err = json.Unmarshal([]byte(trainStr), &v); err == nil {
		return DoTrain(&qtProject, &qtAiFramework, &v)
	} else {
		return nil, errors.Wrap(err, "Unmarshal failed AiTrain")
	}
}

// 训练完成后通用的执行方法
/*
args[0] taskId
args[1] 项目结构体Json字符串
args[2] 框架结构体Json字符串
args[3] 训练任务记录结构体Json字符串
ret 返回[]qtModels数组的Json数据
*/
func Done(args ...interface{}) (ret interface{}, err error) {
	logrus.WithFields(logrus.Fields{
		"pluginName":    pluginName,
		"pluginVersion": pluginVersion,
	}).Info("插件调用成功...")
	taskId := args[0].(string)
	var qtProject plugins.QtProjects; err = json.Unmarshal([]byte(args[1].(string)), &qtProject); if err != nil {return nil, errors.Wrap(err, "Unmarshal failed QtProjects")}
	var qtAiFramework plugins.QtAiFramework; err = json.Unmarshal([]byte(args[2].(string)), &qtAiFramework); if err != nil {return nil, errors.Wrap(err, "Unmarshal failed QtAiFramework")}
	var qtTrainRecord plugins.QtTrainRecord; errQtTrainRecord := json.Unmarshal([]byte(args[3].(string)), &qtTrainRecord)
	modelBasePath := path.Join(beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"),
		qtProject.ProjectName,
		"backup",
		qtAiFramework.FrameworkName)

	var models []plugins.QtModels

	for _, v := range qtProject.Labels {
		modelLabelBasePath := path.Join(modelBasePath, v.LabelName)
		modelPath := path.Join(modelLabelBasePath, fmt.Sprintf("%s-%s.weights", taskId, v.LabelName))
		if _, err := os.Stat(modelPath); err == nil {
			suggest := getSuggest(taskId, v.LabelName, modelLabelBasePath)
			labelStr := getLabelStr(taskId, v.LabelName, modelLabelBasePath)
			var tempName = ""
			// 先查询是否存在相关的训练记录
			if errQtTrainRecord == nil {
				tempName = "(" + qtTrainRecord.TaskName + ")"
			}
			qtModel := &plugins.QtModels{
				TaskId:        taskId,
				ProjectId:     &qtProject,
				AiFrameworkId: &qtAiFramework,
				ModelName:     time.Now().Format("2006-01-02 03:04") + tempName,
				ModelPath:     modelPath,
				ModelBasePath: modelLabelBasePath,
				SuggestScore:  suggest,
				LabelStr:      labelStr,
				IsMultilabel:  0,
				Status:        0,
				CreateTime:    time.Now(),
			}
			models = append(models, *qtModel)
		}
	}
	qtModelsByte, err := json.Marshal(&models)
	//_, err = models.AddQtModels(qtModels)
	if err != nil {
		logrus.WithField("taskId", taskId).Error(fmt.Sprintf("%+v", errors.Wrap(err, "模型转换json失败")))
		return nil, err
	}
	return string(qtModelsByte), nil
}
