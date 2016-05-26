package main

import (
	"flag"
	"net/http"
	"os"

	"path/filepath"

	"middleware/session"
	"models"
	"modules/captcha"
	"modules/log"
	"modules/pongor"
	"setting"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
)

var configPath string

func init() {
	pwd, _ := os.Getwd()
	flag.StringVar(&configPath, "c", filepath.Join(pwd, "./src/cmd/main/app.toml"), "-c /path/to/app.toml config gile")
}

func bootstrap() {
	if err := setting.InitConf(configPath); err != nil {
		panic(err)
	}

	log.InitLogService(setting.Conf.Logs)

	// 初始化 mysql
	if err := models.InitModels(setting.Conf.Models); err != nil {
		panic(err)
	}
	// 初始化 redis
	models.InitRedis(setting.Conf.Redis)

	if err := session.InitSession(setting.Conf.Session); err != nil {
		panic(err)
	}

	captcha.InitCaptcha()
}

func main() {
	flag.Parse()

	bootstrap()

	start()

	go func() {
		log.Info("pprof listen on %s", setting.Conf.Web.Pprof)
		log.Error("%v", http.ListenAndServe(setting.Conf.Web.Pprof, nil))
	}()

	// Echo instance
	e := echo.New()

	// render
	render := pongor.GetRenderer(pongor.PongorOption{Directory: setting.Conf.Pongo.Directory, Reload: setting.Conf.Pongo.Reload})

	e.SetRenderer(render)

	// 固定返回值
	e.SetHTTPErrorHandler(func(err error, c echo.Context) {

		code := http.StatusInternalServerError
		msg := "服务器错误"

		switch err.(type) {
		case *echo.HTTPError:
			he := err.(*echo.HTTPError)
			code = he.Code
			msg = he.Message
		case *models.Err:
			log.Error("models ERR %s", err.Error())

			msg = "数据操作失败"

		default:
			// panic(err)
			log.Error("unknown ERR %T %v", err, err)
		}

		if e.Debug() {
			msg = err.Error()
			log.Error("%T, %v", err, err)
		}

		c.Render(http.StatusOK, "error.html", map[string]interface{}{"code": code, "msg": msg})
		// c.JSON(http.StatusOK, api.RetErr(code, msg))
	})

	if setting.Conf.Web.Debug {
		e.SetDebug(true)
	}

	// Middleware
	if setting.Conf.AccessLog.Enable {
		file := setting.Conf.AccessLog.FilePath
		if len(file) > 0 {

			os.Mkdir(filepath.Dir(file), os.ModePerm)

			f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
			if err != nil {
				log.Error("AccessLog OpenFile failed : %v ", err)
				err = nil
			}
			e.Use(mw.LoggerWithConfig(mw.LoggerConfig{Output: f}))
		} else {
			e.Use(mw.Logger())
		}
	}

	////////////////// middleware ///////////////

	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	e.Use(session.Sessioner())

	e.Static("/public", setting.Conf.Web.StaticDir)

	// e.File("/favicon.ico", "public/favicon.ico")
	e.Get("/favicon.ico", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// 验证码
	e.Get("/captcha/*.png", captcha.Server())

	// 路由
	router(e)
	// 测试路由
	tester(e)

	var server engine.Server

	if setting.Conf.Web.Fasthttp {
		server = fasthttp.New(setting.Conf.Web.Listen)
	} else {
		server = standard.New(setting.Conf.Web.Listen)
	}

	log.Info("server use %T, listen on %s", server, setting.Conf.Web.Listen)
	// Start server
	e.Run(server)

}
