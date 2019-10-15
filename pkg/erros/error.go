package erros

import "net/http"

type (
	Error struct {
		Code     int         `json:"-"`
		Message  interface{} `json:"message"`
		Internal error       `json:"-"`
	}
)

var (
	ErrNotFound           = NewError(http.StatusNotFound)
	ErrBadRequest         = NewError(http.StatusBadRequest)
	ErrServiceUnavailable = NewError(http.StatusServiceUnavailable)
	ErrMkdir              = NewError(http.StatusServiceUnavailable, "目录创建失败")
)

func NewError(code int, message ...interface{}) *Error {
	he := &Error{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}
