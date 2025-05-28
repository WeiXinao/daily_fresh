package options

import (
	"errors"

	"github.com/spf13/pflag"
)


type RegisteryOptions struct {
	Address string `mapstructure:"address" json:"address"`
	Scheme string `mapstructure:"scheme" json:"scheme"`
}

func NewRegistryOptions() *RegisteryOptions {
	return &RegisteryOptions{
		Address: "127.0.0.1:8500",
		Scheme: "http",
	}
}

func (o *RegisteryOptions) Validate() []error {
	errs := []error{}
	if o.Address == "" || o.Scheme == "" {
		errs = append(errs, errors.New("address an scheme is empty"))
	}
	return errs
}

func (o *RegisteryOptions) AddFlags(fs pflag.FlagSet) {
	fs.StringVar(&o.Address, "consul.address", o.Address, "consul address, if left, default is 127.0.0.1:8500")

	fs.StringVar(&o.Scheme, "consul.scheme", o.Scheme, "consul scheme, if left, default is http")
} 