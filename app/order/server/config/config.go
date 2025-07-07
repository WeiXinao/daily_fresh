package config

import (
	"encoding/json"

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
	Telemtry *options.TelemtryOptions  `json:"telemetry" mapstructure:"telemetry"`
	MySQL    *options.MysqlOptions     `json:"mysql" mapstructure:"mysql"`
	Redis    *options.RedisOptions     `json:"redis" mapstructure:"redis"`
	Dtm *options.DtmOptions `json:"dtm" mapstructure:"dtm"`
}

func New() *Config {
	return &Config{
		Log:      log.NewOptions(),
		Server:   options.NewServerOptions(),
		Registry: options.NewRegistryOptions(),
		Telemtry: options.NewTelemtryOptions(),
		MySQL:    options.NewMySQLOptions(),
		Redis: options.NewRedisOptions(),
		Dtm: options.NewDtmOptions(),
	}
}

// Flags implements app.CliOptions.
func (c *Config) Flags() (fss flag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Telemtry.AddFlags(fss.FlagSet("telemetry"))
	c.MySQL.AddFlags(fss.FlagSet("mysql"))
	c.Redis.AddFlags(fss.FlagSet("redis"))
	c.Dtm.AddFlags(fss.FlagSet("dtm"))
	return fss
}

// Validate implements app.CliOptions.
func (c *Config) Validate() []error {
	errs := make([]error, 0)
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	errs = append(errs, c.Telemtry.Validate()...)
	errs = append(errs, c.MySQL.Validate()...)
	errs = append(errs, c.Redis.Validate()...)
	return errs
}

func (c *Config) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}
