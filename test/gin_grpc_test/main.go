package main

import (
	"reflect"
	"strings"

	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/WeiXinao/daily_fresh/test/gin_grpc_test/proto"
	"github.com/WeiXinao/daily_fresh/test/gin_grpc_test/server"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	universal_translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func main() {
	e := gin.Default()
	svc := server.NewSayHelloServer()
	tanslator, err := initTrans("zh")
	if err != nil {
		panic(err)
	}

	proto.RegisterHelloServiceServerHTTPServer(tanslator, svc, e)
	e.Run(":8080")
}

func initTrans(locale string) (trans universal_translator.Translator, err error) {
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
		uni := universal_translator.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return nil, errors.Errorf("uni.GetTranslator(%s)", locale)
		}	

		switch locale {
			case "en":
				err = en_translations.RegisterDefaultTranslations(v, trans)
				if err != nil {
					return nil, err
				}
			case "zh":
				err = zh_translations.RegisterDefaultTranslations(v, trans)
				if err != nil {
					return nil, err
				}	
			default:
				err = en_translations.RegisterDefaultTranslations(v, trans)
				if err != nil {
					return nil, err
				}
		}
	}
	return trans, nil
}