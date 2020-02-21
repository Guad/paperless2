package main

import (
	"github.com/guad/paperless2/backend/broker"
	"github.com/guad/paperless2/backend/db"
	"github.com/guad/paperless2/backend/storage"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	storage.InitStorage()
	db.InitDB()
	broker.InitBroker()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	e.Use(middleware.Logger())

	route(e)

	e.Start(":8080")
}
