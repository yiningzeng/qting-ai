package main

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"path"
	"qting-ai/models"
	"strings"
	"time"
)

var (
	pluginName    = "QTing-tiny-3l-single"
	pluginVersion = 10
)

func Version(args ...interface{}) (ret interface{}, err error) {
	t := models.QtPlugins{
		Module:      pluginName,
		VersionCode: pluginVersion,
		Symbol:      "Run",
		Args:        "taskId string, projectName string, aiFrameworkName string",
		PluginName:  args[0].(string),
		CreateTime:  time.Now(),
	}
	return t, nil
}
func getLabelStr(taskId string, labelName string, modelBasePath string) (res string) {
	labelFileName := path.Join(modelBasePath, fmt.Sprintf("%s-%s.names", taskId, labelName))
	if by, err := afero.ReadFile(afero.NewOsFs(), labelFileName); err == nil {
		labelStr := "," + strings.ReplaceAll(string(by), "\n", ",") // 此处首位加了个,为了到时候筛选的时候过滤掉个别标签中包含搜索的关键词，到时候搜索直接,关键词,
		if !strings.HasSuffix(labelStr, ",") { // 需要确保以,结尾， 因为有些情况不存在最后一个空行
			labelStr += ","
		}
		return labelStr
	}
	return ""
}
func getSuggest(taskId string, labelName string, modelBasePath string) (res string) {
	//读取labels.names 如果不存在就直接不执行
	suggestFileName := path.Join(modelBasePath, fmt.Sprintf("%s-%s.suggest", taskId, labelName))
	if by, err := afero.ReadFile(afero.NewOsFs(), suggestFileName); err == nil {
		//temp := strings.Trim(string(by))
		jsonStr := strings.ReplaceAll(string(by), "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, " ", "")
		logrus.Info(jsonStr)
		return jsonStr
	}
	return ""
}

// 插件通用的执行方法
func Run(args ...interface{}) (ret interface{}, err error) {
	taskId := args[0].(string)
	projectName := args[1].(string)
	aiFrameworkName := args[2].(string)
	modelBasePath := path.Join(beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"),
		projectName,
		"backup",
		aiFrameworkName)

	logrus.WithFields(logrus.Fields{
		"taskId": taskId,
		"projectName": projectName,
		"aiFrameworkName": aiFrameworkName,
		"modelBasePath": modelBasePath,
	}).Debug("QTing-tiny-3l-single参数")
	if qtProject, err := models.GetQtProjectsByProjectName(projectName); err == nil {
		for _, v := range qtProject.Labels {
			modelLabelBasePath := path.Join(modelBasePath, v.LabelName)
			modelPath := path.Join(modelLabelBasePath, fmt.Sprintf("%s-%s.weights", taskId, v.LabelName))
			if _, err := os.Stat(modelPath); err == nil {
				suggest := getSuggest(taskId, v.LabelName, modelLabelBasePath)
				labelStr := getLabelStr(taskId, v.LabelName, modelLabelBasePath)
				qtAiFramework, _ := models.GetQtAiFrameworkByFrameworkName(aiFrameworkName)
				var tempName = ""
				// 先查询是否存在相关的训练记录
				if qtTrainRecord, err := models.GetQtTrainRecordByTaskId(taskId); err == nil {
					tempName = "(" + qtTrainRecord.TaskName + ")"
				}

				qtModels := &models.QtModels{
					TaskId:        taskId,
					ProjectId:     qtProject,
					AiFrameworkId: qtAiFramework,
					ModelName:     time.Now().Format("2006-01-02 03:04") + tempName,
					ModelPath:     modelPath,
					ModelBasePath: modelLabelBasePath,
					SuggestScore:  suggest,
					LabelStr:      labelStr,
					IsMultilabel:  0,
					Status:        0,
					CreateTime:    time.Now(),
				}
				_, err = models.AddQtModels(qtModels)
				if err != nil {
					logrus.WithField("taskId", taskId).Error(fmt.Sprintf("%+v", errors.Wrap(err, "新增模型失败")))
				}
			}
		}
		return nil, nil
	} else {
		logrus.WithField("ProjectName", projectName).Error(fmt.Sprintf("%+v", errors.Wrap(err, "未找到项目")))
		return nil, err
	}
}
