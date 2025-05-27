package options

import "github.com/spf13/pflag"

type ServerOptions struct {
	// 是否开启 pprof
	EnableProfiling   bool   `json:"profiling" mapstructure:"profiling"`

	// 是否开启 metrics 
	EnableMetrics     bool   `json:"enable-metrics" mapstructure:"enable-metrics"`

	// 是否开启 health check
	EnableHealthCheck bool   `json:"enable-health-check" mapstructure:"enable-health-check"`

	// host
	Host              string    `json:"host" mapstructure:"host"`

	// port
	Port              int    `json:"port" mapstructure:"port"`

	// http port
	HttpPort          int    `json:"http-port" mapstructure:"http-port"`

	// 名称
	Name              string `json:"name" mapstructure:"name"`

	// 中间件
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		EnableHealthCheck: true,
		EnableProfiling: true,
		EnableMetrics: true,
		Host: "127.0.0.1",
		Port: 8078,
		HttpPort: 8079,
		Name: "daily-your-go-user-srv",
	}
}

func (s *ServerOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (s *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.EnableProfiling, "server.enable-profiling", s.EnableProfiling, 
	"enable profiling, if true, will add <host>:<port>/debug/pprof, default is true")
	fs.BoolVar(&s.EnableMetrics, "server.enable-matrics", s.EnableMetrics, 
	"enable metrics, if true, will add <host>:<port>/metrics, default is true")
	fs.BoolVar(&s.EnableHealthCheck, "server.enable-health-check", s.EnableHealthCheck, 
	"enable health check, if true, will add health check route, default is true")
	fs.StringVar(&s.Host, "server.host", s.Host, "server host, default is 127.0.0.1")
	fs.IntVar(&s.Port, "server.port", s.Port, "server port, default is 8078")
	fs.IntVar(&s.HttpPort, "server.http-port", s.HttpPort, "server http port, default is 8079")
	fs.StringVar(&s.Name, "server.name", s.Name, "server host, default is daily-your-go-user-srv")
}