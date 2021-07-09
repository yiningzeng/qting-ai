package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"path"
	"qting-ai/plugins"
	"strings"
)

func getLabelStr(taskId string, modelBasePath string) (res string) {
	labelFileName := path.Join(modelBasePath, taskId + ".names")
	if by, err := afero.ReadFile(afero.NewOsFs(), labelFileName); err == nil {
		labelStr := "," + strings.ReplaceAll(string(by), "\n", ",") // 此处首位加了个,为了到时候筛选的时候过滤掉个别标签中包含搜索的关键词，到时候搜索直接,关键词,
		if !strings.HasSuffix(labelStr, ",") { // 需要确保以,结尾， 因为有些情况不存在最后一个空行
			labelStr += ","
		}
		return labelStr
	}
	return ""
}
func getSuggest(taskId string, modelBasePath string) (res string) {
	//读取labels.names 如果不存在就直接不执行
	suggestFileName :=  path.Join(modelBasePath, taskId + ".suggest")
	if by, err := afero.ReadFile(afero.NewOsFs(), suggestFileName); err == nil {
		//temp := strings.Trim(string(by))
		jsonStr := strings.ReplaceAll(string(by), "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, " ", "")
		logrus.Info(jsonStr)
		return jsonStr
	}
	return ""
}

func DoTrain(project *plugins.QtProjects, aiFramework *plugins.QtAiFramework, m *plugins.AiTrain) (data string, err error) {
	fs := afero.NewOsFs()
	if b, err := afero.ReadFile(fs, "conf/" + m.ProviderType + ".yaml"); err == nil {
		var temp plugins.TrainConfig
		if err = yaml.Unmarshal(b, &temp); err == nil {
			// region 第一步 先查询相关数据 并更新数据
			if err != nil {return "", err}

			temp.Modeldcfgname = aiFramework.Cfg
			temp.Image = fmt.Sprintf("%s:%s", aiFramework.BaseImageUrl, aiFramework.ImageVersion)
			temp.Volume = aiFramework.Volume
			temp.ProjectPath = project.AssetsPath
			// endregion
			// region 第二步 更新数据
			temp.Batchsize =  m.Batchsize
			temp.Maxiter =  m.Maxiter
			temp.Imagesize =  []int {m.Imagewidth, m.Imageheight}
			temp.Modelname =  m.TaskId
			temp.Triantype =  m.Triantype
			temp.Pretrainweight =  m.Pretrainweight
			temp.Pretraincfgname = m.Pretrainweight
			temp.Gpus =  m.Gpus
			temp.Singletrain =  m.Singletrain[:]
			temp.Angle =  m.Angle
			temp.Cell_stride =  m.Cell_stride
			temp.Cellsize =  m.Cellsize
			temp.Expand_size =  m.Expand_size[:]
			temp.Ignore_size =  m.Ignore_size[:]
			temp.Resizearrange =  m.Resizearrange[:]
			temp.Trainwithnolabelpic =  m.Trainwithnolabelpic
			temp.Subdivisionssize =  m.Subdivisionssize
			temp.Rmgeneratedata =  m.Rmgeneratedata
			temp.Split_ratio =  m.Split_ratio
			temp.Recalldatum =  m.Recalldatum
			temp.Otherlabeltraintype =  m.Otherlabeltraintype
			temp.MergeTrainSymbol =  m.Mergetrainsymbol

			temp.ProjectName = project.ProjectName
			temp.ProjectId = project.Id
			temp.TaskId = m.TaskId
			temp.ProviderType = m.ProviderType

			temp.Isshuffle = m.Isshuffle
			temp.Otherlabelstride = m.Otherlabelstride
			temp.Learnrate = m.Learnrate
			data, err := yaml.Marshal(temp)
			if err != nil {
				return "", err
			}
			// endregion
			return string(data), nil
		} else {
			return "", errors.Wrap(err, "yaml Unmarshal failed")
		}
	} else {
		return "", errors.Wrap(err, fmt.Sprintf("can't find %s file", "conf/" + m.ProviderType + ".yaml"))
	}
}
