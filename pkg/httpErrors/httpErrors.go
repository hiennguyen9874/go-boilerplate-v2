package httpErrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
)

var (
	ErrorBadRequest                = errors.New("bad_request")
	ErrorNotFound                  = errors.New("not_found")
	ErrorUnauthorized              = errors.New("unauthorized")
	ErrorInternalServerError       = errors.New("internal_server_error")
	ErrorRequestTimeoutError       = errors.New("request_timeout")
	ErrorExistsEmailError          = errors.New("email_exists")
	ErrorInvalidJWTToken           = errors.New("invalid_jwt_token")
	ErrorInvalidJWTClaims          = errors.New("invalid_jwt_claims")
	ErrorValidation                = errors.New("validation")
	ErrorWrongPassword             = errors.New("wrong_password")
	ErrorTokenNotFound             = errors.New("token_not_found")
	ErrorInactiveUser              = errors.New("inactive_user")
	ErrorNotEnoughPrivileges       = errors.New("not_enough_privileges")
	ErrorGenToken                  = errors.New("generate_token_error")
	ErrorJson                      = errors.New("error_json_marshal")
	ErrorNotFoundRefreshTokenRedis = errors.New("not_found_refresh_token_redis")
	ErrorUserAlreadyVerified       = errors.New("user_already_verified")
	ErrorUserNotVerified           = errors.New("user_not_verified")
)

// Rest error interface
type ErrRest interface {
	GetErr() error
	GetStatus() int
	GetStatusText() string
	GetMsg() string
	Error() string
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err        error  `json:"-"`                                      // low-level runtime error
	Status     int    `json:"status" example:"404"`                   // http response status code
	StatusText string `json:"statusText" example:"not_found"`         // user-level status message
	Msg        string `json:"msg,omitempty" example:"not found user"` // application-level error message, for debugging
}

func (e *ErrResponse) GetErr() error {
	return e.Err
}

func (e *ErrResponse) GetStatus() int {
	return e.Status
}

func (e *ErrResponse) GetStatusText() string {
	return e.StatusText
}

func (e *ErrResponse) GetMsg() string {
	return e.Msg
}

// Error Error() interface method
func (e *ErrResponse) Error() string {
	return fmt.Sprintf("status: %d - statusText: %s - msg: %s - error: %v", e.Status, e.StatusText, e.Msg, e.Err)
}

func ErrBadRequest(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorBadRequest.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusNotFound,
		StatusText: ErrorNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrUnauthorized(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUnauthorized.Error(),
		Msg:        err.Error(),
	}
}

func ErrInternalServer(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusInternalServerError,
		StatusText: ErrorInternalServerError.Error(),
		Msg:        err.Error(),
	}
}

func ErrValidation(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnprocessableEntity,
		StatusText: ErrorValidation.Error(),
		Msg:        err.Error(),
	}
}

func ErrRequestTimeoutError(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusRequestTimeout,
		StatusText: ErrorRequestTimeoutError.Error(),
		Msg:        err.Error(),
	}
}

func ErrInactiveUser(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorInactiveUser.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotEnoughPrivileges(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusForbidden,
		StatusText: ErrorNotEnoughPrivileges.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrInvalidJWTClaims(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorInvalidJWTClaims.Error(),
		Msg:        err.Error(),
	}
}

func ErrWrongPassword(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorWrongPassword.Error(),
		Msg:        err.Error(),
	}
}

func ErrGenToken(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusBadRequest,
		StatusText: ErrorGenToken.Error(),
		Msg:        err.Error(),
	}
}

func ErrTokenNotFound(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorTokenNotFound.Error(),
		Msg:        err.Error(),
	}
}

func ErrJson(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorJson.Error(),
		Msg:        err.Error(),
	}
}

func ErrNotFoundRefreshTokenRedis(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorNotFoundRefreshTokenRedis.Error(),
		Msg:        err.Error(),
	}
}

func ErrUserAlreadyVerified(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUserAlreadyVerified.Error(),
		Msg:        err.Error(),
	}
}

func ErrUserNotVerified(err error) ErrRest {
	return &ErrResponse{
		Err:        err,
		Status:     http.StatusUnauthorized,
		StatusText: ErrorUserNotVerified.Error(),
		Msg:        err.Error(),
	}
}

// Parser of error string messages ,returns RestError
func ParseErrors(err error) ErrRest {
	switch {
	case ent.IsNotFound(err):
		return ErrNotFound(err)
	case errors.Is(err, context.DeadlineExceeded):
		return ErrRequestTimeoutError(err)
	default:
		if restErr, ok := err.(ErrRest); ok {
			return restErr
		}
		return ErrBadRequest(err)
	}
}
