package controllers

import (
	"qting-ai/models"
)

// QtLabelsController operations for QtLabels
type AIController struct {
	BaseController
}

// URLMapping ...
func (c *AIController) URLMapping() {
	c.Mapping("Post", c.Train)
	//c.Mapping("GetOne", c.GetOneRabMsg)
	c.Mapping("Put", c.StopTrain)
	c.Mapping("AoiGetModels", c.AoiGetModels)
	//c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create AiTrain
// @Param	body		body 	models.AiTrain	true		"body for AiTrain content"
// @Success 200 {int} models.AiTrain
// @Failure 403 body is empty
// @Failure 704 json error
// @router / [post]
func (c *AIController) Train() {
	var v models.TrainBaseConfig
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if  err == nil {
		c.Ret(v, models.StartTrain(&v, string(c.Ctx.Input.RequestBody)))
	} else {
		c.ErrorJson(701, "json数据有误", err.Error())
	}
}

// Put ...
// @Title Put
// @Description update the QtAiFramework
// @Param	taskId		path 	string	true		"The id you want to update"
// @Success 200 {object} models.QtTrainRecord
// @Failure 403 :id is not int
// @router /stop/:taskId [put]
func (c *AIController) StopTrain() {
	idStr := c.Ctx.Input.Param(":taskId")
	err := models.StopTrain(idStr, models.Stopped, models.StoppedInt)
	c.Ret("停止成功", err)
}

// AoiGetModels ...
// @Title 用于Aoi软件获取发布的项目模型接口
// @Description 用于Aoi软件获取发布的项目模型接口
// @Success 200 {object} models.QtLabels
// @Failure 403 :id is empty
// @router /models/ [get]
func (c *AIController) AoiGetModels() {
	c.RetAOI(models.V1GetModelsForAoi())
}
