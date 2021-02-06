// @APIVersion 3.1.3
// @Title 训练中心
// @Description 自训练中心后台接口
// @Contact yiningzeng@outlook.com
// @TermsOfServiceUrl http://www.qtingvision.com/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"qting-ai/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/qt_ai_framework",
			beego.NSInclude(
				&controllers.QtAiFrameworkController{},
			),
		),

		beego.NSNamespace("/qt_labels",
			beego.NSInclude(
				&controllers.QtLabelsController{},
			),
		),

		beego.NSNamespace("/qt_models",
			beego.NSInclude(
				&controllers.QtModelsController{},
			),
		),

		beego.NSNamespace("/qt_projects",
			beego.NSInclude(
				&controllers.QtProjectsController{},
			),
		),

		beego.NSNamespace("/qt_rabbitmq_info",
			beego.NSInclude(
				&controllers.QtRabbitmqInfoController{},
			),
		),

		beego.NSNamespace("/qt_train_record",
			beego.NSInclude(
				&controllers.QtTrainRecordController{},
			),
		),
		beego.NSNamespace("/tools",
			beego.NSInclude(
				&controllers.ToolsController{},
			),
		),
		beego.NSNamespace("/ai",
			beego.NSInclude(
				&controllers.AIController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
