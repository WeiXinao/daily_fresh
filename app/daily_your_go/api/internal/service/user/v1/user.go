package user

import (
	"context"
	"time"

	"github.com/WeiXinao/daily_your_go/app/daily_your_go/api/internal/data"
	"github.com/WeiXinao/daily_your_go/app/pkg/code"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	baseCode "github.com/WeiXinao/daily_your_go/gmicro/code"
	gjwt "github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares/jwt"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/WeiXinao/daily_your_go/pkg/log"
	"github.com/WeiXinao/daily_your_go/pkg/storage"
	"github.com/dgrijalva/jwt-go"
)

type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, pwd string) (*UserDTO, error)
	Register(ctx context.Context, mobile, pwd, codes string) (*UserDTO, error)
	Update(ctx context.Context, userDTO *UserDTO) error
	Get(ctx context.Context, id uint64) (*UserDTO, error)
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
func (u *UserService) Get(ctx context.Context, id uint64) (*UserDTO, error) {
	user, err := u.ud.Get(ctx, id)
	if err != nil {
		return nil, err	
	}
	return &UserDTO{User: user}, nil
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
		return nil, errors.WithCode(baseCode.ErrPasswordIncorrect, "密码错误")
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
		return nil, errors.WithCode(baseCode.ErrGenerateTokenFailed, "生成token失败，err: %w", err)
	}

	return &UserDTO{
		User:      user,
		Token:     token,
		ExpiresAt: int64(expiresAt),
	}, nil
}

// Rigister implements UserSrv.
func (u *UserService) Register(ctx context.Context, mobile string, pwd, codes string) (*UserDTO, error) {
	rstore := storage.RedisCluster{}
	value, err := rstore.GetKey(ctx, mobile)
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不存在")
	}
	if value != codes {
		return nil, errors.WithCode(code.ErrCodeIncorrent, "验证码错误")
	}

	user := &data.User{
		Mobile:   mobile,
		Password: pwd,
	}
	err = u.ud.Create(ctx, user)
	if err != nil {
		log.Errorf("user register failed: %v", err)
		return nil, err
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
		return nil, errors.WithCode(baseCode.ErrGenerateTokenFailed, "生成token失败，err: %w", err)
	}
	return &UserDTO{
		User:      *user,
		Token:     token,
		ExpiresAt: int64(expiresAt),
	}, nil
}

// Update implements UserSrv.
func (u *UserService) Update(ctx context.Context, userDTO *UserDTO) error {
	return u.ud.Update(ctx, &data.User{
		ID: userDTO.ID,
		NickName: userDTO.NickName,
		Password: userDTO.Password,
		Birthday: userDTO.Birthday,
		Gender:   userDTO.Gender,
		Role:     userDTO.Role,
	})
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
