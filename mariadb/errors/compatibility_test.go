package errors_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/pedidopago/go-common/errorcode"
	dberrors "github.com/pedidopago/go-common/mariadb/errors"
)

// A assinatura pública continua compatível com errorcode.Error.
var _ func(error) errorcode.Error = dberrors.ToErrorCode

func TestToErrorCodeCompatibility(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want errorcode.Error
	}{
		{name: "nil"},
		{name: "sql no rows", err: sql.ErrNoRows, want: errorcode.ErrNotFound},
		{name: "wrapped sql no rows", err: fmt.Errorf("consulta: %w", sql.ErrNoRows), want: errorcode.ErrNotFound},
		{name: "no results", err: dberrors.ErrNoResults, want: errorcode.ErrNotFound},
		{name: "wrapped no results", err: fmt.Errorf("consulta: %w", dberrors.ErrNoResults), want: errorcode.ErrNotFound},
		{name: "record not found", err: errors.New("record not found"), want: errorcode.ErrNotFound},
		{name: "query failure", err: errors.New("falha sintética"), want: errorcode.ErrInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dberrors.ToErrorCode(tt.err)
			if got != tt.want {
				t.Fatalf("ToErrorCode = %v, esperado %v", got, tt.want)
			}
		})
	}
}

func TestDatabaseErrorHTTPAdapters(t *testing.T) {
	for _, tt := range []struct {
		err  error
		want int
	}{
		{err: sql.ErrNoRows, want: http.StatusNotFound},
		{err: errors.New("falha sintética"), want: http.StatusInternalServerError},
	} {
		code := dberrors.ToErrorCode(tt.err)
		if got := errorcode.GetHTTPStatus(code); got != tt.want {
			t.Fatalf("status HTTP = %d, esperado %d", got, tt.want)
		}
		if got := errorcode.GetEchoHTTPError(code).Code; got != tt.want {
			t.Fatalf("status Echo = %d, esperado %d", got, tt.want)
		}
	}
}
