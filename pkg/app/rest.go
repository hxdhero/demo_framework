package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lls_api/common"
	"lls_api/pkg/log"
	"lls_api/pkg/rerr"
	"net/http"
	"strconv"
)

type Rest struct {
	gin *gin.Context
	log *log.LoggerContext // 后续可以替换成接口
}

func NewRest(ginContext *gin.Context, log *log.LoggerContext) *Rest {
	return &Rest{gin: ginContext, log: log}
}

// == 业务方法
func (r *Rest) HttpError(err error) {
	if err != nil {
		r.gin.JSON(http.StatusBadRequest, err)
		r.log.InfofSkip(1, "err: %s", err.Error())
		var oe IOriginalError
		if errors.As(err, &oe) {
			if oe.OriginalErr() != nil {
				r.log.InfofSkip(1, "err: %s", oe.OriginalErr().Error())
			}
		}
		return
	}
	r.gin.JSON(http.StatusBadRequest, "未知错误.")
}

func (r *Rest) HttpSuccess(data any) {
	r.gin.JSON(http.StatusOK, data)
}

func (r *Rest) HttpSuccessStr(str string) {
	r.Gin().String(http.StatusOK, str)
}

func (r *Rest) bizErrWithErr(err error) {
	if err == nil {
		r.Gin().JSON(http.StatusInternalServerError, "服务器内部错误.")
		return
	}
	format := rerr.NewDefaultStringFormat(rerr.FormatOptions{
		InvertOutput:   true, // flag that inverts the error output (wrap errors shown first)
		WithTrace:      true, // flag that enables stack trace output
		InvertTrace:    true, // flag that inverts the stack trace output (top of call stack shown first)
		WithExternal:   true,
		HideWrapFrames: true,
	})
	errStr := rerr.ToCustomString(err, format)
	r.log.Error(errStr)
	r.Gin().JSON(http.StatusInternalServerError, err.Error())
}

// GetUrlID 获取url中的id /user/:id
func (r *Rest) GetUrlID() (common.ID, error) {
	var id common.ID
	idstr := r.gin.Param("id")
	if idstr == "" {
		return id, rerr.New("url 缺少参数id")
	}
	idInt, err := strconv.Atoi(idstr)
	if err != nil {
		return id, rerr.WrapS(err, "url 参数id类型错误")
	}
	return common.ID(idInt), nil
}

// === 包装方法
func (r *Rest) Gin() *gin.Context {
	return r.gin
}

func (r *Rest) ShouldBind(o any) error {
	err := r.gin.ShouldBind(o)
	return wrapCaller(err, 1)
}

func (r *Rest) ShouldBindQuery(o any) error {
	err := r.gin.ShouldBindQuery(o)
	return wrapCaller(err, 1)
}

func (r *Rest) JSON(code int, obj any) {
	r.gin.JSON(code, obj)
}

func (r *Rest) String(code int, format string, values ...any) {
	r.gin.String(code, format, values...)
}

func (r *Rest) Header(key, value string) {
	r.gin.Header(key, value)
}

func (r *Rest) Abort() {
	r.gin.Abort()
}
