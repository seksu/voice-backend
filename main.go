package main

import (
	"fmt"
	"go-thai-dialect/controllers"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Version float64

func main() {

	var v Version = 0.9

	fmt.Println(time.Now().UnixNano() / int64(time.Millisecond))

	fmt.Println("API SERVER V.", v)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Static("/api/image", "image")
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	_publicAPI := e.Group("/api")

	controllers.VolunteerDock(_publicAPI)
	controllers.DialectDock(_publicAPI)
	controllers.RecordDock(_publicAPI)
	controllers.UserDock(_publicAPI)

	e.Logger.Fatal(e.Start(":9001"))

}
