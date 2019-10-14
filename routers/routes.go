package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"image-server/pkg/conf"
)

func Register(echo *echo.Echo) {

	bGroup := echo.Group("/b")
	boss(bGroup)
	gGroup := echo.Group("/images")
	get(gGroup)
}

func boss(group *echo.Group) {

	group.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    []byte(conf.Config.JWT.Secret),
		SigningMethod: conf.Config.JWT.Algorithm,
	}))
	group.POST("/upload", upload)
	group.DELETE("/delete", del)
}

func get(group *echo.Group) {

	group.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	group.GET("/:folder0/:filename", getImage)
	group.GET("/:folder0/:folder1/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:filename", getImage)
}
