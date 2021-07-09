package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type QtPlugins struct {
	Id          int       `json:"Id" orm:"column(id);auto"`
	Module      string    `json:"Module" orm:"column(module);size(50)" description:"插件模块"`
	VersionCode int       `json:"VersionCode" orm:"column(version_code)" description:"插件版本"`
	Symbol      string    `json:"Symbol" orm:"column(symbol);size(50)" description:"入口函数"`
	Args        string    `json:"Args" orm:"column(args);size(255);null" description:"参数说明"`
	PluginName  string    `json:"PluginName" orm:"column(plugin_name);size(255)" description:"插件文件名"`
	CreateTime  time.Time `json:"CreateTime" orm:"column(create_time);type(datetime);null" time_format:"sql_datetime" time_location:"shanghai" time_utc:"false"`
}

func (t *QtPlugins) TableName() string {
	return "qt_plugins"
}

func init() {
	orm.RegisterModel(new(QtPlugins))
}

// AddQtPlugins insert a new QtPlugins into database and returns
// last inserted Id on success.
func AddQtPlugins(m *QtPlugins) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetQtPluginsById retrieves QtPlugins by Id. Returns error if
// Id doesn't exist
func GetQtPluginsById(id int) (v *QtPlugins, err error) {
	o := orm.NewOrm()
	v = &QtPlugins{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllQtPlugins retrieves all QtPlugins matches certain condition. Returns empty list if
// no records exist
func GetAllQtPlugins(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(QtPlugins))
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

	var l []QtPlugins
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

// UpdateQtPlugins updates QtPlugins by Id and returns error if
// the record to be updated doesn't exist
func UpdateQtPluginsById(m *QtPlugins) (err error) {
	o := orm.NewOrm()
	v := QtPlugins{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

func UpdateQtPluginsByMV(module string, versionCode int, newFileName string) (err error) {
	o := orm.NewOrm()
	v := QtPlugins{Module: module, VersionCode: versionCode}
	// ascertain id exists in the database
	if err = o.Read(&v, "Module", "VersionCode"); err == nil {
		var num int64
		v.PluginName = newFileName
		if num, err = o.Update(&v); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
// DeleteQtPlugins deletes QtPlugins by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQtPlugins(id int) (err error) {
	o := orm.NewOrm()
	v := QtPlugins{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&QtPlugins{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func DeleteQtPluginsByFileName(fileName string) (err error) {
	o := orm.NewOrm()
	v := QtPlugins{PluginName: fileName}
	// ascertain id exists in the database
	if err = o.Read(&v, "PluginName"); err == nil {
		var num int64
		if num, err = o.Delete(&v); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
