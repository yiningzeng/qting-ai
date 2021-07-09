package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"reflect"
	"strings"
	"time"
)

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

func (t *QtProjects) TableName() string {
	return "qt_projects"
}

func init() {
	orm.RegisterModel(new(QtProjects))
}

// AddQtProjects insert a new QtProjects into database and returns
// last inserted Id on success.
func AddQtProjects(m *QtProjects) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtProjectsById retrieves QtProjects by Id. Returns error if
// Id doesn't exist
func GetQtProjectsById(id int) (v *QtProjects, err error) {
	o := orm.NewOrm()
	v = &QtProjects{Id: id}
	if err = o.Read(v); err == nil {
		_, _ = o.LoadRelated(v, "Labels")
		return v, nil
	}
	return nil, err
}

// GetQtProjectsById retrieves QtProjects by Id. Returns error if
// Id doesn't exist
func GetQtProjectsByProjectName(name string) (v *QtProjects, err error) {
	o := orm.NewOrm()
	v = &QtProjects{ProjectName: name}
	if err = o.Read(v, "ProjectName"); err == nil {
		_, _ = o.LoadRelated(v, "Labels")
		return v, nil
	}
	return nil, err
}

// GetAllQtProjects retrieves all QtProjects matches certain condition. Returns empty list if
// no records exist
func GetAllQtProjects(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtProjects))
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

	var l []QtProjects
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).RelatedSel().All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "Labels")
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				_, _ = o.LoadRelated(&v, "Labels")
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		//for _, v := range ml {
		//	_, _ = o.LoadRelated(v, "labels")
		//}
		return ml, nil
	}
	return nil, err
}

// UpdateQtProjects updates QtProjects by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtProjectsById(m *QtProjects) (err error) {
	o := orm.NewOrm()
	v := QtProjects{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQtProjects deletes QtProjects by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtProjects(id int) (err error) {
	o := orm.NewOrm()
	v := QtProjects{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtProjects{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
