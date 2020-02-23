package main

import (
	"github.com/guad/paperless2/backend/api"
	"github.com/guad/paperless2/backend/api/rest"
	"github.com/guad/paperless2/backend/api/user"
	"github.com/labstack/echo"
)

func route(e *echo.Echo) {
	g := e.Group("/api")

	g.POST("/login", user.Login)

	g.POST("/push", api.PushFile, user.AuthMiddleware)
	g.GET("/fetch/:doc/:fname", api.FetchFile, user.AuthMiddleware)
	g.GET("/thumb/:doc", api.GetThumbnail, user.AuthMiddleware)

	doc := g.Group("/document")

	doc.Use(user.AuthMiddleware)

	doc.GET("", rest.ListDocuments)
	doc.GET("/:id", rest.GetDocument)
	doc.PUT("/:id", rest.UpdateDocument)
	doc.DELETE("/:id", rest.DeleteDocument)

	tag := g.Group("/tag")

	tag.Use(user.AuthMiddleware)

	tag.GET("", rest.ListTags)
	tag.POST("", rest.CreateTag)
	tag.GET("/:id", rest.GetTag)
	tag.PUT("/:id", rest.UpdateTag)
	tag.DELETE("/:id", rest.DeleteTag)
}
