package merror

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type info struct {
	code int
	err  error
}

func TestCode(t *testing.T) {
	assert.Equal(t, NotFoundMessage, GetCode(ErrNotFoundCode, "blank"))
	assert.Equal(t, ForbiddenMessage, GetCode(ErrForbiddenCode, "blank"))
	assert.Equal(t, UnauthorizedMessage, GetCode(ErrUnauthorizedCode, "blank"))
	assert.Equal(t, InternalServerErrorMessage, GetCode(ErrInternalServerCode, "blank"))
	assert.Equal(t, UnprocessableMessage, GetCode(ErrUnprocessableEntity, "blank"))
	assert.Equal(t, "blank", GetCode(UndefinedErrorCode, "blank"), "any error")
}

func TestIsCode(t *testing.T) {
	assert.Equal(t, true, IsCode(NewNotFound(), ErrNotFoundCode))
	assert.Equal(t, false, IsCode(NewNotFound(), ErrInternalServerCode))
	assert.Equal(t, false, IsCode(nil, ErrInternalServerCode))
	assert.Equal(t, true, IsCode(nil, UndefinedErrorCode))
	assert.Equal(t, true, IsCode(errors.New("another error"), UndefinedErrorCode))
	assert.Equal(t, false, IsCode(New(0, nil, "comment"), UndefinedErrorCode))

	assert.Equal(t, true, IsNotFound(NewNotFound()))
	assert.Equal(t, true, IsUnauthorized(NewUnauthorized()))
	assert.Equal(t, true, IsForbidden(NewForbidden()))
	assert.Equal(t, true, IsInternalServerError(NewInternalServer()))
	assert.Equal(t, true, IsUnprocessableEntity(NewUnprocessableEntity()))
}

func TestSet(t *testing.T) {
	causeErr := errors.New("cause error message")
	e := New(1, causeErr, "comment")

	assert.Equal(t, 1, e.Code())
	assert.Equal(t, causeErr, e.Cause())
	assert.Equal(t, "comment", e.Error())

	// change values
	assert.Equal(t, 2, e.SetCode(2).Code())
	assert.Equal(t, DefaultMessageIfEmpty, e.SetComment("").Error())
	assert.Equal(t, "comment 1", e.SetCommentf("comment %d", 1).Error())
	assert.Equal(t, "no comment", e.SetComment("no comment").Error())

	assert.Equal(t, format(1, "text", causeErr), New(1, causeErr, "text").String())
}

func TestNew(t *testing.T) {
	assert.Equal(t, ErrNotFoundCode, Code(NewNotFound()), "not found error code")
	assert.Equal(t, ErrForbiddenCode, Code(NewForbidden()), "forbidden error code")
	assert.Equal(t, ErrUnauthorizedCode, Code(NewUnauthorized()), "forbidden error code")
	assert.Equal(t, ErrInternalServerCode, Code(NewInternalServer()), "internal server error code")
	assert.Equal(t, ErrUnprocessableEntity, Code(NewUnprocessableEntity()), "unprocessable entity error code")
	assert.Equal(t, UndefinedErrorCode, Code(errors.New("any error")), "any error")
	assert.Equal(t, UndefinedErrorCode, Code(nil), "nil")
}

func TestWrapAndCause(t *testing.T) {
	wrapped := errors.New("wrapped error message")
	assert.Equal(t, info{ErrNotFoundCode, wrapped}, getInfo(NewNotFoundWrap(wrapped)), "not found error code")
	assert.Equal(t, info{ErrForbiddenCode, wrapped}, getInfo(NewForbiddenWrap(wrapped)), "forbidden error code")
	assert.Equal(t, info{ErrUnauthorizedCode, wrapped}, getInfo(NewUnauthorizedWrap(wrapped)), "forbidden error code")
	assert.Equal(t, info{ErrInternalServerCode, wrapped}, getInfo(NewInternalServerWrap(wrapped)), "internal server error code")
	assert.Equal(t, info{ErrUnprocessableEntity, wrapped}, getInfo(NewUnprocessableEntityWrap(wrapped)), "unprocessable entity error code")
}

func TestUnWrap(t *testing.T) {
	wrapped := errors.New("wrapped error message")
	assert.Equal(t, wrapped, UnWrap(NewNotFoundWrap(wrapped)), "not found error code")
	assert.Equal(t, wrapped, UnWrap(NewForbiddenWrap(wrapped)), "forbidden error code")
	assert.Equal(t, wrapped, UnWrap(NewUnauthorizedWrap(wrapped)), "forbidden error code")
	assert.Equal(t, wrapped, UnWrap(NewInternalServerWrap(wrapped)), "internal server error code")
	assert.Equal(t, wrapped, UnWrap(NewUnprocessableEntityWrap(wrapped)), "unprocessable entity error code")
	assert.Equal(t, nil, UnWrap(nil))

}

func getInfo(err error) info {
	if !Is(err) {
		return info{UndefinedErrorCode, nil}
	}
	return info{Code(err), Cause(err)}
}
