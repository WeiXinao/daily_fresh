package options

import "github.com/spf13/pflag"

type DtmOptions struct {
	GrpcServer string `mapstructure:"grpc" json:"grpc"`
	HttpServer string `mapstructure:"http" json:"http"`
}

func NewDtmOptions() *DtmOptions {
	return &DtmOptions{
		HttpServer: "http://localhost:36789/api/dtmsvr",
		GrpcServer: "192.168.5.52:36790",
	}
}

func (o *DtmOptions) Validate() []error {
	errs := make([]error, 0)
	return errs
}

func (o *DtmOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.GrpcServer, "dtm.grpc", o.GrpcServer, "dtm grpc server url")
	fs.StringVar(&o.HttpServer, "dtm.http", o.HttpServer, "dtm http server url")
}
