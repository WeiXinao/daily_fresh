package es

import (
	"sync"

	search "github.com/WeiXinao/daily_your_go/app/goods/srv/internal/data_search/v1"
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/pkg/db"
	"github.com/WeiXinao/daily_your_go/pkg/errors"
	"github.com/olivere/elastic/v7"
)

var (
	searchFactory search.SearchFactory
	once     sync.Once
)

func GetSearchFactoryOr(opts *options.EsOptions) (search.SearchFactory, error) {
	if opts == nil && searchFactory == nil {
		return nil, errors.New("failed to get es client")
	}

	if searchFactory != nil {
		return searchFactory, nil
	}

	var (
		err error
		esClient *elastic.Client
	)
	once.Do(func() {
		opts := db.EsOptions{
			Host: opts.Host,
			Port: opts.Port,
		}
		esClient, err = db.NewEsClient(&opts)
		if err != nil {
			return
		}
		searchFactory = &dataSearch{
			esClient: esClient,
		}
	})
	if searchFactory == nil {
		return nil, errors.Wrap(err, "failed to get es client")
	}
	return searchFactory, nil
}

var _ search.SearchFactory = (*dataSearch)(nil)

type dataSearch struct {
	esClient *elastic.Client
}

// Goods implements search.SearchFactory.
func (d *dataSearch) Goods() search.GoodsStore {
	return newGoods(d)
}
