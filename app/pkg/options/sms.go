package options

import "github.com/spf13/pflag"

type SmsOptions struct {
	APIKey string `mapstructure:"key" json:"key"`
	APISecret string `mapstructure:"secret" json:"secret"`
}

func NewSmsOptions() *SmsOptions {
	return &SmsOptions{
		APIKey: "",
		APISecret: "",
	}
}

func (s *SmsOptions) Validate() []error {
	errs := []error{}
	return errs
}

func (s *SmsOptions) AddFlags(fs *pflag.FlagSet)  {
	fs.StringVar(&s.APIKey, "sms.apikey", s.APIKey, "sms api key")	
	fs.StringVar(&s.APISecret, "sms.secret", s.APISecret, "sms api secret")
}