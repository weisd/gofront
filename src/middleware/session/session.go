package session

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/session"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/labstack/echo"
	"net/url"
	// "time"
)

var GlobalSessions *session.Manager

var defaultOtions = Options{"memory", `{}`}

const (
	CONTEXT_SESSION_KEY = "_SESSION_STORE"
	COOKIE_FLASH_KEY    = "_COOKIE_FLASH"
	CONTEXT_FLASH_KEY   = "_FLASH_STORE"
	SESSION_FLASH_KEY   = "_SESSION_FLASH"
)

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

			fmt.Println("session start")
			sess, err := GlobalSessions.SessionStart(c.Response(), c.Request())
			if err != nil {
				return err
			}

			flashVals := url.Values{}

			flashIf := sess.Get(SESSION_FLASH_KEY)
			if flashIf != nil {
				vals, _ := url.QueryUnescape(flashIf.(string))
				flashVals, _ = url.ParseQuery(vals)
				if len(flashVals) > 0 {
					flash := Flash{}
					flash.ErrorMsg = flashVals.Get("error")
					flash.WarningMsg = flashVals.Get("warning")
					flash.InfoMsg = flashVals.Get("info")
					flash.SuccessMsg = flashVals.Get("success")
					c.SetData("FLASH", flash)
				}
			}

			f := NewFlash()

			sess.Set(SESSION_FLASH_KEY, f)

			c.Set(CONTEXT_SESSION_KEY, sess)

			defer func() {
				sess.Set(SESSION_FLASH_KEY, url.QueryEscape(f.Encode()))
				sess.SessionRelease(c.Response())
			}()

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

func GetFlash(c echo.Context) *Flash {
	return GetStore(c).Get(SESSION_FLASH_KEY).(*Flash)
}

func NewFlash() *Flash {
	return &Flash{url.Values{}, "", "", "", ""}
}

type Flash struct {
	url.Values
	ErrorMsg, WarningMsg, InfoMsg, SuccessMsg string
}

func (f *Flash) set(name, msg string) {
	f.Set(name, msg)
}

func (f *Flash) Error(msg string) {
	f.ErrorMsg = msg
	f.set("error", msg)
}

func (f *Flash) Warning(msg string) {
	f.WarningMsg = msg
	f.set("warning", msg)
}

func (f *Flash) Info(msg string) {
	f.InfoMsg = msg
	f.set("info", msg)
}

func (f *Flash) Success(msg string) {
	f.SuccessMsg = msg
	f.set("success", msg)
}
