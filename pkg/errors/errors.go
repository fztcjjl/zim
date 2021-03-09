package errors

import (
	"encoding/json"
)

type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func New(args ...string) error {
	err := &Error{Code: ErrCustom}
	if len(args) >= 1 {
		err.Message = args[0]
	}
	if len(args) >= 2 {
		err.Detail = args[1]
	}
	return err
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Code = -1
		e.Detail = err
	}
	if e.Code == 0 && e.Message == "" {
		e.Code = -1
		e.Detail = err
	}
	return e
}

func newError(code int32, message string) *Error {
	return &Error{Code: code, Message: message}
}
