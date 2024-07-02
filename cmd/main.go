package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "It's work")
	})

	e.POST("/login", Greetings)

	e.Logger.Fatal(e.Start(":1323"))
}

func Greetings(c echo.Context) error {
	return c.JSON(http.StatusOK, HelloWorld{
		Message: "Hello World",
	})
}
