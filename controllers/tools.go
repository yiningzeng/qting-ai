package controllers

import (
	"github.com/sirupsen/logrus"
	"qting-ai/models"
	"qting-ai/version"
)

// QtLabelsController operations for QtLabels
type ToolsController struct {
	BaseController
}

// URLMapping ...
func (c *ToolsController) URLMapping() {
	c.Mapping("Put", c.FsnotifyAddWatch)
	c.Mapping("GetVersion", c.GetVersion)
}

// GetVersion ...
// @Title GetVersion
// @Description 获取程序的版本信息
// @Success 200 获取成功
// @Failure 703 获取队列失败
// @router /version [get]
func (c *ToolsController) GetVersion() {
	c.SuccessJson(logrus.Fields{
		"VersionId": version.ID,
		"BuildDate": version.BuildDate,
		})
}

// put ...
// @Title Put
// @Description create QtLabels
// @Param	dir		path 	string	true		"需要监控的文件夹目录"
// @Success 201 {int} 监控成功
// @Failure 403 body is empty
// @router /fsnotify/:dir [put]
func (c *ToolsController) FsnotifyAddWatch() {
	dir := c.Ctx.Input.Param(":dir")
	c.Ret(nil, models.AddWatchDir(dir))
}
