package code

import "github.com/muixstudio/clio/pkg/errorx"

var (
	BadRequest          = errorx.NewHttpErrorWithStatusCode(400, 400, "Bad Request")
	Unauthorized        = errorx.NewHttpErrorWithStatusCode(401, 401, "Unauthorized")
	Forbidden           = errorx.NewHttpErrorWithStatusCode(403, 403, "Forbidden")
	NotFound            = errorx.NewHttpErrorWithStatusCode(404, 404, "Not Found")
	Conflict            = errorx.NewHttpErrorWithStatusCode(409, 409, "Conflict")
	UnprocessableEntity = errorx.NewHttpErrorWithStatusCode(422, 422, "Unprocessable Entity")
	TooManyRequests     = errorx.NewHttpErrorWithStatusCode(429, 429, "Too Many Requests")

	internalServerError = errorx.NewHttpErrorWithStatusCode(500, 500, "Internal Server Error")
	ServiceUnavailable  = errorx.NewHttpErrorWithStatusCode(503, 503, "Service Unavailable")
	GatewayTimeout      = errorx.NewHttpErrorWithStatusCode(504, 504, "Gateway Timeout")
)
