
# swaggerui-echo
Embedded, self-hosted [Swagger Ui](https://swagger.io/tools/swagger-ui/) for go servers

This module provides `swaggerui.EchoHandler`, which you can use to serve an embedded copy of [Swagger UI](https://swagger.io/tools/swagger-ui/) as well as an embedded specification for your API.

## Example usage
```go
package main

import (
	_ "embed"
	"log"

	"github.com/labstack/echo/v4"
	swaggerui "github.com/gcerrato/swaggerui-echo"
)

//go:embed swagger.json
var spec []byte

func main() {
	e := echo.New()
	
	e.Any("/pet/*", petHandler)
	e.GET("/swagger/*", swaggerui.EchoHandler("/swagger", spec))
	
	log.Println("serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
