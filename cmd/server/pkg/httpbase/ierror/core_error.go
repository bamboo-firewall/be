package ierror

import "fmt"

func NewCoreError(name string, message string) *CoreError {
	return &CoreError{
		Name:    name,
		Message: message,
	}
}

type CoreError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Child   error  `json:"-"`
}

func (e *CoreError) Error() string {
	if e.Child != nil {
		return fmt.Sprintf("%s:%s child(%s)", e.Name, e.Message, e.Child)
	} else {
		return fmt.Sprintf("%s:%s", e.Name, e.Message)
	}
}

func (e *CoreError) WithChild(child error) *CoreError {
	e.Child = child
	return e
}
