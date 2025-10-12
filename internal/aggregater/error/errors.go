package errors

import (
	"net/http"

	"github.com/muixstudio/clio/pkg/errorx"
)

var (
	ParseErr = parseErr{}
)

type parseErr struct {
	errorx.HttpError
}

func (e parseErr) WithParseErr() string {
	return e.HttpError.Error()
}

func NewParseErr(err error) error {
	return errorx.NewHttpErrorWithStatusCode(http.StatusBadRequest, 10, err.Error())
}
