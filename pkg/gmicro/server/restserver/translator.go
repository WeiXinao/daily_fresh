package restserver

import (
	"reflect"
	"strings"

	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func (s *Server) initTrans(locale string) (err error) {
	// 修改 gin 框架中的引擎属性实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取 json 的 tag 的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		// 第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		s.trans, ok = uni.GetTranslator(locale)
		if !ok {
			return errors.Errorf("uni.GetTranslator(%s)", locale)
		}	

		switch locale {
			case "en":
				err = en_translations.RegisterDefaultTranslations(v, s.trans)
				if err != nil {
					return err
				}
			case "zh":
				err = zh_translations.RegisterDefaultTranslations(v, s.trans)
				if err != nil {
					return err
				}	
			default:
				err = en_translations.RegisterDefaultTranslations(v, s.trans)
				if err != nil {
					return err
				}
		}
	}
	return
}