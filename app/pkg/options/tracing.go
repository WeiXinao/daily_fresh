package options

import (
	"slices"

	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/spf13/pflag"
)

type TelemtryOptions struct {
	Name string `json:"name" mapstructure:"name"`
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
	Sampler float64 `json:"sampler" mapstructure:"sampler"`
	Batcher string `json:"batcher" mapstructure:"batcher"`
}

func NewTelemtryOptions() *TelemtryOptions {
	return &TelemtryOptions{
		Name: "daily_your_go",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler: 1.0,
		Batcher: "jaeger",
	}
}

func (t *TelemtryOptions) Validate() []error {
	errs := []error{}
	batchers := []string{"jaeger", "zipkin"}
	if !slices.Contains(batchers, t.Batcher) {
		errs = append(errs, errors.New("opentelemetry only support jaeger or zipkin"))
	}
	return errs
}

// AddFlags add flags reload to open telemetry for specific tracing to the specific flagset.
func (t *TelemtryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&t.Name, "telemetry.name", t.Name, "opentelemetry name")
	fs.StringVar(&t.Endpoint, "telemetry.endpoint", t.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&t.Sampler, "telemetry.sampler", t.Sampler, "opentelemetry sampler")
	fs.StringVar(&t.Batcher, "telemetry.batcher", t.Batcher, "opentelemetry batcher, only support jaeger or zipkin")
}