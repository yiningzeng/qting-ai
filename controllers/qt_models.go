package controllers

import (
	"errors"
	"qting-ai/models"
	"strconv"
	"strings"
)

// QtModelsController operations for QtModels
type QtModelsController struct {
	BaseController
}

// URLMapping ...
func (c *QtModelsController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("OnlineModel", c.OnlineModel)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAllQtModelsByLabelsAndMulti", c.GetAllQtModelsByLabelsAndMulti)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post 发布模型
// @Description create QtModels
// @Param	body		body 	models.OnlineModelPar	true		"body for QtModels content"
// @Success 201 {int} models.QtModels
// @Failure 403 body is empty
// @router /OnlineModel/ [post]
func (c *QtModelsController) OnlineModel() {
	var v models.OnlineModelPar
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		c.Ret(nil, models.OnlineModel(&v))
	} else {
		c.ErrorJson(509, "发布失败" + err.Error(), nil)
	}
}

// Post ...
// @Title Post
// @Description create QtModels
// @Param	body		body 	models.QtModels	true		"body for QtModels content"
// @Success 201 {int} models.QtModels
// @Failure 403 body is empty
// @router / [post]
func (c *QtModelsController) Post() {
	var v models.QtModels
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddQtModels(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get QtModels by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.QtModels
// @Failure 403 :id is empty
// @router /:id [get]
func (c *QtModelsController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetQtModelsById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAllQtModelsByLabelsAndMulti ...
// @Title 获取模型的所有模型通过字段
// @Description get QtModels
// @Param	projectId	query	string	false	"项目id"
// @Param	aiFrameworkId	query	string	false	"AI框架id"
// @Param	label	query	string	false	"单个标签名"
// @Param	isMultilabel	query	string	false	"是否是多类单标签，如果该值!=0 那么label上面的值会失效, 只有在QTing-tiny-3l-single框架下才=0"
// @Success 200 {object} models.QtModels
// @Failure 403
// @router /GetAllQtModelsByLabelsAndMulti/ [get]
func  (c *QtModelsController) GetAllQtModelsByLabelsAndMulti() {
	projectId, _ := strconv.Atoi(c.GetString("projectId"))
	aiFrameworkId, _ := strconv.Atoi(c.GetString("aiFrameworkId"))
	c.Ret(models.GetAllQtModelsByLabelsAndMulti(projectId, aiFrameworkId, c.GetString("label"), c.GetString("isMultilabel")))
}

// GetAll ...
// @Title Get All
// @Description get QtModels
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.QtModels
// @Failure 403
// @router / [get]
func (c *QtModelsController) GetAll() {
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

	l, err := models.GetAllQtModels(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the QtModels
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.QtModels	true		"body for QtModels content"
// @Success 200 {object} models.QtModels
// @Failure 403 :id is not int
// @router /:id [put]
func (c *QtModelsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.QtModels{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateQtModelsById(&v); err == nil {
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
// @Description delete the QtModels
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *QtModelsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteQtModels(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

