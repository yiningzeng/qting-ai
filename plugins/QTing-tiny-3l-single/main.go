package main
//package debugTest

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"qting-ai/models"
	"strings"
	"time"
)

var (
	pluginName    = "QTing-tiny-3l-single"
	pluginVersion = 11
)

type AiTrain struct {
	AiFrameworkId       int       `json:"aiFrameworkId"`
	ProjectId           int       `yaml:"projectId"` // 训练中心使用 项目名
	TaskId              string    `json:"taskId"`
	TaskName            string    `json:"taskName"`
	ProjectName         string    `json:"projectName"`
	AssetsType          string    `json:"assetsType"`
	ProviderType        string    `json:"providerType"`

	Batchsize           int       `json:"batchSize"`
	Imagewidth          int       `json:"imageWidth"`
	Imageheight         int       `json:"imageHeight"`
	Maxiter             int       `json:"maxIter"`
	Pretrainweight      string    `json:"pretrainWeight"`
	Gpus                string    `json:"gpus"`
	Triantype           int       `json:"trianType"`
	Singletrain         []string  `json:"singleTrain"`
	Angle               int       `json:"angle"`
	Cell_stride         int       `json:"cell_stride"`
	Cellsize            int       `json:"cellsize"`
	Expand_size         []int     `json:"expand_size"`
	Ignore_size         []int     `json:"ignore_size"`
	Resizearrange       []float32 `json:"resizearrange"`
	Trainwithnolabelpic int       `json:"trainwithnolabelpic"`
	Subdivisionssize    int       `json:"subdivisionssize"`
	Rmgeneratedata      int       `json:"rmgeneratedata"`
	Split_ratio         float32   `json:"split_ratio"`
	Recalldatum         int       `json:"recalldatum"`
	Otherlabeltraintype int       `json:"otherlabeltraintype"`
	Mergetrainsymbol    int       `json:"mergeTrainSymbol"`
	Learnrate           float32   `json:"learnrate"`
	Otherlabelstride    int       `json:"otherlabelstride"`
	Isshuffle           bool      `json:"isshuffle"`
}

type TrainConfig struct {
	ProjectId  int `yaml:"projectId"` // 训练中心使用 项目名
	ProjectName string `yaml:"projectName"` // 训练中心使用 项目名
	ProjectPath string `yaml:"projectPath"` // 训练中心使用 项目路径
	ProviderType string `yaml:"providerType"` // 训练中心使用 框架名称
	Image string `yaml:"image"` // 训练中心使用 框架的docker镜像
	Volume string `yaml:"volume"` // 训练中心使用 框架镜像内部映射的地址
	TaskId string `yaml:"taskId"` // 训练中心使用 项目标识

	Acc_suggest_pro []float32 `yaml:"acc_suggest_pro"`
	Anchortype int `yaml:"anchortype"`
	Angle int `yaml:"angle"`
	Batchsize int `yaml:"batchsize"`
	Cell_stride int `yaml:"cell_stride"`
	Cellsize int `yaml:"cellsize"`
	Expand_size []int `yaml:"expand_size"`
	F1_suggest_pro []float32 `yaml:"f1_suggest_pro"`
	Gpus string `yaml:"gpus"`
	Ignore_size []int `yaml:"ignore_size"`
	Imagesize []int `yaml:"imagesize"`
	Labelsfilename string `yaml:"labelsfilename"`
	Lablelist []string `yaml:"lablelist"`
	Maxiter int `yaml:"maxiter"`
	MergeTrainSymbol int `yaml:"mergeTrainSymbol"`
	Modeldcfgname string `yaml:"modeldcfgname"`
	Modelname string `yaml:"modelname"`
	Otherlabeltraintype int `yaml:"otherlabeltraintype"`
	Path_prefix string `yaml:"path_prefix"`
	Pixelareafpthresh []float32 `yaml:"pixelareafpthresh"`
	Pretraincfgname string `yaml:"pretraincfgname"`
	Pretrainweight string `yaml:"pretrainweight"`
	Random int `yaml:"random"`
	Recall_suggest_pro []float32 `yaml:"recall_suggest_pro"`
	Recalldatum int `yaml:"recalldatum"`
	Resizearrange []float32 `yaml:"resizearrange"`
	Rmgeneratedata int `yaml:"rmgeneratedata"`
	Scale int `yaml:"scale"`
	Singletrain []string `yaml:"singletrain"`
	Split_ratio float32 `yaml:"split_ratio"`
	Subdivisionssize int `yaml:"subdivisionssize"`
	Sum_suggest_pro []float32 `yaml:"sum_suggest_pro"`
	Trainwithnolabelpic int `yaml:"trainwithnolabelpic"`
	Triantype int `yaml:"triantype"`
	Learnrate float32 `yaml:"learnrate"`
	Otherlabelstride int `yaml:"otherlabelstride"`
	Isshuffle bool `yaml:"isshuffle"`
}

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


func DoTrain(m *AiTrain) (err error) {
	fs := afero.NewOsFs()
	if b, err := afero.ReadFile(fs, "conf/" + m.ProviderType + ".yaml"); err == nil {
		var temp TrainConfig
		if err = yaml.Unmarshal(b, &temp); err == nil {
			// region 第一步 先查询相关数据 并更新数据
			aiFramwork, err := models.GetQtAiFrameworkById(m.AiFrameworkId)
			if err != nil {return err}

			temp.Modeldcfgname = aiFramwork.Cfg
			temp.Image = fmt.Sprintf("%s:%s", aiFramwork.BaseImageUrl, aiFramwork.ImageVersion)
			temp.Volume = aiFramwork.Volume

			project, err := models.GetQtProjectsById(m.ProjectId)
			if err != nil {return err}
			temp.ProjectPath = project.AssetsPath
			// endregion
			// region 第二步 更新数据 发布到队列
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
				return err
			}
			savePath := beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/") + temp.ProjectName + "/training_data/"
			statusFile := fmt.Sprintf("%s/train_status_%s.log", savePath, temp.TaskId)

			logrus.WithField("taskId", temp.TaskId).Info("开始训练")
			if err = fs.MkdirAll(savePath, 0755); err == nil {

				if err = afero.WriteFile(fs, statusFile, []byte(models.Waiting), 0755); err == nil {
					err = afero.WriteFile(fs, savePath+"config.yaml", data, 0755)
					//go tools.GetMsg(true)
				} else {
					logrus.WithField("taskId", temp.TaskId).Error("训练失败")
				}
			}

			models.Publish(string(data))
			// endregion
			// region 第三步 插入数据到数据库
			record := &models.QtTrainRecord{
				TaskId:        m.TaskId,
				TaskName:      m.TaskName,
				ProjectId:     &models.QtProjects{Id: m.ProjectId},
				Status:        models.WaitingInt,
				AiFrameworkId: aiFramwork,
				AssetsType:    m.AssetsType,
				IsJump:        0,
				DrawUrl: fmt.Sprintf("..%s/%s/training_data/train_%s/chart.png",
					beego.AppConfig.DefaultString("ProjectPathStaticDir", "/qting"),
					project.ProjectName, m.TaskId),
				CreateTime: time.Now(),
			}
			_, err = models.AddQtTrainRecord(record)
			if err != nil {return err}
			if beego.AppConfig.DefaultBool("saveRabbitmqInfo", false) {
				qqtTrainRecord, _ := models.GetQtTrainRecordByTaskId(m.TaskId)
				_, _ = models.AddQtRabbitmqInfo(&models.QtRabbitmqInfo{
					RecordId: qqtTrainRecord,
					Queue:    beego.AppConfig.DefaultString("queue_name", "ai.train.topic-queue"),
					Message:  string(data),
				})
			}
			// endregion
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

// 框架开始训练执行方法
func Train(args ...interface{}) (ret interface{}, err error) {
	trainStr := args[0].(string)
	var v AiTrain
	if err = json.Unmarshal([]byte(trainStr), &v); err == nil {
		err = DoTrain(&v)
	}
	return nil, nil
}

// 插件通用的执行方法
func Done(args ...interface{}) (ret interface{}, err error) {
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

func main() {
	if _, err := os.Stat(`/qtingvisionfolder/Projects/前道/backup/QTing-tiny-3l-multilabel/20210218152531.weights`); err == nil {
		logrus.Info("啦啦啦")
	}
	//str := "{\"aiFrameworkId\":1,\"taskId\":\"20210623174205\",\"taskName\":\"啊啊所多\",\"projectId\":1,\"assetsType\":\"powerAi\",\"providerType\":\"QTing-tiny-3l-single\",\"singleTrain\":[\"luhou-zhanxi\",\"xiqiu\",\"heidian\",\"aasds\",\"zhanxi\",\"op\",\"o2p\",\"o2p2\",\"xizhu\",\"dianzhuang-yiwu\",\"yanghua\",\"asdsadsddd\",\"消息\",\"可爱\",\"可爱2\",\"阿萨德\",\"二套\",\"让他\",\"阿萨德3\",\"阿萨德36\",\"小可爱\",\"战三\",\"李四\",\"小吴\"],\"mergeTrainSymbol\":0,\"learnrate\":0.00261,\"recalldatum\":2,\"otherlabeltraintype\":1,\"batchSize\":64,\"subdivisionssize\":2,\"imageWidth\":-1,\"angle\":360,\"split_ratio\":0.95,\"maxIter\":-1,\"trainwithnolabelpic\":50000,\"cell_stride\":1,\"otherlabelstride\":1,\"cellsize\":16,\"learnratestepsratio\":[0.9,0.95],\"expand_size\":[8,8],\"ignore_size\":[6,6],\"resizearrange\":[0.3,1.6],\"gpus\":\"0,1\",\"trianType\":0,\"rmgeneratedata\": 0,\"isshuffle\":true}"

	//var trainConfig TrainConfigTest
	//if err := yaml.Unmarshal([]byte(str), &trainConfig); err == nil {
	//	_, _ = Train(str)
	//}

}

