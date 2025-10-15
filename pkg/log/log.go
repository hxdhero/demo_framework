package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"lls_api/pkg/config"
	"lls_api/pkg/util/notification"
)

type logger struct {
	zerolog zerolog.Logger
}

type LoggerContext struct {
	RequestID string `json:"request_id"`
	User      string `json:"user"`
	UserErr   string `json:"user_err"` // 解析用户失败时存储的错误信息
	skipFrame int
}

func DefaultContext() *LoggerContext {
	return &LoggerContext{
		RequestID: "-",
		User:      "-",
	}
}

var myLogger *logger

func InitLog() {
	if myLogger != nil {
		return
	}
	myLogger = &logger{zerolog: newZerolog()}
}

func (ctx *LoggerContext) addField(event *zerolog.Event) *zerolog.Event {
	return event.Str("request_id", "request_id:"+ctx.RequestID).Str("user", "user:"+ctx.User)
}

// InfoSkip 打印时跳过栈帧
func (ctx *LoggerContext) InfoSkip(skip int, str string) {
	ctx.addField(myLogger.zerolog.Info().CallerSkipFrame(skip)).Msg(str)
}

func (ctx *LoggerContext) Info(str string) {
	ctx.addField(myLogger.zerolog.Info()).Msg(str)
}

func (ctx *LoggerContext) InfofSkip(skip int, template string, args ...interface{}) {
	ctx.addField(myLogger.zerolog.Info().CallerSkipFrame(skip)).Msgf(template, args...)
}

func (ctx *LoggerContext) Infof(template string, args ...interface{}) {
	ctx.addField(myLogger.zerolog.Info()).Msgf(template, args...)
}

func (ctx *LoggerContext) Error(str string) {
	ctx.notifyFeishu(str)
	ctx.addField(myLogger.zerolog.Error()).Msg(str)
}

func (ctx *LoggerContext) ErrorErr(err error) {
	ctx.notifyFeishu(err.Error())
	ctx.addField(myLogger.zerolog.Error()).Msg(err.Error())
}

func (ctx *LoggerContext) Errorf(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	ctx.notifyFeishu(msg)
	ctx.addField(myLogger.zerolog.Error()).Msg(msg)
}

func (ctx *LoggerContext) Fatalf(template string, args ...interface{}) {
	ctx.addField(myLogger.zerolog.Fatal()).Msgf(template, args...)
}

func (ctx *LoggerContext) Fatal(str string) {
	ctx.addField(myLogger.zerolog.Fatal()).Msg(str)
}

func (ctx *LoggerContext) FatalErr(err error) {
	ctx.addField(myLogger.zerolog.Fatal()).Msg(err.Error())
}

func (ctx *LoggerContext) notifyFeishu(msg string) {
	if config.C.Env != "prod" {
		return
	}
	go func() {
		if err := notification.SendFeishuChatCard(config.C.Notification.Feishu.ErrChatID, msg); err != nil {
			ctx.Infof("HARAKIRI %s", err.Error())
		}
	}()
}
