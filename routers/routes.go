package routers

import (
	"github.com/buzzxu/boys/types"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"image-server/pkg/conf"
	"net/http"
	"time"
)

func Register(echo *echo.Echo) {

	bGroup := echo.Group("/b")
	boss(bGroup)
	gGroup := echo.Group("/images")
	images(gGroup)
	funsGroup := echo.Group("/funs")
	funs(funsGroup)

	qrCodeGroup := echo.Group("/qrCode")
	funs(qrCodeGroup)
}

var jwt = echojwt.JWT(echojwt.Config{
	SigningKey:    []byte(conf.Config.JWT.Secret),
	SigningMethod: conf.Config.JWT.Algorithm,
})
var cors = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
})

var rateLimiter = middleware.RateLimiterConfig{
	Skipper: middleware.DefaultSkipper,
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
	),
	ErrorHandler: func(context echo.Context, err error) error {
		return context.JSON(http.StatusForbidden, types.Result{Code: http.StatusForbidden, Message: "未授权操作"})
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return context.JSON(http.StatusTooManyRequests, types.Result{Code: http.StatusTooManyRequests, Message: "已被限流，请稍后尝试"})
	},
}

func boss(group *echo.Group) {
	group.Use(cors, jwt)
	group.POST("/upload", upload, middleware.BodyLimit(conf.Config.BodyLimit))
	group.DELETE("/delete", del)
}

func images(group *echo.Group) {

	group.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	group.GET("/:folder0/:filename", getImage)
	group.GET("/:folder0/:folder1/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:folder4/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:folder4/:folder5/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:folder4/:folder5/:folder6/:filename", getImage)
	group.GET("/:folder0/:folder1/:folder2/:folder3/:folder4/:folder5/:folder6/:folder7/:filename", getImage)
	group.POST("/upload", upload, jwt, middleware.BodyLimit(conf.Config.BodyLimit))
	group.DELETE("/delete", del, jwt)
}

func funs(group *echo.Group) {
	//裁剪
	group.POST("/crop", crop)
	//合并
	group.POST("/composite", composite)
}

func qrCode(group *echo.Group) {

}
