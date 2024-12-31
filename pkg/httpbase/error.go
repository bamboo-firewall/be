package httpbase

import (
	"context"
	"net/http"

	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

const (
	DefaultServerName = ierror.ErrorName("err_internal_server")
)

// ErrorCode with 5xx status. Begin Code from 1000 -> 1999
const (
	ErrorCodeInternalServer ierror.ErrorCode = iota + 1000
	ErrorCodeDatabase
)

// ErrorCode with 4xx status. Begin code from 2000 -> 2999
const (
	ErrorCodeBadRequest ierror.ErrorCode = iota + 2000
	ErrorCodeNotFound
	ErrorCodeValidateRequest
	ErrorCodeForBidden
	ErrorCodeUnauthorized
)

func toName(id ierror.ErrorCode) ierror.ErrorName {
	switch id {
	// 5xx code
	case ErrorCodeInternalServer:
		return DefaultServerName
	case ErrorCodeDatabase:
		return "err_database"
	// 4xx code
	case ErrorCodeBadRequest:
		return "err_bad_request"
	case ErrorCodeUnauthorized:
		return "err_unauthorized"
	case ErrorCodeNotFound:
		return "err_not_found"
	case ErrorCodeValidateRequest:
		return "err_validate_request"
	default:
		return "err_common"
	}
}

var (
	ErrNotFound = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeNotFound, msgID).SetHTTPStatus(http.StatusNotFound)
	}

	ErrBindRequest = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeBadRequest, msgID).SetHTTPStatus(http.StatusBadRequest)
	}

	ErrValidateRequest = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeBadRequest, msgID).SetHTTPStatus(http.StatusBadRequest)
	}

	ErrBadRequest = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeBadRequest, msgID).SetHTTPStatus(http.StatusBadRequest)
	}

	ErrDatabase = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeDatabase, msgID).SetHTTPStatus(http.StatusInternalServerError)
	}

	ErrInternal = func(ctx context.Context, msgID string) *ierror.Error {
		return newClientIError(ctx, ErrorCodeInternalServer, msgID).SetHTTPStatus(http.StatusInternalServerError)
	}
)

func newClientIError(ctx context.Context, errorCode ierror.ErrorCode, msgID string) *ierror.Error {
	return ierror.NewError(errorCode, toName(errorCode), msgID)
}
