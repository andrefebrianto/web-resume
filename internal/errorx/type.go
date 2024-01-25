package errorx

import (
	"errors"
	"net/http"
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	TypeInvalidParameter   Type = "BAD_REQUEST"
	TypeUnauthorized       Type = "UNAUTHORIZED"
	TypeForbidden          Type = "FORBIDDEN"
	TypeNotFound           Type = "NOT_FOUND"
	TypeInternal           Type = "INTERNAL_ERROR"
	TypeNotImplemented     Type = "NOT_IMPLEMENTED"
	TypeBadGateway         Type = "BAD_GATEWAY"
	TypeServiceUnavailable Type = "SERVICE_UNAVAILABLE"
	TypeNotMatch           Type = "NOT_MATCH"
)

var TypeToCode = map[Type]int{
	TypeInvalidParameter:   http.StatusBadRequest,
	TypeUnauthorized:       http.StatusUnauthorized,
	TypeForbidden:          http.StatusForbidden,
	TypeNotFound:           http.StatusNotFound,
	TypeInternal:           http.StatusInternalServerError,
	TypeNotImplemented:     http.StatusNotImplemented,
	TypeBadGateway:         http.StatusBadGateway,
	TypeServiceUnavailable: http.StatusServiceUnavailable,
	TypeNotMatch:           http.StatusBadRequest,
}

// Predefined Errors
var (
	ErrInvalidParameter = &Error{
		Type:    TypeInvalidParameter,
		Message: TypeInvalidParameter.String(),
		Code:    TypeToCode[TypeInvalidParameter],
		Err:     errors.New(TypeInvalidParameter.String()),
	}
	ErrUnauthorized = &Error{
		Type:    TypeUnauthorized,
		Message: TypeUnauthorized.String(),
		Code:    TypeToCode[TypeUnauthorized],
		Err:     errors.New(TypeUnauthorized.String()),
	}
	ErrForbidden = &Error{
		Type:    TypeForbidden,
		Message: TypeForbidden.String(),
		Code:    TypeToCode[TypeForbidden],
		Err:     errors.New(TypeForbidden.String()),
	}
	ErrNotFound = &Error{
		Type:    TypeNotFound,
		Message: TypeNotFound.String(),
		Code:    TypeToCode[TypeNotFound],
		Err:     errors.New(TypeNotFound.String()),
	}
	ErrInternal = &Error{
		Type:    TypeInternal,
		Message: TypeInternal.String(),
		Code:    TypeToCode[TypeInternal],
		Err:     errors.New(TypeInternal.String()),
	}
	ErrNotImplemented = &Error{
		Type:    TypeNotImplemented,
		Message: TypeNotImplemented.String(),
		Code:    TypeToCode[TypeNotImplemented],
		Err:     errors.New(TypeNotImplemented.String()),
	}
	ErrBadGateway = &Error{
		Type:    TypeBadGateway,
		Message: TypeBadGateway.String(),
		Code:    TypeToCode[TypeBadGateway],
		Err:     errors.New(TypeBadGateway.String()),
	}
	ErrServiceUnavailable = &Error{
		Type:    TypeServiceUnavailable,
		Message: TypeServiceUnavailable.String(),
		Code:    TypeToCode[TypeServiceUnavailable],
		Err:     errors.New(TypeServiceUnavailable.String()),
	}
	ErrNotMatch = &Error{
		Type:    TypeNotMatch,
		Message: TypeNotMatch.String(),
		Code:    TypeToCode[TypeNotMatch],
		Err:     errors.New(TypeNotMatch.String()),
	}
)
