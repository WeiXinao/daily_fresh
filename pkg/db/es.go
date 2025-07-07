package db

import (
	"fmt"

	"github.com/WeiXinao/daily_fresh/pkg/log"
	"github.com/olivere/elastic/v7"
)

type EsOptions struct {
	Host string
	Port string
}

var _ elastic.Logger = esLog(func(format string, v ...interface{}) {})

type esLog func(format string, v ...interface{})

// Printf implements elastic.Logger.
func (e esLog) Printf(format string, v ...interface{}) {
	e(format, v...)
}

func NewEsClient(opts *EsOptions) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetErrorLog(esLog(func(format string, v ...interface{}) {
			log.Errorf("ELASTIC " + format, v...)
		})),
		elastic.SetInfoLog(esLog(func(format string, v ...interface{}) {
			log.Infof("ELASTIC " + format, v...)
		})),
		elastic.SetTraceLog(esLog(func(format string, v ...interface{}) {
			log.Debugf("ELASTIC " + format, v...)
		})),
		elastic.SetURL(fmt.Sprintf("http://%s:%s", opts.Host, opts.Port)),
		elastic.SetSniff(false),
	)
}
