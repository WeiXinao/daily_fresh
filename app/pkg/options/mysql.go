package options

import (
	"time"

	"github.com/spf13/pflag"
)

type MysqlOptions struct {
	Host string `json:"host" mapstructure:"host"`
	Port string `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`	
	Database string `json:"database" mapstructure:"database"`
	MaxIdleConnections int `json:"max-idle-connections" mapstructure:"max-idle-connections"`
	MaxOpenConnections int `json:"max-open-connections" mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time" mapstructure:"max-connection-life-time"`
	LogLevel int	`json:"log-level" mapstructure:"log-level"`
}

// NewMySQLOptions create a zero value instance.
func NewMySQLOptions() *MysqlOptions {
	return &MysqlOptions{
		Host: "127.0.0.1",
		Port: "3306",
		Username: "",
		Password: "",
		Database: "",
		MaxIdleConnections: 100,
		MaxOpenConnections: 100,
		MaxConnectionLifeTime: 10 * time.Second,
		LogLevel: 1, // Silent
	}
}

// Validate verifies flags passed to MySQLOptions.
func (mo *MysqlOptions) Validate() []error {
	errs := []error{}
	return errs
}

// AddFlags adds flags related to MySQLOptions for a specific api server to the specified FlagSet.
func (mo *MysqlOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&mo.Host, "mysql.host", mo.Host, "MySQL service host address. If left blank, the following related mysql options will be ignored.")
	fs.StringVar(&mo.Port, "mysql.port", mo.Port, "MySQL service host port.")
	fs.StringVar(&mo.Username, "mysql.username", mo.Username, "Username for access to MySQL service.")
	fs.StringVar(&mo.Password, "mysql.password", mo.Password, "Password for access to MySQL service.")
	fs.StringVar(&mo.Database, "mysql.database", mo.Database, "Database for access to MySQL service.")
	fs.IntVar(&mo.MaxIdleConnections, "mysql.max-idle-connections", mo.MaxIdleConnections, "Maximum idle connections to allocate.")
	fs.IntVar(&mo.MaxOpenConnections, "mysql.max-open-connections", mo.MaxOpenConnections, "Maximum open connections to allocate.")
	fs.DurationVar(&mo.MaxConnectionLifeTime, "mysql.max-connection-life-time", mo.MaxConnectionLifeTime, "Maximum connection life time.")
	fs.IntVar(&mo.LogLevel, "mysql.log-level", mo.LogLevel, "MySQL log level.")
}