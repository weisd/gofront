package test

import (
	"github.com/labstack/echo"
	"net/http"

	m "models/test"
)

func Index(c echo.Context) error {

	list, err := m.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, list)
}

func Info(c echo.Context) error {

	info, err := m.InfoByName("da")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info)
}

func Add(c echo.Context) error {

	info, err := m.Add("da", "sdfsdf")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, info)
}
