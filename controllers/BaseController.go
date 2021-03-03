package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	jsontime "github.com/liamylian/jsontime/v2/v2"
)
var json = jsontime.ConfigWithCustomTimeFormat
type BaseController struct {
	beego.Controller
}

type ReturnMsg struct {
	Code int
	Msg  string
	Data interface{}
}

type ReturnMsgAoi struct {
	Res         int       `json:"res"`
	Message  string   `json:"message"`
	ProjectList interface{} `json:"project_list"`
}

func (bc *BaseController) SuccessJson(data interface{}) {
	res := ReturnMsg{200, "success", data}
	bc.Data["json"] = res
	_ = bc.ServeJSON() //对json进行序列化输出
	bc.StopRun()
}

func (bc *BaseController) ErrorJson(code int, msg string, data interface{}) {
	res := ReturnMsg{code, msg, data}
	bc.Data["json"] = res
	_ = bc.ServeJSON() //对json进行序列化输出
	bc.StopRun()
}

func  (bc *BaseController) Ret(data interface{}, err error) {
	if err != nil {
		switch err.Error() {
		case "<QuerySeter> no row found":
			bc.ErrorJson(204, "未查询到内容", nil)
		default:
			bc.ErrorJson(500, err.Error(), nil)
		}
	} else {
		bc.SuccessJson(data)
	}
}

func  (bc *BaseController) RetAOI(data interface{}, err error) {
	if err != nil {
		var res ReturnMsgAoi
		switch err.Error() {
		case "<QuerySeter> no row found":
			res = ReturnMsgAoi{0, "success", data}
		default:
			res = ReturnMsgAoi{-1, "获取失败", nil}
		}
		bc.Data["json"] = res
		_ = bc.ServeJSON() //对json进行序列化输出
		bc.StopRun()
	} else {
		res := ReturnMsgAoi{0, "success", data}
		bc.Data["json"] = res
		_ = bc.ServeJSON() //对json进行序列化输出
		bc.StopRun()
	}
}
