package main

import (
	"github.com/guad/paperless2/backend/api"
	"github.com/labstack/echo"
)

func route(e *echo.Echo) {
	g := e.Group("/api")

	g.POST("/login", api.Login)

	g.POST("/push", api.PushFile, api.AuthMiddleware)
	g.GET("/fetch/:doc/:fname", api.FetchFile, api.AuthMiddleware)
	g.GET("/thumb/:doc", api.GetThumbnail, api.AuthMiddleware)

	doc := g.Group("/document")

	doc.Use(api.AuthMiddleware)

	doc.GET("", api.ListDocuments)
	doc.GET("/:id", api.GetDocument)
	doc.PUT("/:id", api.UpdateDocument)
	doc.DELETE("/:id", api.DeleteDocument)

	tag := g.Group("/tag")

	tag.Use(api.AuthMiddleware)

	tag.GET("", api.ListTags)
	tag.POST("", api.CreateTag)
	tag.GET("/:id", api.GetTag)
	tag.PUT("/:id", api.UpdateTag)
	tag.DELETE("/:id", api.DeleteTag)
}
