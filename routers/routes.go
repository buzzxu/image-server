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
	images(gGroup)
}

var jwt = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey:    []byte(conf.Config.JWT.Secret),
	SigningMethod: conf.Config.JWT.Algorithm,
})
var cors = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
})

func boss(group *echo.Group) {
	group.Use(cors, jwt)
	group.POST("/upload", upload)
	group.DELETE("/delete", del)
}

func images(group *echo.Group) {

	group.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	group.GET("/:folder0/:filename", getImage)
	group.GET("/:folder0/:folder1/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:filename", getImage)
	group.POST("/upload", upload, jwt)
	group.DELETE("/delete", del, jwt)
}
