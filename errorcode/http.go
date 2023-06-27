package errorcode

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Error error

var (
	ErrBadRequest                    Error = fmt.Errorf("bad request")
	ErrUnauthorized                  Error = fmt.Errorf("unauthorized")
	ErrPaymentRequired               Error = fmt.Errorf("payment required")
	ErrForbidden                     Error = fmt.Errorf("forbidden")
	ErrNotFound                      Error = fmt.Errorf("not found")
	ErrMethodNotAllowed              Error = fmt.Errorf("method not allowed")
	ErrNotAcceptable                 Error = fmt.Errorf("not acceptable")
	ErrProxyAuthRequired             Error = fmt.Errorf("proxy authentication required")
	ErrRequestTimeout                Error = fmt.Errorf("request timeout")
	ErrConflict                      Error = fmt.Errorf("conflict")
	ErrGone                          Error = fmt.Errorf("gone")
	ErrLengthRequired                Error = fmt.Errorf("length required")
	ErrPreconditionFailed            Error = fmt.Errorf("precondition failed")
	ErrRequestEntityTooLarge         Error = fmt.Errorf("request entity too large")
	ErrURITooLong                    Error = fmt.Errorf("URI too long")
	ErrUnsupportedMediaType          Error = fmt.Errorf("unsupported media type")
	ErrRequestedRangeNotSatisfiable  Error = fmt.Errorf("requested range not satisfiable")
	ErrExpectationFailed             Error = fmt.Errorf("expectation failed")
	ErrTeapot                        Error = fmt.Errorf("i'm a teapot")
	ErrMisdirectedRequest            Error = fmt.Errorf("misdirected request")
	ErrUnprocessableEntity           Error = fmt.Errorf("unprocessable entity")
	ErrLocked                        Error = fmt.Errorf("locked")
	ErrFailedDependency              Error = fmt.Errorf("failed dependency")
	ErrTooEarly                      Error = fmt.Errorf("too early")
	ErrUpgradeRequired               Error = fmt.Errorf("upgrade required")
	ErrPreconditionRequired          Error = fmt.Errorf("precondition required")
	ErrTooManyRequests               Error = fmt.Errorf("too many requests")
	ErrRequestHeaderFieldsTooLarge   Error = fmt.Errorf("request header fields too large")
	ErrUnavailableForLegalReasons    Error = fmt.Errorf("unavailable for legal reasons")
	ErrInternalServerError           Error = fmt.Errorf("internal server error")
	ErrNotImplemented                Error = fmt.Errorf("not implemented")
	ErrBadGateway                    Error = fmt.Errorf("bad gateway")
	ErrServiceUnavailable            Error = fmt.Errorf("service unavailable")
	ErrGatewayTimeout                Error = fmt.Errorf("gateway timeout")
	ErrHTTPVersionNotSupported       Error = fmt.Errorf("HTTP version not supported")
	ErrVariantAlsoNegotiates         Error = fmt.Errorf("variant also negotiates")
	ErrInsufficientStorage           Error = fmt.Errorf("insufficient storage")
	ErrLoopDetected                  Error = fmt.Errorf("loop detected")
	ErrNotExtended                   Error = fmt.Errorf("not extended")
	ErrNetworkAuthenticationRequired Error = fmt.Errorf("network authentication required")
)

func NewError(code Error, msg string) error {
	return fmt.Errorf("%s: %s", code.Error(), msg)
}

func CodeFromError(err error) Error {
	if !strings.Contains(err.Error(), ":") {
		return err
	}
	return errors.New(strings.Split(err.Error(), ":")[0])
}

func MsgFromError(err error) string {
	if !strings.Contains(err.Error(), ":") {
		return ""
	}
	return strings.Replace(strings.Split(err.Error(), ":")[1], " ", "", 1)
}

func GetHTTPStatus(err Error) int {
	switch CodeFromError(err) {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrPaymentRequired:
		return http.StatusPaymentRequired
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case ErrNotAcceptable:
		return http.StatusNotAcceptable
	case ErrProxyAuthRequired:
		return http.StatusProxyAuthRequired
	case ErrRequestTimeout:
		return http.StatusRequestTimeout
	case ErrConflict:
		return http.StatusConflict
	case ErrGone:
		return http.StatusGone
	case ErrLengthRequired:
		return http.StatusLengthRequired
	case ErrPreconditionFailed:
		return http.StatusPreconditionFailed
	case ErrRequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge
	case ErrURITooLong:
		return http.StatusRequestURITooLong
	case ErrUnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case ErrRequestedRangeNotSatisfiable:
		return http.StatusRequestedRangeNotSatisfiable
	case ErrExpectationFailed:
		return http.StatusExpectationFailed
	case ErrTeapot:
		return http.StatusTeapot
	case ErrMisdirectedRequest:
		return http.StatusMisdirectedRequest
	case ErrUnprocessableEntity:
		return http.StatusUnprocessableEntity
	case ErrLocked:
		return http.StatusLocked
	case ErrFailedDependency:
		return http.StatusFailedDependency
	case ErrTooEarly:
		return http.StatusTooEarly
	case ErrUpgradeRequired:
		return http.StatusUpgradeRequired
	case ErrPreconditionRequired:
		return http.StatusPreconditionRequired
	case ErrTooManyRequests:
		return http.StatusTooManyRequests
	case ErrRequestHeaderFieldsTooLarge:
		return http.StatusRequestHeaderFieldsTooLarge
	case ErrUnavailableForLegalReasons:
		return http.StatusUnavailableForLegalReasons
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotImplemented:
		return http.StatusNotImplemented
	case ErrBadGateway:
		return http.StatusBadGateway
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrGatewayTimeout:
		return http.StatusGatewayTimeout
	case ErrHTTPVersionNotSupported:
		return http.StatusHTTPVersionNotSupported
	case ErrVariantAlsoNegotiates:
		return http.StatusVariantAlsoNegotiates
	case ErrInsufficientStorage:
		return http.StatusInsufficientStorage
	case ErrLoopDetected:
		return http.StatusLoopDetected
	case ErrNotExtended:
		return http.StatusNotExtended
	case ErrNetworkAuthenticationRequired:
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusInternalServerError
	}
}

func GetEchoHTTPError(err Error) *echo.HTTPError {
	switch CodeFromError(err) {
	case ErrBadRequest:
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: MsgFromError(err),
		}
	case ErrUnauthorized:
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: MsgFromError(err),
		}
	case ErrPaymentRequired:
		return &echo.HTTPError{
			Code:    http.StatusPaymentRequired,
			Message: MsgFromError(err),
		}
	case ErrForbidden:
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: MsgFromError(err),
		}
	case ErrNotFound:
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: MsgFromError(err),
		}
	case ErrMethodNotAllowed:
		return &echo.HTTPError{
			Code:    http.StatusMethodNotAllowed,
			Message: MsgFromError(err),
		}
	case ErrNotAcceptable:
		return &echo.HTTPError{
			Code:    http.StatusNotAcceptable,
			Message: MsgFromError(err),
		}
	case ErrProxyAuthRequired:
		return &echo.HTTPError{
			Code:    http.StatusProxyAuthRequired,
			Message: MsgFromError(err),
		}
	case ErrRequestTimeout:
		return &echo.HTTPError{
			Code:    http.StatusRequestTimeout,
			Message: MsgFromError(err),
		}
	case ErrConflict:
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: MsgFromError(err),
		}
	case ErrGone:
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: MsgFromError(err),
		}
	case ErrLengthRequired:
		return &echo.HTTPError{
			Code:    http.StatusLengthRequired,
			Message: MsgFromError(err),
		}
	case ErrPreconditionFailed:
		return &echo.HTTPError{
			Code:    http.StatusPreconditionFailed,
			Message: MsgFromError(err),
		}
	case ErrRequestEntityTooLarge:
		return &echo.HTTPError{
			Code:    http.StatusRequestEntityTooLarge,
			Message: MsgFromError(err),
		}
	case ErrURITooLong:
		return &echo.HTTPError{
			Code:    http.StatusRequestURITooLong,
			Message: MsgFromError(err),
		}
	case ErrUnsupportedMediaType:
		return &echo.HTTPError{
			Code:    http.StatusUnsupportedMediaType,
			Message: MsgFromError(err),
		}
	case ErrRequestedRangeNotSatisfiable:
		return &echo.HTTPError{
			Code:    http.StatusRequestedRangeNotSatisfiable,
			Message: MsgFromError(err),
		}
	case ErrExpectationFailed:
		return &echo.HTTPError{
			Code:    http.StatusExpectationFailed,
			Message: MsgFromError(err),
		}
	case ErrTeapot:
		return &echo.HTTPError{
			Code:    http.StatusTeapot,
			Message: MsgFromError(err),
		}
	case ErrMisdirectedRequest:
		return &echo.HTTPError{
			Code:    http.StatusMisdirectedRequest,
			Message: MsgFromError(err),
		}
	case ErrUnprocessableEntity:
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: MsgFromError(err),
		}
	case ErrLocked:
		return &echo.HTTPError{
			Code:    http.StatusLocked,
			Message: MsgFromError(err),
		}
	case ErrFailedDependency:
		return &echo.HTTPError{
			Code:    http.StatusFailedDependency,
			Message: MsgFromError(err),
		}
	case ErrTooEarly:
		return &echo.HTTPError{
			Code:    http.StatusTooEarly,
			Message: MsgFromError(err),
		}
	case ErrUpgradeRequired:
		return &echo.HTTPError{
			Code:    http.StatusUpgradeRequired,
			Message: MsgFromError(err),
		}
	case ErrPreconditionRequired:
		return &echo.HTTPError{
			Code:    http.StatusPreconditionRequired,
			Message: MsgFromError(err),
		}
	case ErrTooManyRequests:
		return &echo.HTTPError{
			Code:    http.StatusTooManyRequests,
			Message: MsgFromError(err),
		}
	case ErrRequestHeaderFieldsTooLarge:
		return &echo.HTTPError{
			Code:    http.StatusRequestHeaderFieldsTooLarge,
			Message: MsgFromError(err),
		}
	case ErrUnavailableForLegalReasons:
		return &echo.HTTPError{
			Code:    http.StatusUnavailableForLegalReasons,
			Message: MsgFromError(err),
		}
	case ErrInternalServerError:
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: MsgFromError(err),
		}
	case ErrNotImplemented:
		return &echo.HTTPError{
			Code:    http.StatusNotImplemented,
			Message: MsgFromError(err),
		}
	case ErrBadGateway:
		return &echo.HTTPError{
			Code:    http.StatusBadGateway,
			Message: MsgFromError(err),
		}
	case ErrServiceUnavailable:
		return &echo.HTTPError{
			Code:    http.StatusServiceUnavailable,
			Message: MsgFromError(err),
		}
	case ErrGatewayTimeout:
		return &echo.HTTPError{
			Code:    http.StatusGatewayTimeout,
			Message: MsgFromError(err),
		}
	case ErrHTTPVersionNotSupported:
		return &echo.HTTPError{
			Code:    http.StatusHTTPVersionNotSupported,
			Message: MsgFromError(err),
		}
	case ErrVariantAlsoNegotiates:
		return &echo.HTTPError{
			Code:    http.StatusVariantAlsoNegotiates,
			Message: MsgFromError(err),
		}
	case ErrInsufficientStorage:
		return &echo.HTTPError{
			Code:    http.StatusInsufficientStorage,
			Message: MsgFromError(err),
		}
	case ErrLoopDetected:
		return &echo.HTTPError{
			Code:    http.StatusLoopDetected,
			Message: MsgFromError(err),
		}
	case ErrNotExtended:
		return &echo.HTTPError{
			Code:    http.StatusNotExtended,
			Message: MsgFromError(err),
		}
	case ErrNetworkAuthenticationRequired:
		return &echo.HTTPError{
			Code:    http.StatusNetworkAuthenticationRequired,
			Message: MsgFromError(err),
		}
	default:
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: MsgFromError(err),
		}
	}
}
