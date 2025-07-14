package config

import (
	"github.com/WeiXinao/daily_fresh/app/pkg/options"
	"github.com/WeiXinao/daily_fresh/pkg/app"
	"github.com/WeiXinao/daily_fresh/pkg/common/cli/flag"
	"github.com/WeiXinao/daily_fresh/pkg/log"
)

var _ app.CliOptions = (*Config)(nil)

type Config struct {
	Log      *log.Options              `json:"log" mapstructure:"log"`
	Server   *options.ServerOptions    `json:"server" mapstructure:"server"`
	Registry *options.RegisteryOptions `json:"registry" mapstructure:"registry"`
}

func New() *Config {
	return &Config{
		Log:      log.NewOptions(),
		Server:   options.NewServerOptions(),
		Registry: options.NewRegistryOptions(),
	}
}

// Flags implements app.CliOptions.
func (c *Config) Flags() (fss flag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	return fss
}

// Validate implements app.CliOptions.
func (c *Config) Validate() []error {
	errs := make([]error, 0)
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	return errs
}
