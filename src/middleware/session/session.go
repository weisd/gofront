package session

import (
	"errors"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

var GlobalSessions *session.Manager

var defaultOtions = Options{"memory", `{}`}

const CONTEXT_SESSION_KEY = "_SESSION_STORE"

type Options struct {
	Provider string
	Config   string
}

/*
type managerConfig struct {
	CookieName      string `json:"cookieName"`
	EnableSetCookie bool   `json:"enableSetCookie,omitempty"`
	Gclifetime      int64  `json:"gclifetime"`
	Maxlifetime     int64  `json:"maxLifetime"`
	Secure          bool   `json:"secure"`
	CookieLifeTime  int    `json:"cookieLifeTime"`
	ProviderConfig  string `json:"providerConfig"`
	Domain          string `json:"domain"`
	SessionIDLength int64  `json:"sessionIDLength"`
}
*/

func InitSession(op ...Options) error {
	option := defaultOtions
	if len(op) > 0 {
		option = op[0]
	}

	if len(option.Provider) == 0 {
		option.Provider = defaultOtions.Provider
	}

	if len(option.Config) == 0 {
		option.Config = defaultOtions.Config
	}

	var err error
	GlobalSessions, err = session.NewManager(option.Provider, option.Config)
	if err != nil {
		return err
	}
	go GlobalSessions.GC()

	return nil
}

func Sessioner() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if GlobalSessions == nil {

				return errors.New("session manager not found, use session middleware but not init ?")
			}

			sess, err := GlobalSessions.SessionStart(c.Response().(*standard.Response).ResponseWriter, c.Request().(*standard.Request).Request)
			if err != nil {
				return err
			}
			defer sess.SessionRelease(c.Response().(http.ResponseWriter))

			c.Set(CONTEXT_SESSION_KEY, sess)

			return next(c)
		}
	}
}

func GetStore(c echo.Context) session.Store {
	store := c.Get(CONTEXT_SESSION_KEY)
	if store != nil {
		return store.(session.Store)
	}

	return nil
}
