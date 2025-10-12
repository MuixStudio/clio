package errorx

import "net/http"

type HttpError interface {
	ErrorX
	GetStatusCode() int
	GetCode() int
}

type httpErrorBody struct {
	Code    int
	Message string
}
type httpError struct {
	StatusCode int
	Body       httpErrorBody
}

func (e *httpError) Error() string {
	return e.Body.Message
}

func (e *httpError) Tag() string {
	return "http"
}

func (e *httpError) GetStatusCode() int {
	return e.StatusCode
}

func (e *httpError) GetCode() int {
	return e.Body.Code
}

func NewHttpError(code int, message string) HttpError {
	return &httpError{
		StatusCode: http.StatusInternalServerError,
		Body: httpErrorBody{
			Code:    code,
			Message: message,
		},
	}
}

func NewHttpErrorWithStatusCode(statusCode int, code int, message string) HttpError {
	return &httpError{
		StatusCode: statusCode,
		Body: httpErrorBody{
			Code:    code,
			Message: message,
		},
	}
}
