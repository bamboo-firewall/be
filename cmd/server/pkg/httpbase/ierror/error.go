package ierror

import "fmt"

type ErrorCode uint16
type ErrorName string

func NewError(code ErrorCode, name ErrorName, message string) *Error {
	return &Error{
		Code:    code,
		Name:    name,
		Message: message,
	}
}

type Error struct {
	Code    ErrorCode   `json:"code"`
	Name    ErrorName   `json:"name"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`

	SubError       *CoreError `json:"-"`
	HTTPStatusCode int        `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d:%s msg[%s] subErr[%s]", e.Code, e.Name, e.Message, e.SubError)
}

func (e *Error) SetSubError(subErr *CoreError) *Error {
	e.SubError = subErr
	return e
}

func (e *Error) SetHTTPStatus(status int) *Error {
	e.HTTPStatusCode = status
	return e
}

func (e *Error) SetDetail(detail interface{}) *Error {
	e.Detail = detail
	return e
}
