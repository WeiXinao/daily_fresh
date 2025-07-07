package options

import (
	"time"

	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/asaskevich/govalidator"
)

type JwtOptions struct {
	Realm string `json:"realm" mapstructure:"realm"`
	Key string `json:"key" mapstructure:"key"`
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm: "daily_your_go",
		Key: "daily_your_go",
		Timeout: time.Hour,
		MaxRefresh: time.Hour * 24 * 30,
	}
}

func (j *JwtOptions) Validate() []error {
	var errs []error
	if !govalidator.StringLength(j.Key, "6", "32") {
		errs = append(errs, errors.New("secret-key 长度必须在6到32之间"))
	}
	return errs
}

func (j *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&j.Realm, "jwt.realm", j.Realm, "Realm name to display to the user.")
	fs.StringVar(&j.Key, "jwt.key", j.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&j.Timeout, "jwt.timeout", j.Timeout, "JWT token timeout.")

	fs.DurationVar(&j.MaxRefresh, "jwt.max-refresh", j.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")	
}	