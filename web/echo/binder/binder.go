package binder

import "github.com/labstack/echo/v4"

type Binder interface {
	QueryBinder
	PathBinder
	BodyBinder
	Bind(any) error
}

type QueryBinder interface {
	BindQueryParams(any) error
}

type PathBinder interface {
	BindPathParams(any) error
}

type BodyBinder interface {
	BindBody(any) error
}

type EchoContextBinder interface {
	EchoContextQueryBinder
	EchoContextPathBinder
	EchoContextBodyBinder
	Bind(any, echo.Context) error
}

type EchoContextQueryBinder interface {
	BindQueryParams(any, echo.Context) error
}

type EchoContextPathBinder interface {
	BindPathParams(any, echo.Context) error
}

type EchoContextBodyBinder interface {
	BindBody(any, echo.Context) error
}
