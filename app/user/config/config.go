package config

import (
	"github.com/WeiXinao/daily_your_go/pkg/app"
	"github.com/WeiXinao/daily_your_go/pkg/common/cli/flag"
	"github.com/WeiXinao/daily_your_go/pkg/log"
)

var _ app.CliOptions = (*Config)(nil)

type Config struct {
	Log *log.Options `json:"log" mapstructure:"log"`
}

func New() *Config {
	return &Config{
		Log: log.NewOptions(),
	}
}

// Flags implements app.CliOptions.
func (c *Config) Flags() (fss flag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	return fss
}

// Validate implements app.CliOptions.
func (c *Config) Validate() []error {
	errs := make([]error, 0)
	errs = append(errs, c.Log.Validate()...)
	return errs
}

