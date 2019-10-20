package routers

import (
	"github.com/buzzxu/boys/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
)

type (
	Result struct {
		*types.Result
		File []string `json:"file"`
	}
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisablePrintStack: true,
		DisableStackAll:   true,
		StackSize:         4 << 10,
	}))
	Register(e)
	return e
}

// JSON 输出json
func R(c echo.Context, data []string) error {
	return c.JSON(http.StatusOK, &Result{Result: &types.Result{Code: http.StatusOK, Success: true}, File: data})
}
func RNullData(c echo.Context) error {
	return c.JSON(http.StatusOK, types.ResultNilData(http.StatusOK))
}
func E(c echo.Context, result *types.Error) error {
	return c.JSON(result.Code, result)
}

func jsonErrorHandler(err error, c echo.Context) error {
	var (
		code = http.StatusInternalServerError
		erro interface{}
	)
	if e, ok := err.(*types.Error); ok {
		code = e.Code
		erro = e
	} else if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
		erro = types.NewError(e.Code, e.Error())
	} else {
		erro = types.ErrorOf(err)
	}
	return c.JSON(code, erro)
}
