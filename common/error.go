package common

import "encoding/json"

type StandardDetail struct {
	Msg        string `json:"msg"`
	Code       string `json:"code"`
	XRequestID string `json:"x_request_id"`
}

type DisplayError struct {
	OriginErr      error          `json:"-"`
	StandardDetail StandardDetail `json:"standard_detail"`
}

func (e DisplayError) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e DisplayError) Error() string {
	return e.StandardDetail.Msg
}

func (e DisplayError) OriginalErr() error {
	return e.OriginErr
}
