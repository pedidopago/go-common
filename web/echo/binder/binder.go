package binder

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
