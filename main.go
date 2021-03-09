package main

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"qting-ai/models"
	"qting-ai/plugins"
	_ "qting-ai/routers"
	"qting-ai/version"
	"strings"
	"time"
)


func LoggerInit(debug bool) {
	path := "log/yiningzeng.log"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	if debug {
		logrus.SetFormatter(&logrus.TextFormatter{})
		//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
		//同时写文件和屏幕
		//fileAndStdoutWriter := io.MultiWriter([]io.Writer{writer, os.Stdout}...)
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		writer, _ := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			// rotatelogs.WithMaxAge(time.Duration(180)*time.),
			rotatelogs.WithRotationCount(60),
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		)
		logrus.SetOutput(writer)
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func Start(state overseer.State) {
	LoggerInit(beego.AppConfig.DefaultBool("debug", false))
	sqlConn, er := beego.AppConfig.String("sqlconn")
	if er != nil {
		logrus.Panic("读取数据库连接出错出错")
	}
	if err := orm.RegisterDataBase("default", "mysql", sqlConn); err != nil {
		logrus.WithField("err", err).Panic("初始化数据库失败")
	}

	fs := afero.NewOsFs()
	_ = fs.Remove("update")
	version.PrintVersion()

	//go models.WatchDir(qTingFolder)
	go models.StartCron()
	go plugins.Watcher() // 监听插件目录


	//if beego.BConfig.RunMode == "dev" {
	//	beego.BConfig.WebConfig.DirectoryIndex = true
	//	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	//}
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.BConfig.WebConfig.StaticDir["/manage"] = "manage"
	beego.BConfig.WebConfig.StaticDir[beego.AppConfig.DefaultString("ProjectPathStaticDir", "/qting")] = beego.AppConfig.DefaultString("ProjectPath", "/qtingvisionfolder/Projects/")
	//InsertFilter是提供一个过滤函数
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//允许访问所有源
		AllowAllOrigins: true,
		//可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" ()
		//其中Options跨域复杂请求预检
		AllowMethods:   []string{"GET","POST","OPTIONS","PUT","DELETE"}, // 跨域问题，PUT和DELETE会报错， 所以这里不能用*代替，一定要写全
		//指的是允许的Header的种类
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "content-type"},
		//公开的HTTP标头列表
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		//如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
		AllowOrigins: []string{"http://10.*.*.*:*","http://localhost:*","http://127.0.0.1:*"},
	}))
	if state.Enabled {
		_ = fs.Remove("update-info")
		_ = afero.WriteFile(fs, "update-success", []byte(time.Now().Format("2006-01-02 15:04:05")), 0755)
	}
	beego.Run()
}

func main() {
	timeZoneShanghai, _ := time.LoadLocation("Asia/Shanghai")
	jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
	jsontime.AddLocaleAlias("shanghai", timeZoneShanghai)
	//_, err := fsnotify.NewWatcher()
	//if err != nil {
	//	logrus.Error(fmt.Sprintf("%+v", errors.Wrap(err, "")))
	//}

	mode := beego.AppConfig.DefaultString("runmode", "prod")
	if strings.EqualFold("prod", mode) {
		overseer.Run(overseer.Config{
			Program: Start,
			PreUpgrade: func(tempBinaryPath string) error {
				logrus.WithFields(logrus.Fields{
					"老版本ID": version.ID,
					"编译日期":  version.BuildDate,
				}).Info("已做好老版本备份，删除更新成功状态，开始执行更新")
				fs := afero.NewOsFs()
				_ = afero.WriteFile(fs, "update-info", []byte("已做好老版本备份，删除更新成功状态，开始执行更新"), 0755)
				_ = fs.Remove("update-success")
				if b, err := afero.ReadFile(fs, "qting-ai"); err == nil {
					_ = afero.WriteFile(fs, "qting-ai(old)", b, 0755)
				}
				return nil
			},
			Fetcher: &fetcher.File{Path: "update"},
			Debug:   false, //display log of overseer actions
		})
	} else {
		Start(overseer.State{Enabled: false})
	}
}

