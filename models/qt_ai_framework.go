package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type QtAiFramework struct {
	Id            int       `orm:"column(id);auto"`
	FrameworkName string    `orm:"column(framework_name);size(30)" description:"ai框架名称"`
	BaseImageUrl  string    `orm:"column(base_image_url);size(255)" description:"基础镜像地址"`
	ImageVersion  string    `orm:"column(image_version);size(255)" description:"镜像版本"`
	Volume        string    `orm:"column(volume);size(255)" description:"数据映射的镜像内部地址"`
	Cfg           string    `orm:"column(cfg);size(255)" description:"对应的配置文件"`
	Remarks       string    `orm:"column(remarks);size(100);null" description:"备注"`
	ParsJson      string    `orm:"column(pars_json);null" description:"动态参数JSON"`
	CreateTime    time.Time `orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}

func (t *QtAiFramework) TableName() string {
	return "qt_ai_framework"
}

func init() {
	orm.RegisterModel(new(QtAiFramework))
}

// AddQtAiFramework insert a new QtAiFramework into database and returns
// last inserted Id on success.
func AddQtAiFramework(m *QtAiFramework) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtAiFrameworkById retrieves QtAiFramework by Id. Returns error if
// Id doesn't exist
func GetQtAiFrameworkById(id int) (v *QtAiFramework, err error) {
	o := orm.NewOrm()
	v = &QtAiFramework{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetQtAiFrameworkById retrieves QtAiFramework by Id. Returns error if
// Id doesn't exist
func GetQtAiFrameworkByFrameworkName(name string) (v *QtAiFramework, err error) {
	o := orm.NewOrm()
	v = &QtAiFramework{FrameworkName: name}
	if err = o.Read(v, "FrameworkName"); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllQtAiFramework retrieves all QtAiFramework matches certain condition. Returns empty list if
// no records exist
func GetAllQtAiFramework(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtAiFramework))
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

	var l []QtAiFramework
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

// UpdateQtAiFramework updates QtAiFramework by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtAiFrameworkById(m *QtAiFramework) (err error) {
	o := orm.NewOrm()
	v := QtAiFramework{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteQtAiFramework deletes QtAiFramework by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtAiFramework(id int) (err error) {
	o := orm.NewOrm()
	v := QtAiFramework{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtAiFramework{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
