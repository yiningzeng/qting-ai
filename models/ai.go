package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"os/exec"
	"strings"
	"time"
)

const (
	Failed   = "训练失败"
	Ready    = "准备完成"
	Waiting  = "等待训练"
	Training = "正在训练"
	Stopped  = "停止训练"
	Done     = "训练完成"
	Unknown  = "未知状态"
)
const (
	FailedInt   = -1
	ReadyInt    = 0
	WaitingInt  = 1
	TrainingInt = 2
	StoppedInt  = 3
	DoneInt     = 4
	UnknownInt  = 110
)

type AiTrain struct {
	AiFrameworkId       int       `json:"aiFrameworkId"`
	ProjectId           int       `yaml:"projectId"` // 训练中心使用 项目名
	Taskid              string    `json:"taskId"`
	Taskname            string    `json:"taskName"`
	Projectname         string    `json:"projectName"`
	Assetstype          string    `json:"assetsType"`
	Providertype        string    `json:"providerType"`
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

type QTingTiny3L struct {
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

// region Aoi获取的接口
type OneAoiProject struct {
	ProjectName  string        `json:"project_name"`
	List        []OneNetFramework `json:"list"`
}
type OneNetFramework struct {
	NetFramework string      `json:"net_framework"`
	Models       []string `json:"models"`
}
// endregion

func V1GetModelsForAoi() (ml []interface{}, err error) {
	var projects []QtProjects
	var ais []QtAiFramework
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtProjects))
	qa := o.QueryTable(new(QtAiFramework))
	if _, err = qa.All(&ais); err == nil { // 获取所有的ai框架
		if _, err = qs.All(&projects); err == nil {
			for _, oneProject := range projects {
				oneAoiProject := OneAoiProject{
					ProjectName: oneProject.ProjectName,
					List:        []OneNetFramework{},
				}
				for _, oneAiF := range ais {
					oneNetFramework := OneNetFramework{
						NetFramework: oneAiF.FrameworkName,
						Models:       []string{},
					}
					var tempModels []QtModels
					qModel := o.QueryTable(new(QtModels))
					// 再查找相关的模型
					qModel = qModel.Filter("ProjectId", oneProject).Filter("AiFrameworkId", oneAiF).Filter("Status", 2)
					if _, err = qModel.All(&tempModels); err == nil {
						// 这里查询的就是 oneProject 和 oneAiF的模型了
						for _, oneModel := range tempModels {
							oneNetFramework.Models = append(oneNetFramework.Models, oneModel.PublishUrl)
						}
					}
					oneAoiProject.List = append(oneAoiProject.List, oneNetFramework)
				}
				ml = append(ml, oneAoiProject)
			}
			return ml, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

var c = cron.New(cron.WithSeconds())

func DoTrain(m *AiTrain) (err error) {
	fs := afero.NewOsFs()
	if b, err := afero.ReadFile(fs, "conf/" + m.Providertype + ".yaml"); err == nil {
		var temp QTingTiny3L
		if err = yaml.Unmarshal(b, &temp); err == nil {
			// region 第一步 先查询相关数据 并更新数据
			aiFramwork, err := GetQtAiFrameworkById(m.AiFrameworkId)
			if err != nil {return err}

			temp.Modeldcfgname = aiFramwork.Cfg
			temp.Image = fmt.Sprintf("%s:%s", aiFramwork.BaseImageUrl, aiFramwork.ImageVersion)
			temp.Volume = aiFramwork.Volume

			project, err := GetQtProjectsById(m.ProjectId)
			if err != nil {return err}
			temp.ProjectPath = project.AssetsPath
			// endregion
			// region 第二步 更新数据 发布到队列
			temp.Batchsize =  m.Batchsize
			temp.Maxiter =  m.Maxiter
			temp.Imagesize =  []int {m.Imagewidth, m.Imageheight}
			temp.Modelname =  m.Taskid
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
			temp.TaskId = m.Taskid
			temp.ProviderType = m.Providertype

			temp.Isshuffle = m.Isshuffle
			temp.Otherlabelstride = m.Otherlabelstride
			temp.Learnrate = m.Learnrate

			data, _ := yaml.Marshal(temp)
			Publish(string(data))
			// endregion
			// region 第三步 插入数据到数据库
			record := &QtTrainRecord{
				TaskId:        m.Taskid,
				TaskName:      m.Taskname,
				ProjectId:     &QtProjects{Id: m.ProjectId},
				Status:        WaitingInt,
				AiFrameworkId: aiFramwork,
				AssetsType:    m.Assetstype,
				IsJump:        0,
				DrawUrl: fmt.Sprintf("%s/%s/training_data/train_%s/chart.png",
					beego.AppConfig.DefaultString("drawBaseUrl", "http://localhost:1121/"),
					project.ProjectName, m.Taskid),
				CreateTime: time.Now(),
			}
			_, err = AddQtTrainRecord(record)
			if err != nil {return err}
			if beego.AppConfig.DefaultBool("saveRabbitmqInfo", false) {
				qqtTrainRecord, _ := GetQtTrainRecordByTaskId(m.Taskid)
				_, _ = AddQtRabbitmqInfo(&QtRabbitmqInfo{
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

func StopTrain(taskId string, status string, statusCode int) (err error) {
	if data, err := GetQtTrainRecordByTaskId(taskId); err == nil {
		data.Status = statusCode
		savePath := beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/") + data.ProjectId.ProjectName + "/training_data/"
		fs := afero.NewOsFs()
		statusFile := fmt.Sprintf("%s/train_status_%s.log", savePath, taskId)
		updateStatus(fs, statusFile, status, statusCode, taskId, "")
	}
	return err
}

func updateStatus(fs afero.Fs, statusFile string, status string, statusCode int, taskId string, dockerId string) {
	//_ = afero.WriteFile(fs, "/aaa.txt", []byte(Training), 0755)
	if err := afero.WriteFile(fs, statusFile, []byte(status), 0755); err == nil {
		if data, err := GetQtTrainRecordByTaskId(taskId); err == nil {
			data.Status = statusCode
			data.ContainerId = dockerId
			if err = UpdateQtTrainRecordById(data); err == nil {
				if statusCode == StoppedInt || statusCode == DoneInt || statusCode == FailedInt || statusCode == UnknownInt {
					go GetMsg(true)
				}
			} else {
				logrus.WithField("errType", "更新记录状态出错").Error(err.Error())
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"errType": "无法通过taskId找到记录",
				"taskId": taskId,
				"statusCode": statusCode,
			}).Error(err.Error())
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"statusFile": statusFile,
			"status":     status,
			"taskId":     taskId,
		}).Error(err.Error())
	}


}

func cronFunc() {
	if ok, msgBody := GetMsg(false); ok == true {
		var trainConfig QTingTiny3L
		if err := yaml.Unmarshal(msgBody, &trainConfig); err == nil {
			if data, err := yaml.Marshal(trainConfig); err == nil {
				//logrus.Debug(string(data))
				savePath := beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/") + trainConfig.ProjectName + "/training_data/"
				fs := afero.NewOsFs()
				statusFile := fmt.Sprintf("%s/train_status_%s.log", savePath, trainConfig.TaskId)
				if isEx, err := afero.Exists(fs, statusFile); err == nil {
					if isEx {
						byteStatus, _ := afero.ReadFile(fs, statusFile)
						//logrus.Debug(string(byteStatus))
						switch strings.ReplaceAll(strings.ReplaceAll(string(byteStatus), "\n", ""), " ", "") {
						case Waiting:
							cmdStr := fmt.Sprintf("%s --name 'qtingTrain-%s' -v /etc/localtime:/etc/localtime:ro -v %s:%s %s",
								beego.AppConfig.DefaultString("dockerRunPrefix", "docker run --shm-size 32G --memory-swap -1 --gpus all --rm -d"),
								trainConfig.TaskId, trainConfig.ProjectPath, trainConfig.Volume, trainConfig.Image)
							cmd := exec.Command("bash", "-c", cmdStr)
							if output, err := cmd.CombinedOutput(); err == nil {
								updateStatus(fs, statusFile, Training, TrainingInt, trainConfig.TaskId, string(output))
							} else {
								logrus.Error(cmdStr)
								logrus.WithField("out",  string(output)).Error(err.Error())
							}
							break
						case Training:
							// 检测一下容器是否停止, 如果停止了,判断下是否是训练完成，如果不是训练完成那么就是未知状态
							cmd := exec.Command("bash", "-c", fmt.Sprintf("docker ps |grep %s", trainConfig.TaskId))
							output, err := cmd.CombinedOutput()
							if err != nil {
								if strings.Contains(err.Error(), "exit status 1") {
									logrus.WithFields(logrus.Fields{
										"cmd": cmd,
										"output": string(output),
									}).Info("正在训练 > 未查询到正在运行的训练容器")
									if len(string(output)) < 10 {
										if by, _ := afero.ReadFile(fs, statusFile); string(by) != Done {
											// 未训练完成
											updateStatus(fs, statusFile, Failed, FailedInt, trainConfig.TaskId, "")
										}
									}
								}
							}
							//updateStatus(fs, statusFile, Training, TrainingInt, &trainConfig)
							break
						case Done:
							updateStatus(fs, statusFile, Done, DoneInt, trainConfig.TaskId, "")
							break
						case Failed:
							updateStatus(fs, statusFile, Failed, FailedInt, trainConfig.TaskId, "")
							break
						case Stopped:
							updateStatus(fs, statusFile, Stopped, StoppedInt, trainConfig.TaskId, "")
							break
						default:
							updateStatus(fs, statusFile, Unknown, UnknownInt, trainConfig.TaskId, "")
							break
						}
					} else { // 状态文件不存在
						if err = fs.MkdirAll(savePath, 0755); err == nil {
							if err = afero.WriteFile(fs, savePath+"config.yaml", data, 0755); err == nil {
								_ = afero.WriteFile(fs, statusFile, []byte(Waiting), 0755)
								//go tools.GetMsg(true)
							}
						}
					}
				}
			}
		} else {
			logrus.Error("队列训练数据有误->需要处理执行出错")
			go GetMsg(true)
		}
	} else {
		//logrus.Debug("队列空数据")
	}
	//logrus.Debug(time.Now())
}


func StartCron() (err error) {
	_, err = c.AddFunc(beego.AppConfig.DefaultString("cron", "0/5 * * * * ?"), cronFunc)
	if err == nil {
		c.Start()
		return nil
	} else {
		return err
	}
}

func StopCron() {
	c.Stop()
}
