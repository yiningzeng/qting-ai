## qting-ai
go版本训练中心 整合了beego 框架，包含API、Web 及后端服务等应用，是一个 RESTful 的框架

### 架构设计
![架构图](./doc/architecture.png)
### 执行逻辑
![执行逻辑](./doc/flow.png)

### 接口文档
> 接口内网地址 `http://192.168.31.77:8080/swagger/`

![api](./doc/api.png)

### web demo
> 内网地址 `http://192.168.31.77`

![api](./doc/demo.png)

### 插件说明
插件默认格式，否则无法正确加载

```golang
var (
	pluginName    = "QTing-tiny-3l" // 插件名-对应框架名称
	pluginVersion = 10 // 版本id
)

// 自动加载插件会首先读取版本信息，包含
func Version(args ...interface{}) (ret interface{}, err error) {
	t := models.QtPlugins{
		Module:      pluginName, 
		VersionCode: pluginVersion,
		Symbol:      "Run", // 该字段就是对应的下面执行的函数
		Args:        "taskId string, modelPath string", // 该字段表示执行函数的参数说明
		PluginName:  args[0].(string),
		CreateTime:  time.Now(),
	}
	return t, nil
}
// 默认插件执行的
func Run(args ...interface{}) (ret interface{}, err error) {
    return something, nil
}
```
         


