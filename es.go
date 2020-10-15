package eskeeper

import (
	"context"

	"github.com/olivere/elastic"
)

type esclient struct {
	client *elastic.Client
}

func newEsClient(url, user, pass string) (*esclient, error) {
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		// Handle error
		panic(err)
	}
	return &esclient{
		client: client,
	}, nil

}

func (c *esclient) createIndex(ctx context.Context, conf Config) error {
	return nil
}

func (c *esclient) syncAlias(ctx context.Context, conf Config) error {
	return nil
}
