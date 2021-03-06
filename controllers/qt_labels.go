package controllers

import (
	"errors"
	"qting-ai/models"
	"strconv"
	"strings"
)

// QtLabelsController operations for QtLabels
type QtLabelsController struct {
	BaseController
}

// URLMapping ...
func (c *QtLabelsController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetOneByProjectId", c.GetOneByProjectId)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create QtLabels
// @Param	body		body 	models.QtLabels	true		"body for QtLabels content"
// @Success 201 {int} models.QtLabels
// @Failure 403 body is empty
// @router / [post]
func (c *QtLabelsController) Post() {
	var v models.QtLabels
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddQtLabels(&v); err == nil {
			c.SuccessJson(v)
		} else {
			c.ErrorJson(501,err.Error(),nil)
		}
	} else {
		c.ErrorJson(501,err.Error(),nil)
	}
}

// GetOne ...
// @Title Get One
// @Description get QtLabels by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.QtLabels
// @Failure 403 :id is empty
// @router /:id [get]
func (c *QtLabelsController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	c.Ret(models.GetQtLabelsById(id))
}

// GetOneByProjectId ...
// @Title GetOneByProjectId
// @Description get QtLabels by ProjectId
// @Param	projectId		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.QtLabels
// @Failure 403 :projectId is empty
// @router /GetOneByProjectId/:projectId [get]
func (c *QtLabelsController) GetOneByProjectId() {
	idStr := c.Ctx.Input.Param(":projectId")
	id, _ := strconv.Atoi(idStr)
	c.Ret(models.GetQtLabelsByProjectId(id))
}

// GetAll ...
// @Title Get All
// @Description get QtLabels
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.QtLabels
// @Failure 403
// @router / [get]
func (c *QtLabelsController) GetAll() {
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
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}
	c.Ret(models.GetAllQtLabels(query, fields, sortby, order, offset, limit))
}

// Put ...
// @Title Put
// @Description update the QtLabels
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.QtLabels	true		"body for QtLabels content"
// @Success 200 {object} models.QtLabels
// @Failure 403 :id is not int
// @router /:id [put]
func (c *QtLabelsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.QtLabels{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateQtLabelsById(&v); err == nil {
			c.SuccessJson(v)
		} else {
			c.ErrorJson(501,err.Error(),nil)
		}
	} else {
		c.ErrorJson(501,err.Error(),nil)
	}
}

// Delete ...
// @Title Delete
// @Description delete the QtLabels
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *QtLabelsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	c.Ret(nil, models.DeleteQtLabels(id))
}
