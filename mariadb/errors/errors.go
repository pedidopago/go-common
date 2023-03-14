package errors

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidQuery = errors.New("invalid query")
	ErrNoResults    = errors.New("no results")
	ErrQueryResult  = errors.New("query result error")
)

func CustomInvalidInput(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, msg)
}

func InvalidQuery(original error) error {
	if original == nil {
		return ErrInvalidQuery
	}
	return fmt.Errorf("%w: %s", ErrInvalidQuery, original.Error())
}

// WrapSQLX wraps a sqlx error
func WrapSQLX(original error) error {
	if original == nil {
		return nil
	}
	if original == sql.ErrNoRows {
		return ErrNoResults
	}
	return fmt.Errorf("%w: %s", ErrQueryResult, original.Error())
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// New returns an error that formats as the given text. Each call to New returns
// a distinct error value even if the text is identical.
func New(text string) error {
	return errors.New(text)
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	if Is(err, ErrNoResults) {
		return true
	}
	if Is(err, sql.ErrNoRows) {
		return true
	}
	estr := err.Error()
	if estr == "sql: no rows in result set" {
		return true
	}
	if estr == "record not found" {
		return true
	}
	if estr == "no results" {
		return true
	}
	return false
}
