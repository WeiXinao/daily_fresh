package trace

const TraceName = "daily_your_go"

type Options struct {
	Name string `json:"name" mapstructure:"name"`
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
	Sampler float64 `json:"sampler" mapstructure:"sampler"`
	Batcher string `json:"batcher" mapstructure:"batcher"`
}