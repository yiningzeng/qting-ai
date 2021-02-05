package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type QtLabels struct {
	Id         int         `orm:"column(id);auto"`
	ProjectId  *QtProjects `orm:"column(project_id);rel(fk)"`
	LabelName  string      `orm:"column(label_name);size(30);null" description:"标签名称"`
	Remarks    string      `orm:"column(remarks);size(255);null"`
	CreateTime time.Time   `orm:"column(create_time);type(datetime);null"`
}

func (t *QtLabels) TableName() string {
	return "qt_labels"
}

func init() {
	orm.RegisterModel(new(QtLabels))
}

// AddQtLabels insert a new QtLabels into database and returns
// last inserted Id on success.
func AddQtLabels(m *QtLabels) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtLabelsById retrieves QtLabels by Id. Returns error if
// Id doesn't exist
func GetQtLabelsById(id int) (v *QtLabels, err error) {
	o := orm.NewOrm()
	v = &QtLabels{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetQtLabelsByProjectId(projectId int) (v *QtLabels, err error) {
	o := orm.NewOrm()
	v = &QtLabels{ProjectId: &QtProjects{Id: projectId}}
	if err = o.Read(v, "ProjectId"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllQtLabels retrieves all QtLabels matches certain condition. Returns empty list if
// no records exist
func GetAllQtLabels(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtLabels))
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

	var l []QtLabels
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "ProjectId")
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "ProjectId")
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

// UpdateQtLabels updates QtLabels by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtLabelsById(m *QtLabels) (err error) {
	o := orm.NewOrm()
	v := QtLabels{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQtLabels deletes QtLabels by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtLabels(id int) (err error) {
	o := orm.NewOrm()
	v := QtLabels{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtLabels{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
