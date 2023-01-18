package main

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 验证 go 的文件二进制嵌入

//go:embed static/kubernetes.png
var content []byte

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":12345"))
}

// Handler
func hello(c echo.Context) error {
	return c.Blob(http.StatusOK, "image/png", content)
}
