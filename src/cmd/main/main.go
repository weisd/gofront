package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	mw "github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	// Start server
	e.Run(standard.New(":1323"))

}
