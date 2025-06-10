package user

import (
	"context"
	"time"

	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/code"
	gjwt "github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares/jwt"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/dgrijalva/jwt-go"
)

type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, pwd string) (*UserDTO, error)
	Rigister(ctx context.Context, mobile, pwd string) (*UserDTO, error)
	Update(ctx context.Context, userDTO *UserDTO) error
	Get(ctx context.Context, id int64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	CheckPassword(ctx context.Context, pwd, encryptedPwd string) (bool, error)
}

var _ UserSrv = (*UserService)(nil)

type UserService struct {
	ud data.UserData

	jwtOpts *options.JwtOptions
}

// CheckPassword implements UserSrv.
func (u *UserService) CheckPassword(ctx context.Context, pwd string, encryptedPwd string) (bool, error) {
	return u.ud.CheckPassword(ctx, pwd, encryptedPwd)
}

// Get implements UserSrv.
func (u *UserService) Get(ctx context.Context, id int64) (*UserDTO, error) {
	panic("unimplemented")
}

// GetByMobile implements UserSrv.
func (u *UserService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	panic("unimplemented")
}

// MobileLogin implements UserSrv.
func (u *UserService) MobileLogin(ctx context.Context, mobile string, pwd string) (*UserDTO, error) {
	user, err := u.ud.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	// 检查密码是否正确
	b, err := u.CheckPassword(ctx, pwd, user.Password)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, errors.WithCode(code.ErrPasswordIncorrect, "密码错误")
	}

	// 生成 token
	j := gjwt.NewJWT(u.jwtOpts.Key)
	expiresAt := time.Now().Local().Add(u.jwtOpts.Timeout).Unix()
	token, err := j.CreateToken(gjwt.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), // 签名的生效时间
			ExpiresAt: expiresAt,
			Issuer:    u.jwtOpts.Realm,
		},
	})
	if err != nil {
		return nil, errors.WithCode(code.ErrGenerateTokenFailed, "生成token失败，err: %w", err)
	}

	return &UserDTO{
		User:      user,
		Token:     token,
		ExpiresAt: int64(expiresAt),
	}, nil
}

// Rigister implements UserSrv.
func (u *UserService) Rigister(ctx context.Context, mobile string, pwd string) (*UserDTO, error) {
	panic("unimplemented")
}

// Update implements UserSrv.
func (u *UserService) Update(ctx context.Context, userDTO *UserDTO) error {
	panic("unimplemented")
}

type UserDTO struct {
	data.User

	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewUserService(ud data.UserData, jwtOpts *options.JwtOptions) *UserService {
	return &UserService{
		ud:      ud,
		jwtOpts: jwtOpts,
	}
}
