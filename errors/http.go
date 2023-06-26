package errorcode

import (
	"fmt"
	"net/http"
)

var (
	ErrBadRequest                    = fmt.Errorf("bad request")
	ErrUnauthorized                  = fmt.Errorf("unauthorized")
	ErrPaymentRequired               = fmt.Errorf("payment required")
	ErrForbidden                     = fmt.Errorf("forbidden")
	ErrNotFound                      = fmt.Errorf("not found")
	ErrMethodNotAllowed              = fmt.Errorf("method not allowed")
	ErrNotAcceptable                 = fmt.Errorf("not acceptable")
	ErrProxyAuthRequired             = fmt.Errorf("proxy authentication required")
	ErrRequestTimeout                = fmt.Errorf("request timeout")
	ErrConflict                      = fmt.Errorf("conflict")
	ErrGone                          = fmt.Errorf("gone")
	ErrLengthRequired                = fmt.Errorf("length required")
	ErrPreconditionFailed            = fmt.Errorf("precondition failed")
	ErrRequestEntityTooLarge         = fmt.Errorf("request entity too large")
	ErrURITooLong                    = fmt.Errorf("URI too long")
	ErrUnsupportedMediaType          = fmt.Errorf("unsupported media type")
	ErrRequestedRangeNotSatisfiable  = fmt.Errorf("requested range not satisfiable")
	ErrExpectationFailed             = fmt.Errorf("expectation failed")
	ErrTeapot                        = fmt.Errorf("i'm a teapot")
	ErrMisdirectedRequest            = fmt.Errorf("misdirected request")
	ErrUnprocessableEntity           = fmt.Errorf("unprocessable entity")
	ErrLocked                        = fmt.Errorf("locked")
	ErrFailedDependency              = fmt.Errorf("failed dependency")
	ErrTooEarly                      = fmt.Errorf("too early")
	ErrUpgradeRequired               = fmt.Errorf("upgrade required")
	ErrPreconditionRequired          = fmt.Errorf("precondition required")
	ErrTooManyRequests               = fmt.Errorf("too many requests")
	ErrRequestHeaderFieldsTooLarge   = fmt.Errorf("request header fields too large")
	ErrUnavailableForLegalReasons    = fmt.Errorf("unavailable for legal reasons")
	ErrInternalServerError           = fmt.Errorf("internal server error")
	ErrNotImplemented                = fmt.Errorf("not implemented")
	ErrBadGateway                    = fmt.Errorf("bad gateway")
	ErrServiceUnavailable            = fmt.Errorf("service unavailable")
	ErrGatewayTimeout                = fmt.Errorf("gateway timeout")
	ErrHTTPVersionNotSupported       = fmt.Errorf("HTTP version not supported")
	ErrVariantAlsoNegotiates         = fmt.Errorf("variant also negotiates")
	ErrInsufficientStorage           = fmt.Errorf("insufficient storage")
	ErrLoopDetected                  = fmt.Errorf("loop detected")
	ErrNotExtended                   = fmt.Errorf("not extended")
	ErrNetworkAuthenticationRequired = fmt.Errorf("network authentication required")
)

func GetHTTPStatus(err error) int {
	switch err {
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
