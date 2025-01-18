package main

import (
	_ "embed"

	"github.com/gcerrato/swaggerui-echo"
	"github.com/labstack/echo/v4"
)

//go:embed spec/petstore.yml
var spec []byte

func main() {
	e := echo.New()

	e.Any("/pet/*", petHandler)
	e.GET("/swagger/*", swaggerui.EchoHandler("/swagger", spec))

	e.Logger.Fatal(e.Start(":8080"))
}
