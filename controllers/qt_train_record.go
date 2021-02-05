package controllers

import (
	"encoding/json"
	"errors"
	"qting-ai/models"
	"strconv"
	"strings"
)

// QtTrainRecordController operations for QtTrainRecord
type QtTrainRecordController struct {
	BaseController
}

// URLMapping ...
func (c *QtTrainRecordController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOneByTaskId", c.GetOneByTaskId)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create QtTrainRecord
// @Param	body		body 	models.QtTrainRecord	true		"body for QtTrainRecord content"
// @Success 201 {int} models.QtTrainRecord
// @Failure 403 body is empty
// @router / [post]
func (c *QtTrainRecordController) Post() {
	var v models.QtTrainRecord
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddQtTrainRecord(&v); err == nil {
			c.SuccessJson(v)
		} else {
			c.ErrorJson(501, "新增失败", err.Error())
		}
	} else {
		c.ErrorJson(502, "新增失败", err.Error())
	}
}

// GetOne ...
// @Title Get One
// @Description get QtTrainRecord by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.QtTrainRecord
// @Failure 403 :id is empty
// @router /:id [get]
func (c *QtTrainRecordController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetQtTrainRecordById(id)
	c.Ret(v, err)
}

// GetOneByTaskId ...
// @Title Get One By TaskId
// @Description get QtTrainRecord by TaskId
// @Param	taskId		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.QtTrainRecord
// @Failure 403 :id is empty
// @router /GetOneByTaskId/:taskId [get]
func (c *QtTrainRecordController) GetOneByTaskId() {
	idStr := c.Ctx.Input.Param(":taskId")
	v, err := models.GetQtTrainRecordByTaskId(idStr)
	c.Ret(v, err)
}

// GetAll ...
// @Title Get All
// @Description get QtTrainRecord
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.QtTrainRecord
// @Failure 403
// @router / [get]
func (c *QtTrainRecordController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Ret(nil, errors.New("Error: invalid query key/value pair"))
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllQtTrainRecord(query, fields, sortby, order, offset, limit)
	c.Ret(l, err)
}

// Put ...
// @Title Put
// @Description update the QtTrainRecord
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.QtTrainRecord	true		"body for QtTrainRecord content"
// @Success 200 {object} models.QtTrainRecord
// @Failure 403 :id is not int
// @router /:id [put]
func (c *QtTrainRecordController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.QtTrainRecord{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateQtTrainRecordById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the QtTrainRecord
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *QtTrainRecordController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteQtTrainRecord(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
