package user

import "context"

type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, pwd string) (*UserDTO, error)
	Rigister(ctx context.Context, mobile, pwd string) (*UserDTO, error)
	Update(ctx context.Context, userDTO *UserDTO) error
	Get(ctx context.Context, id int64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	CheckPassword(ctx context.Context, pwd, encryptedPwd string) (bool, error)
}

type UserDTO struct {
	
}
