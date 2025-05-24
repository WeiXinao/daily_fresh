package code

import (
	"net/http"
	"slices"

	"github.com/WeiXinao/daily_your_go/pkg/errors"
)

var _ errors.Coder = (*ErrCode)(nil)

type ErrCode struct {
	// 错误码
	C int

	// http 的转态码
	HTTP int

	// 扩展字段
	Ext string

	// 引用文档
	Ref string
}

// Code implements errors.Coder.
func (e *ErrCode) Code() int {
	return e.C
}

// HTTPStatus implements errors.Coder.
func (e *ErrCode) HTTPStatus() int {
	if e.HTTP == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTP
}

// Reference implements errors.Coder.
func (e *ErrCode) Reference() string {
	return e.Ref
}

// String implements errors.Coder.
func (e *ErrCode) String() string {
	return e.Ext
}

func register(code int, httpStatus int, message string, refs ...string) {
	if !slices.Contains([]int{200, 400, 401, 403, 404, 500}, httpStatus) {
		panic("http code not in `200，400, 401, 403, 404, 500`")
	}
	
	var ref string
	if len(refs) <= 0 {
		ref = "reference not provided"
	}

	ref = refs[0]

	coder := ErrCode {
		C: code,
		HTTP: httpStatus,
		Ext: message,
		Ref: ref,
	}

	errors.MustRegister(&coder)
}