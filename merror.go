package merror

import (
	"errors"
	"fmt"
)

const (
	NotFoundMessage            = "not found"
	ForbiddenMessage           = "forbidden"
	UnauthorizedMessage        = "unauthorized"
	InternalServerErrorMessage = "internal server error"
	UnprocessableMessage       = "unprocessable entity"
)

// you can register custom error codes
// below basic codes
var ErrNotFoundCode = 404
var ErrForbiddenCode = 403
var ErrUnauthorizedCode = 401
var ErrInternalServerCode = 501
var ErrUnprocessableEntity = 422
var DefaultMessageIfEmpty = ""
var UndefinedErrorCode = -1

func NewNotFound() *Error {
	return New(ErrNotFoundCode, nil, "")
}
func NewNotFoundWrap(err error) *Error {
	return New(ErrNotFoundCode, err, "")
}

func NewForbidden() *Error {
	return New(ErrForbiddenCode, nil, "")
}

func NewForbiddenWrap(err error) *Error {
	return New(ErrForbiddenCode, err, "")
}
func NewUnauthorized() *Error {
	return New(ErrUnauthorizedCode, nil, "")
}
func NewUnauthorizedWrap(err error) *Error {
	return New(ErrUnauthorizedCode, err, "")
}

func NewInternalServer() *Error {
	return New(ErrInternalServerCode, nil, "")
}
func NewInternalServerWrap(err error) *Error {
	return New(ErrInternalServerCode, err, "")
}

func NewUnprocessableEntity() *Error {
	return New(ErrUnprocessableEntity, nil, "")
}
func NewUnprocessableEntityWrap(err error) *Error {
	return New(ErrUnprocessableEntity, err, "")
}

func New(code int, casue error, comment string) *Error {
	return &Error{
		code, comment, casue,
	}
}

type Error struct {
	code int
	// public safe message about case of error
	comment string

	cause error
}

func (self *Error) SetCode(code int) *Error {
	self.code = code
	return self
}
func (self *Error) SetCommentf(format string, args ...interface{}) *Error {
	self.comment = fmt.Sprintf(format, args...)
	return self
}

func (self *Error) SetComment(comment string) *Error {
	self.comment = comment
	return self
}

func (self *Error) String() string {
	return format(self.Code(), self.Error(), self.Cause())
}

// for better testing
func format(code int, comment string, cause error) string {
	return fmt.Sprintf("code: %v, comment: %s cause: %v", code, comment, cause)
}

func (self *Error) Error() string {
	if len(self.comment) > 0 {
		return self.comment
	}
	return GetCode(self.code, DefaultMessageIfEmpty)
}

func (self *Error) Code() int {
	return self.code
}

func (self *Error) Cause() error {
	return self.cause
}

type CauseHolder interface {
	Cause() error
}

// Cause returns to the parent and returns the root cause of the error
// implements the "causer" interface.
func Cause(err error) error {
	if r, ok := err.(CauseHolder); ok && r.Cause() != nil {
		return Cause(r.Cause())
	}
	return err
}

// The UnWrap function takes an error and return unwraped error.
func UnWrap(err error) error {
	if r, ok := err.(CauseHolder); ok && r != nil {
		return r.Cause()
	}
	return err
}

var pErr = &Error{}

func IsNotFound(err error) bool {
	return IsCode(err, ErrNotFoundCode)
}

func IsForbidden(err error) bool {
	return IsCode(err, ErrForbiddenCode)
}

func IsUnauthorized(err error) bool {
	return IsCode(err, ErrUnauthorizedCode)
}

func IsInternalServerError(err error) bool {
	return IsCode(err, ErrInternalServerCode)
}

func IsUnprocessableEntity(err error) bool {
	return IsCode(err, ErrUnprocessableEntity)
}

func IsCode(err error, code int) bool {
	if err == nil {
		return code == UndefinedErrorCode
	}
	if r, ok := err.(*Error); ok {

		return r.Code() == code
	}
	return code == UndefinedErrorCode
}

func Code(err error) int {
	if err == nil {
		return UndefinedErrorCode
	}
	if r, ok := err.(*Error); ok {
		return r.Code()
	}
	return UndefinedErrorCode
}

func Is(err error) bool {
	return errors.As(err, &pErr)
}

var codesMapOfString = map[int]string{}

func SetCode(code int, value string) {
	codesMapOfString[code] = value
}

func GetCode(code int, dfault string) string {
	if v, ok := codesMapOfString[code]; ok {
		return v
	}
	return dfault
}

func init() {
	SetCode(ErrNotFoundCode, NotFoundMessage)
	SetCode(ErrForbiddenCode, ForbiddenMessage)
	SetCode(ErrUnauthorizedCode, UnauthorizedMessage)
	SetCode(ErrInternalServerCode, InternalServerErrorMessage)
	SetCode(ErrUnprocessableEntity, UnprocessableMessage)
}
