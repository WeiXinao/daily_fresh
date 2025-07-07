package main

import (
	"github.com/WeiXinao/daily_fresh/tools/autoGenerateGinByProto/core/generator"
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {

	// 可能其他导入文件中 在init方法中，注入了 flag.StringVar()函数
	flag.Parse()
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
		}
		return nil
	})
}
