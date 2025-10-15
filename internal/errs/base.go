package errs

import (
	"encoding/json"
	"lls_api/pkg/rerr"
)

type BaseError struct {
	Err error  `json:"-"`
	Msg string `json:"errorMessage"`
}

func NewBaseError(msg string, errs ...error) BaseError {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	return BaseError{Msg: msg, Err: rerr.Wrap(err)}
}

func (e BaseError) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e BaseError) Error() string {
	return e.Msg
}
