package controllers

import (
	"github.com/sirupsen/logrus"
	"qting-ai/tools"
	"qting-ai/version"
)

// QtLabelsController operations for QtLabels
type ToolsController struct {
	BaseController
}

// URLMapping ...
func (c *ToolsController) URLMapping() {
	c.Mapping("TestPlugin", c.TestPlugin)
	c.Mapping("GetVersion", c.GetVersion)
	c.Mapping("GetUpdateInfo", c.GetUpdateInfo)
	c.Mapping("PutUpdateInfo", c.PutUpdateInfo)
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

// GetVersion ...
// @Title GetUpdateInfo
// @Description 获取程序的更新信息
// @Success 200 获取成功
// @Failure 703 获取队列失败
// @router /update [get]
func (c *ToolsController) GetUpdateInfo() {
	c.Ret(tools.CheckUpdate())
}

// GetVersion ...
// @Title GetUpdateInfo
// @Description 获取程序的更新信息
// @Success 200 获取成功
// @Failure 703 获取队列失败
// @router /update [put]
func (c *ToolsController) PutUpdateInfo() {
	c.Ret(tools.Update())
}

// GetVersion ...
// @Title GetVersion
// @Description 获取程序的版本信息
// @Param	module	query	string	false	"module"
// @Param	function	query	string	false	"function"
// @Param	args	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Success 200 获取成功
// @Failure 703 获取队列失败
// @router /testPlugin [get]
func (c *ToolsController) TestPlugin() {
	//module := c.GetString("module")
	//function := c.GetString("function")
	//args := c.GetString("args")
	//c.Ret(models.HotPluginRun(module, function, args))
}

