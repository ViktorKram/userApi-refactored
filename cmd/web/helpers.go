package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/render"
)

type (
	ErrResponse struct {
		Err            error `json:"-"`
		HTTPStatusCode int   `json:"-"`

		StatusText string `json:"status"`
		AppCode    int64  `json:"code,omitempty"`
		ErrorText  string `json:"error,omitempty"`
	}
)

var (
	ErrUserNotFound = errors.New("user_not_found")
)

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func badRequestError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
