package main

import (
	"balance/cmd/controllers"
	"balance/pkg/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := db.MigrateDb()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/balance/pay", controllers.SendMoney)
	e.POST("/balance/users", controllers.PayBetweenUsers)
	e.POST("/balance/get", controllers.SelectBalance)

	err = e.Start(":9999")
	if err != nil {
		return
	}

}
