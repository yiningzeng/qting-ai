package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type QtTrainRecord struct {
	Id          int         `orm:"column(id);auto"`
	TaskId      string      `orm:"column(task_id);size(30);null" description:"任务标识"`
	TaskName    string      `orm:"column(task_name);size(100);null" description:"训练任务名称"`
	ContainerId string      `orm:"column(container_id);size(100);null" description:"训练的容器的id"`
	ProjectId   *QtProjects `orm:"column(project_id);rel(fk)" description:"项目名"`
	Status      int         `orm:"column(status);null" description:"0:准备完成
1:等待训练
2:正在训练
3:暂停训练
4:训练完成
-1:训练有误"`
	AiFrameworkId *QtAiFramework `orm:"column(ai_framework_id);rel(fk)" description:"训练使用的框架ID"`
	AssetsType    string         `orm:"column(assets_type);size(50);null" description:"训练的数据类型"`
	IsJump        int            `orm:"column(is_jump);null" description:"是否插队
0:正常
1:插队"`
	DrawUrl    string    `orm:"column(draw_url);size(255);null" description:"画图的网址"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null"`
}

func (t *QtTrainRecord) TableName() string {
	return "qt_train_record"
}

func init() {
	orm.RegisterModel(new(QtTrainRecord))
}

// AddQtTrainRecord insert a new QtTrainRecord into database and returns
// last inserted Id on success.
func AddQtTrainRecord(m *QtTrainRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtTrainRecordById retrieves QtTrainRecord by Id. Returns error if
// Id doesn't exist
func GetQtTrainRecordById(id int) (v *QtTrainRecord, err error) {
	o := orm.NewOrm()
	v = &QtTrainRecord{Id: id}
	if err = o.Read(v); err == nil {
		_, _ = o.LoadRelated(v, "ProjectId")
		_, _ = o.LoadRelated(v, "AiFrameworkId")
		return v, nil
	}
	return nil, err
}

// GetQtTrainRecordByTaskId retrieves QtTrainRecord by Id. Returns error if
// taskId doesn't exist
func GetQtTrainRecordByTaskId(taskId string) (v *QtTrainRecord, err error) {
	o := orm.NewOrm()
	if strings.Contains(taskId, "-") {
		taskId = strings.Split(taskId, "-")[0]
	}
	v = &QtTrainRecord{TaskId: taskId}
	if err = o.Read(v, "TaskId"); err == nil {
		_, _ = o.LoadRelated(v, "ProjectId")
		_, _ = o.LoadRelated(v, "AiFrameworkId")
		return v, nil
	}
	return nil, err
}

// GetAllQtTrainRecord retrieves all QtTrainRecord matches certain condition. Returns empty list if
// no records exist
func GetAllQtTrainRecord(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtTrainRecord))
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

	var l []QtTrainRecord
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "ProjectId")
				_, _ = o.LoadRelated(&v, "AiFrameworkId")
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "ProjectId")
				_, _ = o.LoadRelated(&v, "AiFrameworkId")
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

// UpdateQtTrainRecord updates QtTrainRecord by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtTrainRecordById(m *QtTrainRecord) (err error) {
	o := orm.NewOrm()
	v := QtTrainRecord{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQtTrainRecord deletes QtTrainRecord by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtTrainRecord(id int) (err error) {
	o := orm.NewOrm()
	v := QtTrainRecord{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtTrainRecord{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
