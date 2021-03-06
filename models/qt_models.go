package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/spf13/afero"
	"os"
	"path"
	"qting-ai/tools"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type QtModels struct {
	Id            int            `orm:"column(id);auto"`
	TaskId        string         `orm:"column(task_id);size(100);null" description:"主要是方便后期更新这个模型通过taskId来查询"`
	ProjectId     *QtProjects    `orm:"column(project_id);rel(fk);null" description:"项目名"`
	AiFrameworkId *QtAiFramework `orm:"column(ai_framework_id);rel(fk);null" description:"训练使用的框架ID"`
	ModelName     string         `orm:"column(model_name);size(100);null" description:"模型名称"`
	ModelPath     string         `orm:"column(model_path);size(255);null" description:"模型的本地路径"`
	ModelBasePath string         `orm:"column(model_base_path);size(255);null" description:"模型的保存目录"`
	IsMultilabel int `orm:"column(is_multilabel);null" description:"该模型是否包含多标签 -1 表示未知 0 表示单标签单模型 1 表示多标签单模型"`
	LabelStr string `orm:"column(label_str);null" description:"缺陷的标签，使用英文的,分离 比如zhanxi,jinmian,wuran"`
	SuggestScore string `orm:"column(suggest_score);null" description:"推荐置信度的json值"`
	Status int `orm:"column(status);null" description:"模型状态 0 表示未发布 1 表示已发布 2 表示已下线"`
	PublishUrl  string    `orm:"column(publish_url);size(255);null" description:"模型的发布地址"`
	PublishTime time.Time `orm:"column(publish_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}

type OnlineModelPar struct {
	ModelId int `json:"ModelId"`
	ModelWidth int `json:"ModelWidth"`
	ModelHeight int `json:"ModelHeight"`
	Label string `json:"Label"` // 用于单个标签的模型的下线操作，否则把这个项目的都下线了
}

func (t *QtModels) TableName() string {
	return "qt_models"
}

func init() {
	orm.RegisterModel(new(QtModels))
}
// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
// 模型上线
func OnlineModel(m *OnlineModelPar) (err error) {
	if qTModel, err := GetQtModelsById(m.ModelId); err == nil {
		if qTModel.Status == 1 {
			qTModel.Status = 2
		} else {
			qTModel.Status = 2
			qTModel.PublishTime = time.Now()
		}
		// 1.0的接口这里需要查询下把之前的模型给下线
		if oldQtModel, err := GetOneQtModelsByProjectAndMulti(qTModel.ProjectId, qTModel.AiFrameworkId, qTModel.IsMultilabel, m.Label); err == nil {
			oldQtModel.Status = 1
			_ = UpdateQtModelsById(oldQtModel)
		}
		// region 打包发布
		zipFile := ""
		label := ""
		assetPath := beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/")

		if strings.HasSuffix(qTModel.ProjectId.AssetsPath, "/") {
			assetPath = assetPath + qTModel.ProjectId.ProjectName
		} else {
			assetPath = assetPath + "/" + qTModel.ProjectId.ProjectName
		}

		if exist, _ := PathExists(path.Join(assetPath, "model_release")); !exist {
			// 创建文件夹
			_ = os.Mkdir(path.Join(assetPath, "model_release"), os.ModePerm)
		}

		if qTModel.IsMultilabel == 0 {
			label = strings.ReplaceAll(qTModel.LabelStr, ",", "")
			zipFile = fmt.Sprintf("%s/model_release/%s.zip", assetPath, label)
		} else {
			label = "allLabels"
			zipFile = fmt.Sprintf("%s/model_release/allLabels.zip",assetPath)
		}
		_ = afero.WriteFile(afero.NewOsFs(),
			  path.Join(qTModel.ModelBasePath, qTModel.TaskId + ".modelInfo"),
			[]byte(fmt.Sprintf("项目名称: %s\n模型名称: %s\n发布日期: %s\n", qTModel.ProjectId.ProjectName, qTModel.ModelName, qTModel.PublishTime)),
			0755)
		if strings.Contains(qTModel.AiFrameworkId.FrameworkName, "QTing-tiny-3l") {
			// 这里是老的yolo的打包方式
			err := tools.ZipV1(zipFile, qTModel.ModelBasePath, qTModel.TaskId, label, strconv.Itoa(m.ModelWidth), strconv.Itoa(m.ModelHeight))
			if err != nil {
				return err
			}
		} else {
			// 这里是最新的打包方式,以后的框架都是这种打包文件夹的形式!!!!!!!!
			if _, err = tools.PluginRun(qTModel.AiFrameworkId.FrameworkName, "Compress", []string{qTModel.ModelPath}, zipFile); err != nil {
				return err
			}
		}

		// endregion
		qTModel.PublishUrl = strings.ReplaceAll(zipFile,
			beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/"),
			beego.AppConfig.DefaultString("DownloadHost", "http://localhost:8080") + beego.AppConfig.DefaultString("ProjectPathStaticDir", "/qting")+"/")
		err = UpdateQtModelsById(qTModel)
		return nil
	} else {
		return err
	}
}

// AddQtModels insert a new QtModels into database and returns
// last inserted Id on success.
func AddQtModels(m *QtModels) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtModelsById retrieves QtModels by Id. Returns error if
// Id doesn't exist
func GetQtModelsById(id int) (v *QtModels, err error) {
	o := orm.NewOrm()
	v = &QtModels{Id: id}
	if err = o.Read(v); err == nil {
		_, _ = o.LoadRelated(v, "ProjectId")
		_, _ = o.LoadRelated(v, "AiFrameworkId")
		return v, nil
	}
	return nil, err
}

// GetQtModelsById retrieves QtModels by Id. Returns error if
// Id doesn't exist
func GetQtModelsByTaskId(taskId string) (v *QtModels, err error) {
	o := orm.NewOrm()
	v = &QtModels{TaskId: taskId}
	if err = o.Read(v, "TaskId"); err == nil {
		_, _ = o.LoadRelated(v, "ProjectId")
		_, _ = o.LoadRelated(v, "AiFrameworkId")
		return v, nil
	}
	return nil, err
}

func GetOneQtModelsByProjectAndMulti(project *QtProjects, aiFramework *QtAiFramework, isMultilabel int, label string) (v *QtModels, err error) {
	o := orm.NewOrm()
	if isMultilabel == 0 {
		v = &QtModels{ProjectId: project, AiFrameworkId: aiFramework, IsMultilabel: isMultilabel, Status: 2, LabelStr: ","+label+","}
		if err = o.Read(v, "ProjectId", "IsMultilabel", "Status", "LabelStr"); err == nil {
			_, _ = o.LoadRelated(v, "ProjectId")
			_, _ = o.LoadRelated(v, "AiFrameworkId")
			return v, nil
		}
		return nil, err
	} else {
		v = &QtModels{ProjectId: project, IsMultilabel: isMultilabel, Status: 2}
		if err = o.Read(v, "ProjectId", "IsMultilabel", "Status"); err == nil {
			_, _ = o.LoadRelated(v, "ProjectId")
			_, _ = o.LoadRelated(v, "AiFrameworkId")
			return v, nil
		}
		return nil, err
	}
}

func GetAllQtModelsByLabelsAndMulti(projectId int, aiFrameworkId int, oneLabel string, isMultilabel string) (ml []interface{}, err error) {
	var l []QtModels
	o := orm.NewOrm()
	aiFramework, aiFrErr := GetQtAiFrameworkById(aiFrameworkId)
	if project, err := GetQtProjectsById(projectId); err == nil {
		qs := o.QueryTable(new(QtModels))
		qs = qs.OrderBy("-Status", "-PublishTime", "-CreateTime")
		if strings.EqualFold(isMultilabel, "0") {
			qs = qs.Filter("ProjectId", project).Filter("IsMultilabel", isMultilabel).Filter("LabelStr__icontains", fmt.Sprintf(",%s,", oneLabel))
		} else {
			qs = qs.Filter("ProjectId", project).Filter("IsMultilabel", isMultilabel)
		}
		if aiFrErr == nil {
			qs = qs.Filter("AiFrameworkId", aiFramework)
		}
		if _, err = qs.All(&l); err == nil {
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "ProjectId")
				_, _ = o.LoadRelated(&v, "AiFrameworkId")
				ml = append(ml, v)
			}
			return ml, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// GetAllQtModels retrieves all QtModels matches certain condition. Returns empty list if
// no records exist
func GetAllQtModels(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtModels))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []QtModels
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateQtModels updates QtModels by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtModelsById(m *QtModels) (err error) {
	o := orm.NewOrm()
	v := QtModels{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQtModels deletes QtModels by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtModels(id int) (err error) {
	o := orm.NewOrm()
	v := QtModels{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtModels{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
