package validate

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"lls_api/pkg/rerr"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// InitTrans 初始化翻译器
func InitTrans(locale string) (trans ut.Translator, err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		V = v
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(zhT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return nil, fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// 注册自定义校验器和翻译
		if err = RegisterCustomValidator(v, trans); err != nil {
			return nil, err
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// RegisterCustomValidator 注册自定义校验器
func RegisterCustomValidator(v *validator.Validate, trans ut.Translator) error {
	if err := v.RegisterValidation("alphaDashValidator", alphaDashValidator); err != nil {
		return rerr.Wrap(err)
	}
	if err := v.RegisterTranslation("alphaDashValidator", trans, func(ut ut.Translator) error {
		return ut.Add("alphaDashValidator", "{0}只能包含字母、数字和下划线", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	}); err != nil {
		return err
	}

	if err := v.RegisterValidation("customDateTime", customDateTimeValidator); err != nil {
		return rerr.Wrap(err)
	}
	if err := v.RegisterTranslation("customDateTime", trans, func(ut ut.Translator) error {
		return ut.Add("customDateTime", "{0}的格式是2006-01-02 15:04:05", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	}); err != nil {
		return err
	}

	if err := v.RegisterValidation("customDate", customDateValidator); err != nil {
		return rerr.Wrap(err)
	}
	if err := v.RegisterTranslation("customDate", trans, func(ut ut.Translator) error {
		return ut.Add("customDate", "{0}的格式是2006-01-02", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	}); err != nil {
		return err
	}

	return nil
}

// alphaDashValidator 数字字母下划线
func alphaDashValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", fl.Field().String())
	return match
}

// customDateTimeValidator 2006-01-02 15:04:05 格式的时间
func customDateTimeValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.ParseInLocation(time.DateTime, dateStr, time.Local)
	return err == nil
}

// customDateValidator  2006-01-02 格式的时间
func customDateValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.ParseInLocation(time.DateOnly, dateStr, time.Local)
	return err == nil
}
