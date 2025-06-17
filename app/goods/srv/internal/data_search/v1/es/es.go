package es

import (
	"sync"

	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/pkg/db"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/olivere/elastic/v7"
)

var (
	esClient *elastic.Client
	once sync.Once
)

func GetSearchFactoryOr(opts *options.EsOptions) (*elastic.Client, error) {
	if opts == nil && esClient == nil {
		return nil, errors.New("failed to get es client")
	}

	if esClient != nil {
		return esClient, nil
	}

	var err error
	once.Do(func() {
		opts := db.EsOptions{
			Host: opts.Host,
			Port: opts.Port,
		}
		esClient, err = db.NewEsClient(&opts)
	})
	if esClient == nil {
		return nil, errors.Wrap(err, "failed to get es client")
	}
	return esClient, nil
}