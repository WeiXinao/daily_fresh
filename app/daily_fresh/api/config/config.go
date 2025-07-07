package config

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/pkg/app"
	"github.com/WeiXinao/daily_your_go/pkg/common/cli/flag"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

var _ app.CliOptions = (*Config)(nil)

type Config struct {
	Log      *log.Options              `json:"log" mapstructure:"log"`
	Server   *options.ServerOptions    `json:"server" mapstructure:"server"`
	Registry *options.RegisteryOptions `json:"registry" mapstructure:"registry"`
	Jwt      *options.JwtOptions       `json:"jwt" mapstructure:"jwt"`
	Sms      *options.SmsOptions       `json:"sms" mapstructure:"sms"`
	Redis    *options.RedisOptions     `json:"redis" mapstructure:"redis"`
}

func New() *Config {
	return &Config{
		Log:      log.NewOptions(),
		Server:   options.NewServerOptions(),
		Registry: options.NewRegistryOptions(),
		Jwt:      options.NewJwtOptions(),
		Sms:      options.NewSmsOptions(),
		Redis:    options.NewRedisOptions(),
	}
}

// Flags implements app.CliOptions.
func (c *Config) Flags() (fss flag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Jwt.AddFlags(fss.FlagSet("jwt"))
	c.Sms.AddFlags(fss.FlagSet("sms"))
	c.Redis.AddFlags(fss.FlagSet("redis"))
	return fss
}

// Validate implements app.CliOptions.
func (c *Config) Validate() []error {
	errs := make([]error, 0)
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	errs = append(errs, c.Jwt.Validate()...)
	errs = append(errs, c.Sms.Validate()...)
	errs = append(errs, c.Jwt.Validate()...)
	return errs
}
