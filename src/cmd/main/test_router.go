package main

import (
	"github.com/dchest/captcha"
	"github.com/labstack/echo"
	"middleware/session"
	"modules/log"
	"net/http"
	"routers/test"
	"time"
)

// 路由
func tester(e *echo.Echo) {

	r := e.Group("/test")

	r.Get("", test.Index)
	r.Get("/add", test.Add)
	r.Get("/info", test.Info)

	r.Get("/render", func(c echo.Context) error {
		return c.Render(200, "hello.html", map[string]interface{}{
			"name": "weisd",
		})
	})

	r.Get("/error", func(c echo.Context) error {
		panic("panic to rendor error.html")
	})

	type Persion struct {
		Age     int
		Name    string
		Created time.Time
	}

	r.Get("/captcha", func(c echo.Context) error {
		data := map[string]interface{}{}
		data["captchaId"] = captcha.New()

		sliceA := []string{"a", "b", "c", "d"}

		mapA := map[string]string{"id": "123", "name": "weisd"}

		// f := session.FlashVal(c)

		// data["FLASH"] = f
		data["sliceA"] = sliceA
		data["mapA"] = mapA
		data["persion"] = Persion{19, "weisd", time.Now()}

		return c.Render(http.StatusOK, "test/captcha.html", data)
	})

	r.Post("/captcha", func(c echo.Context) error {

		log.Debug("%s %s", c.FormValue("captchaId"), c.FormValue("captchaSolution"))
		if !captcha.VerifyString(c.FormValue("captchaId"), c.FormValue("captchaSolution")) {
			return c.Redirect(302, "/test/captcha")
			// return c.String(http.StatusUnauthorized, "captcha err")
			// f := session.FlashObj(c)
			// f.Error("captcha check failed")

			// log.Debug(" set f %v", f)
		}

		// return c.String(http.StatusOK, "ok")
		return c.Redirect(302, "/test/captcha")

	})

	r.Get("/sess/set", func(c echo.Context) error {
		sess := session.GetStore(c)

		err := sess.Set("name", "weisd")
		if err != nil {
			log.Error("sess.set %v", err)
			return err
		}

		return c.String(200, "ok")
	})

	r.Get("/sess/get", func(c echo.Context) error {
		sess := session.GetStore(c)

		name := "nil"
		nameIf := sess.Get("name")
		switch nameIf.(type) {
		case string:
			name = nameIf.(string)
		case nil:
			name = "nil"
		}

		log.Error("get end")
		return c.String(200, name)
	})

}
