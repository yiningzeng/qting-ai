package plugins

import (
	"time"
)
type QtModels struct {
	Id            int            `json:"Id" orm:"column(id);auto"`
	TaskId        string         `json:"TaskId" orm:"column(task_id);size(100);null" description:"主要是方便后期更新这个模型通过taskId来查询"`
	ProjectId     *QtProjects    `json:"ProjectId" orm:"column(project_id);rel(fk);null" description:"项目名"`
	AiFrameworkId *QtAiFramework `json:"AiFrameworkId" orm:"column(ai_framework_id);rel(fk);null" description:"训练使用的框架ID"`
	ModelName     string         `json:"ModelName" orm:"column(model_name);size(100);null" description:"模型名称"`
	ModelPath     string         `json:"ModelPath" orm:"column(model_path);size(255);null" description:"模型的本地路径"`
	ModelBasePath string         `json:"ModelBasePath" orm:"column(model_base_path);size(255);null" description:"模型的保存目录"`
	IsMultilabel  int            `json:"IsMultilabel" orm:"column(is_multilabel);null" description:"该模型是否包含多标签 -1 表示未知 0 表示单标签单模型 1 表示多标签单模型"`
	LabelStr      string         `json:"LabelStr" orm:"column(label_str);null" description:"缺陷的标签，使用英文的,分离 比如zhanxi,jinmian,wuran"`
	SuggestScore  string         `json:"SuggestScore" orm:"column(suggest_score);null" description:"推荐置信度的json值"`
	Status        int            `json:"Status" orm:"column(status);null" description:"模型状态 0 表示未发布 1 表示已发布 2 表示已下线"`
	PublishUrl    string         `json:"PublishUrl" orm:"column(publish_url);size(255);null" description:"模型的发布地址"`
	PublishTime   time.Time      `json:"PublishTime" orm:"column(publish_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
	CreateTime    time.Time      `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}
type QtTrainRecord struct {
	Id            int            `json:"Id" orm:"column(id);auto"`
	TaskId        string         `json:"TaskId" orm:"column(task_id);size(30);null" description:"任务标识"`
	TaskName      string         `json:"TaskName" orm:"column(task_name);size(100);null" description:"训练任务名称"`
	ContainerId   string         `json:"ContainerId" orm:"column(container_id);size(100);null" description:"训练的容器的id"`
	ProjectId     *QtProjects    `json:"ProjectId" orm:"column(project_id);rel(fk)" description:"项目名"`
	Status        int            `json:"Status" orm:"column(status);null" description:"0:准备完成 1:等待训练 2:正在训练 3:暂停训练 4:训练完成 -1:训练有误"`
	AiFrameworkId *QtAiFramework `json:"AiFrameworkId" orm:"column(ai_framework_id);rel(fk)" description:"训练使用的框架ID"`
	AssetsType    string         `json:"AssetsType" orm:"column(assets_type);size(50);null" description:"训练的数据类型"`
	IsJump        int            `json:"IsJump" orm:"column(is_jump);null" description:"是否插队 0:正常 1:插队"`
	DrawUrl       string         `json:"DrawUrl" orm:"column(draw_url);size(255);null" description:"画图的网址"`
	CreateTime    time.Time      `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}

type QtLabels struct {
	Id         int         `json:"Id" orm:"column(id);auto"`
	ProjectId  *QtProjects `json:"ProjectId" orm:"column(project_id);rel(fk)"`
	LabelName  string      `json:"LabelName" orm:"column(label_name);size(30);null" description:"标签名称"`
	Remarks    string      `json:"Remarks" orm:"column(remarks);size(255);null"`
	CreateTime time.Time   `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}
type QtProjects struct {
	Id          int              `json:"Id" orm:"column(id);auto"`
	ProjectName string           `json:"ProjectName" orm:"column(project_name);size(100)" description:"项目名称"`
	AssetsPath  string           `json:"AssetsPath" orm:"column(assets_path);size(255);null" description:"项目的目录默认所有的项目目录都存在/qtingvisionfolder/Projects/"`
	ImageWidth  int              `json:"ImageWidth" orm:"column(image_width);null" description:"默认指的是发布时的图像宽"`
	ImageHeight int              `json:"ImageHeight" orm:"column(image_height);null" description:"默认指的是发布时的图像高"`
	Remarks     string           `json:"Remarks" orm:"column(remarks);size(100);null" description:"备注"`
	Labels      []*QtLabels      `json:"Labels" orm:"reverse(many)"` // 设置一对多的反向关系
	CreateTime  time.Time        `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}
type QtAiFramework struct {
	Id            int       `json:"Id" orm:"column(id);auto"`
	FrameworkName string    `json:"FrameworkName" orm:"column(framework_name);size(30)" description:"ai框架名称"`
	BaseImageUrl  string    `json:"BaseImageUrl" orm:"column(base_image_url);size(255)" description:"基础镜像地址"`
	ImageVersion  string    `json:"ImageVersion" orm:"column(image_version);size(255)" description:"镜像版本"`
	Volume        string    `json:"Volume" orm:"column(volume);size(255)" description:"数据映射的镜像内部地址"`
	Cfg           string    `json:"Cfg" orm:"column(cfg);size(255)" description:"对应的配置文件"`
	Remarks       string    `json:"Remarks" orm:"column(remarks);size(100);null" description:"备注"`
	ParsJson      string    `json:"ParsJson" orm:"column(pars_json);null" description:"动态参数JSON"`
	CreateTime    time.Time `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}
type QtPlugins struct {
	Id          int       `json:"Id" orm:"column(id);auto"`
	Module      string    `json:"Module" orm:"column(module);size(50)" description:"插件模块"`
	VersionCode int       `json:"VersionCode" orm:"column(version_code)" description:"插件版本"`
	Symbol      string    `json:"Symbol" orm:"column(symbol);size(50)" description:"入口函数"`
	Args        string    `json:"Args" orm:"column(args);size(255);null" description:"参数说明"`
	PluginName  string    `json:"PluginName" orm:"column(plugin_name);size(255)" description:"插件文件名"`
	CreateTime  time.Time `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}
type AiTrain struct {
	AiFrameworkId       int       `json:"aiFrameworkId"`
	ProjectId           int       `json:"projectId"` // 训练中心使用 项目名
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

