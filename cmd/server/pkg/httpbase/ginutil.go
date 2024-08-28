package httpbase

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
)

type ErrorResponse struct {
	Error *MinifyError `json:"error"`
}

type MinifyError struct {
	Code       ierror.ErrorCode `json:"code"`
	Name       ierror.ErrorName `json:"name"`
	Message    string           `json:"message"`
	Detail     interface{}      `json:"detail,omitempty"`
	SubName    string           `json:"sub_name,omitempty"`
	SubMessage string           `json:"sub_message,omitempty"`
	TrackID    string           `json:"track_id,omitempty"`
}

func ReturnErrorResponse(ctx *gin.Context, err *ierror.Error) {
	minifyErr := &MinifyError{
		Code:    err.Code,
		Name:    err.Name,
		Message: err.Message,
		Detail:  err.Detail,
	}
	ctx.JSON(err.HTTPStatusCode, ErrorResponse{
		Error: minifyErr,
	})
}

func ReturnSuccessResponse(ctx *gin.Context, status int, data interface{}) {
	if data == nil {
		ctx.Status(status)
	} else {
		ctx.JSON(status, data)
	}
}

func BindInput(ctx *gin.Context, input interface{}) *ierror.Error {
	err := ctx.ShouldBindUri(input)
	if err != nil {
		return bindError(ctx, err)
	}
	err = ctx.ShouldBindHeader(input)
	if err != nil {
		return bindError(ctx, err)
	}
	err = ctx.ShouldBind(input)
	if err != nil {
		return bindError(ctx, err)
	}
	if err = defaultValidator.Struct(input); err != nil {
		return validateError(ctx, err)
	}
	return nil
}

func bindError(ctx *gin.Context, err error) *ierror.Error {
	return ErrBindRequest(ctx, "BindFailed").SetDetail(err.Error())
}

type validatorErrorDetail struct {
	Field         string      `json:"field"`
	RejectedValue interface{} `json:"rejected_value"`
	Tag           string      `json:"tag"`
	Param         string      `json:"param,omitempty"`
}

func validateError(ctx *gin.Context, err error) *ierror.Error {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		return ErrValidateRequest(ctx, "ValidateFailed").SetDetail(err.Error())
	}

	errDetails := make([]validatorErrorDetail, 0, len(errs))
	for _, verr := range errs {
		errDetails = append(errDetails, validatorErrorDetail{
			Field:         verr.Field(),
			RejectedValue: verr.Value(),
			Tag:           verr.Tag(),
			Param:         verr.Param(),
		})
	}
	return ErrValidateRequest(ctx, "ValidateFailed").SetDetail(errDetails)
}
