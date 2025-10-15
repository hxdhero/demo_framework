package errs

import (
	"lls_api/common"
	"lls_api/pkg/app"
)

func NewDisplayErr(ctx *app.Context, err error, msg, code string) common.DisplayError {
	return common.DisplayError{
		OriginErr: err,
		StandardDetail: common.StandardDetail{
			Msg:        msg,
			Code:       code,
			XRequestID: ctx.Session().RequestID,
		}}
}
