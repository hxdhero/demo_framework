package handler

import (
	ut "github.com/go-playground/universal-translator"
	"lls_api/pkg/app"
	"lls_api/pkg/rerr"
	"lls_api/pkg/util/validate"
)

type BaseHandler struct {
	trans ut.Translator
}

func NewBaseHandler() (*BaseHandler, error) {
	trans, err := validate.InitTrans("zh")
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	return &BaseHandler{trans: trans}, nil
}

func (h BaseHandler) RegisterBaseService(s *app.Server) {
	v1 := s.Group("/1.0")
	v1.GET("/sae/checkpreload/", app.WrapHandler(func(ctx *app.Context) {
		ctx.Rest().HttpSuccessStr("successful")
	}))
}
