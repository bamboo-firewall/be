package client

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bamboo-firewall/be/pkg/common/errlist"
	"github.com/bamboo-firewall/be/pkg/httpbase"
	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

type apiServer struct {
	client *httpbase.Client
}

func NewAPIServer(address string) *apiServer {
	return &apiServer{client: httpbase.NewClient(address)}
}

func responseBodyToIError(ctx context.Context, res *httpbase.Result) *ierror.Error {
	errorResponse := new(httpbase.ErrorResponse)
	err := json.Unmarshal(res.Body, errorResponse)
	if err != nil {
		return httpbase.ErrInternal(ctx, "undefined response error").
			SetSubError(
				errlist.ErrUnmarshalFailed.WithChild(errors.New(string(res.Body))),
			)
	}
	return &ierror.Error{
		Code:           errorResponse.Error.Code,
		Name:           errorResponse.Error.Name,
		Message:        errorResponse.Error.Message,
		Detail:         errorResponse.Error.Detail,
		HTTPStatusCode: res.StatusCode,
	}
}
