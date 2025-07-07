type {{$.Name}}HttpServer struct{
	server {{ $.ServiceName }}
	router gin.IRouter
	translator universal_translator.Translator
}

func Register{{ $.ServiceName }}HTTPServer(trans universal_translator.Translator, srv {{ $.ServiceName }}, r gin.IRouter) {
	s := {{.Name}}HttpServer{
		translator: trans,
		server: srv,
		router: r,
	}
	s.RegisterService()
}

{{range .Methods}}
func (s *{{$.Name}}HttpServer) {{ .HandlerName }} (ctx *gin.Context) {
	var in {{.Request}}
{{if eq .Method "GET" "DELETE" }}
	if len(ctx.Request.URL.Query()) > 0 {
		if err := ctx.ShouldBindQuery(&in); err != nil {
			ginx.HandleValidatorError(ctx, s.translator, err)
			return
		}
	}
{{else if eq .Method "POST" "PUT" }}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ginx.HandleValidatorError(ctx, s.translator, err)
		return
	}
{{else}}
	if err := ctx.ShouldBind(&in); err != nil {
		ginx.HandleValidatorError(ctx, s.translator, err)
		return
	}
{{end}}
{{if .HasPathParams }}
	if err := ctx.ShouldBindUri(&in); err != nil {
		ginx.HandleValidatorError(ctx, s.translator, err)
		return
	}
{{end}}
	out, err := s.server.{{.Name}}(ctx, &in)
	if err != nil {
		log.Error(err.Error())
		core.WriteResponse(ctx, errors.WithCode(code.ErrUnknown, err.Error()), nil)
		return
	}

	core.WriteResponse(ctx, nil, out)
}
{{end}}

func (s *{{$.Name}}HttpServer) RegisterService() {
{{range .Methods -}}
		s.router.Handle("{{.Method}}", "{{.Path}}", s.{{ .HandlerName }})
{{end -}}
}