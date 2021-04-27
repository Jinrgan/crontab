package master

import (
	"go.uber.org/zap"
	"net/http"
	"os"
)

type Wrapper struct {
	Logger *zap.Logger
}

type appHandler func(w http.ResponseWriter, r *http.Request) error

func (w *Wrapper) WrapErr(h appHandler) func(
	http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		// panic
		defer func() {
			if r := recover(); r != nil {
				w.Logger.Error("Panic", zap.Any("recover", r))
				http.Error(
					writer,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		err := h(writer, req)

		if err != nil {
			//w.Logger.Error("Error occurred handling request", zap.Error(err))

			// user error
			if userErr, ok := err.(userError); ok {
				http.Error(writer,
					userErr.Message(),
					http.StatusBadRequest)
				return
			}

			// system error
			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}

			http.Error(writer, http.StatusText(code), code)
		}
	}
}

type userError interface {
	error
	Message() string
}
