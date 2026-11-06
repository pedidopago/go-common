// Package core contém os códigos compartilhados sem depender de um framework HTTP.
package core

import "errors"

type Error error

var (
	ErrNotFound            Error = errors.New("not found")
	ErrInternalServerError Error = errors.New("internal server error")
)
