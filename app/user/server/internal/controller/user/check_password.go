package user

import (
	"context"
	"crypto/sha512"
	"strings"

	upb "github.com/WeiXinao/daily_your_go/api/user/v1"
	"github.com/anaskhan96/go-password-encoder"
)

// CheckPassword implements v1.UserServer.
func (u *userServer) CheckPassword(ctx context.Context, request *upb.PasswordCheckInfo) (*upb.CheckResponse, error) {
	// 校验密码
	passwordInfo := strings.Split(request.EncryptedPassword, "$")
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	check := password.Verify(request.Password, passwordInfo[2], passwordInfo[3], options)
	return &upb.CheckResponse{Success: check}, nil
}