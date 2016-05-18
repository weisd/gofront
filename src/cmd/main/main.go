package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"models"
	"modules/log"
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
}

func main() {
	flag.Parse()

	bootstrap()

	go func() {
		log.Info("pprof listen on %s", setting.Conf.Web.Pprof)
		log.Error("%v", http.ListenAndServe(setting.Conf.Web.Pprof, nil))
	}()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Gzip())
	e.Use(mw.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

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
