package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type (
	// Error 异常结构
	Error struct {
		Code    int         `json:"code"`
		Key     string      `json:"error,omitempty"`
		Success bool        `json:"success"`
		Message interface{} `json:"message"`
	}
	// Result 返回结果
	Result struct {
		Code    int         `json:"code"`
		Success bool        `json:"success"`
		Message interface{} `json:"message,omitempty"`
		Data    interface{} `json:"file,omitempty"`
	}
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	Register(e)
	return e
}

// JSON 输出json
func R(c echo.Context, result *Result) error {
	return c.JSON(result.Code, result)
}
func E(c echo.Context, result *Error) error {
	return c.JSON(result.Code, result)
}

// ResultOf 构造Result
func ResultOf(code int, data interface{}) *Result {
	return &Result{
		Code:    code,
		Success: true,
		Data:    data,
	}
}
func ErrorOf(code int, message interface{}) *Error {
	return &Error{
		Code:    code,
		Success: false,
		Message: message,
	}
}
