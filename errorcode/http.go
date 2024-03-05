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
		return Error(err)
	}
	return Error(errors.New(strings.Split(err.Error(), ":")[0]))
}

func MsgFromError(err error) string {
	if !strings.Contains(err.Error(), ":") {
		return err.Error()
	}

	return strings.TrimSpace(strings.SplitN(err.Error(), ":", 2)[1])
}

func GetHTTPStatus(err Error) int {
	switch CodeFromError(err).Error() {
	case ErrBadRequest.Error():
		return http.StatusBadRequest
	case ErrUnauthorized.Error():
		return http.StatusUnauthorized
	case ErrPaymentRequired.Error():
		return http.StatusPaymentRequired
	case ErrForbidden.Error():
		return http.StatusForbidden
	case ErrNotFound.Error():
		return http.StatusNotFound
	case ErrMethodNotAllowed.Error():
		return http.StatusMethodNotAllowed
	case ErrNotAcceptable.Error():
		return http.StatusNotAcceptable
	case ErrProxyAuthRequired.Error():
		return http.StatusProxyAuthRequired
	case ErrRequestTimeout.Error():
		return http.StatusRequestTimeout
	case ErrConflict.Error():
		return http.StatusConflict
	case ErrGone.Error():
		return http.StatusGone
	case ErrLengthRequired.Error():
		return http.StatusLengthRequired
	case ErrPreconditionFailed.Error():
		return http.StatusPreconditionFailed
	case ErrRequestEntityTooLarge.Error():
		return http.StatusRequestEntityTooLarge
	case ErrURITooLong.Error():
		return http.StatusRequestURITooLong
	case ErrUnsupportedMediaType.Error():
		return http.StatusUnsupportedMediaType
	case ErrRequestedRangeNotSatisfiable.Error():
		return http.StatusRequestedRangeNotSatisfiable
	case ErrExpectationFailed.Error():
		return http.StatusExpectationFailed
	case ErrTeapot.Error():
		return http.StatusTeapot
	case ErrMisdirectedRequest.Error():
		return http.StatusMisdirectedRequest
	case ErrUnprocessableEntity.Error():
		return http.StatusUnprocessableEntity
	case ErrLocked.Error():
		return http.StatusLocked
	case ErrFailedDependency.Error():
		return http.StatusFailedDependency
	case ErrTooEarly.Error():
		return http.StatusTooEarly
	case ErrUpgradeRequired.Error():
		return http.StatusUpgradeRequired
	case ErrPreconditionRequired.Error():
		return http.StatusPreconditionRequired
	case ErrTooManyRequests.Error():
		return http.StatusTooManyRequests
	case ErrRequestHeaderFieldsTooLarge.Error():
		return http.StatusRequestHeaderFieldsTooLarge
	case ErrUnavailableForLegalReasons.Error():
		return http.StatusUnavailableForLegalReasons
	case ErrInternalServerError.Error():
		return http.StatusInternalServerError
	case ErrNotImplemented.Error():
		return http.StatusNotImplemented
	case ErrBadGateway.Error():
		return http.StatusBadGateway
	case ErrServiceUnavailable.Error():
		return http.StatusServiceUnavailable
	case ErrGatewayTimeout.Error():
		return http.StatusGatewayTimeout
	case ErrHTTPVersionNotSupported.Error():
		return http.StatusHTTPVersionNotSupported
	case ErrVariantAlsoNegotiates.Error():
		return http.StatusVariantAlsoNegotiates
	case ErrInsufficientStorage.Error():
		return http.StatusInsufficientStorage
	case ErrLoopDetected.Error():
		return http.StatusLoopDetected
	case ErrNotExtended.Error():
		return http.StatusNotExtended
	case ErrNetworkAuthenticationRequired.Error():
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusInternalServerError
	}
}

func GetEchoHTTPError(err Error) *echo.HTTPError {
	switch CodeFromError(err).Error() {
	case ErrBadRequest.Error():
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: MsgFromError(err),
		}
	case ErrUnauthorized.Error():
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: MsgFromError(err),
		}
	case ErrPaymentRequired.Error():
		return &echo.HTTPError{
			Code:    http.StatusPaymentRequired,
			Message: MsgFromError(err),
		}
	case ErrForbidden.Error():
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: MsgFromError(err),
		}
	case ErrNotFound.Error():
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: MsgFromError(err),
		}
	case ErrMethodNotAllowed.Error():
		return &echo.HTTPError{
			Code:    http.StatusMethodNotAllowed,
			Message: MsgFromError(err),
		}
	case ErrNotAcceptable.Error():
		return &echo.HTTPError{
			Code:    http.StatusNotAcceptable,
			Message: MsgFromError(err),
		}
	case ErrProxyAuthRequired.Error():
		return &echo.HTTPError{
			Code:    http.StatusProxyAuthRequired,
			Message: MsgFromError(err),
		}
	case ErrRequestTimeout.Error():
		return &echo.HTTPError{
			Code:    http.StatusRequestTimeout,
			Message: MsgFromError(err),
		}
	case ErrConflict.Error():
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: MsgFromError(err),
		}
	case ErrGone.Error():
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: MsgFromError(err),
		}
	case ErrLengthRequired.Error():
		return &echo.HTTPError{
			Code:    http.StatusLengthRequired,
			Message: MsgFromError(err),
		}
	case ErrPreconditionFailed.Error():
		return &echo.HTTPError{
			Code:    http.StatusPreconditionFailed,
			Message: MsgFromError(err),
		}
	case ErrRequestEntityTooLarge.Error():
		return &echo.HTTPError{
			Code:    http.StatusRequestEntityTooLarge,
			Message: MsgFromError(err),
		}
	case ErrURITooLong.Error():
		return &echo.HTTPError{
			Code:    http.StatusRequestURITooLong,
			Message: MsgFromError(err),
		}
	case ErrUnsupportedMediaType.Error():
		return &echo.HTTPError{
			Code:    http.StatusUnsupportedMediaType,
			Message: MsgFromError(err),
		}
	case ErrRequestedRangeNotSatisfiable.Error():
		return &echo.HTTPError{
			Code:    http.StatusRequestedRangeNotSatisfiable,
			Message: MsgFromError(err),
		}
	case ErrExpectationFailed.Error():
		return &echo.HTTPError{
			Code:    http.StatusExpectationFailed,
			Message: MsgFromError(err),
		}
	case ErrTeapot.Error():
		return &echo.HTTPError{
			Code:    http.StatusTeapot,
			Message: MsgFromError(err),
		}
	case ErrMisdirectedRequest.Error():
		return &echo.HTTPError{
			Code:    http.StatusMisdirectedRequest,
			Message: MsgFromError(err),
		}
	case ErrUnprocessableEntity.Error():
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: MsgFromError(err),
		}
	case ErrLocked.Error():
		return &echo.HTTPError{
			Code:    http.StatusLocked,
			Message: MsgFromError(err),
		}
	case ErrFailedDependency.Error():
		return &echo.HTTPError{
			Code:    http.StatusFailedDependency,
			Message: MsgFromError(err),
		}
	case ErrTooEarly.Error():
		return &echo.HTTPError{
			Code:    http.StatusTooEarly,
			Message: MsgFromError(err),
		}
	case ErrUpgradeRequired.Error():
		return &echo.HTTPError{
			Code:    http.StatusUpgradeRequired,
			Message: MsgFromError(err),
		}
	case ErrPreconditionRequired.Error():
		return &echo.HTTPError{
			Code:    http.StatusPreconditionRequired,
			Message: MsgFromError(err),
		}
	case ErrTooManyRequests.Error():
		return &echo.HTTPError{
			Code:    http.StatusTooManyRequests,
			Message: MsgFromError(err),
		}
	case ErrRequestHeaderFieldsTooLarge.Error():
		return &echo.HTTPError{
			Code:    http.StatusRequestHeaderFieldsTooLarge,
			Message: MsgFromError(err),
		}
	case ErrUnavailableForLegalReasons.Error():
		return &echo.HTTPError{
			Code:    http.StatusUnavailableForLegalReasons,
			Message: MsgFromError(err),
		}
	case ErrInternalServerError.Error():
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: MsgFromError(err),
		}
	case ErrNotImplemented.Error():
		return &echo.HTTPError{
			Code:    http.StatusNotImplemented,
			Message: MsgFromError(err),
		}
	case ErrBadGateway.Error():
		return &echo.HTTPError{
			Code:    http.StatusBadGateway,
			Message: MsgFromError(err),
		}
	case ErrServiceUnavailable.Error():
		return &echo.HTTPError{
			Code:    http.StatusServiceUnavailable,
			Message: MsgFromError(err),
		}
	case ErrGatewayTimeout.Error():
		return &echo.HTTPError{
			Code:    http.StatusGatewayTimeout,
			Message: MsgFromError(err),
		}
	case ErrHTTPVersionNotSupported.Error():
		return &echo.HTTPError{
			Code:    http.StatusHTTPVersionNotSupported,
			Message: MsgFromError(err),
		}
	case ErrVariantAlsoNegotiates.Error():
		return &echo.HTTPError{
			Code:    http.StatusVariantAlsoNegotiates,
			Message: MsgFromError(err),
		}
	case ErrInsufficientStorage.Error():
		return &echo.HTTPError{
			Code:    http.StatusInsufficientStorage,
			Message: MsgFromError(err),
		}
	case ErrLoopDetected.Error():
		return &echo.HTTPError{
			Code:    http.StatusLoopDetected,
			Message: MsgFromError(err),
		}
	case ErrNotExtended.Error():
		return &echo.HTTPError{
			Code:    http.StatusNotExtended,
			Message: MsgFromError(err),
		}
	case ErrNetworkAuthenticationRequired.Error():
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
