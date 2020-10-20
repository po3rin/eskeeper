package eskeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)

type esclient struct {
	client *elasticsearch.Client
}

func newEsClient(urls []string, user, pass string) (*esclient, error) {
	conf := elasticsearch.Config{
		Addresses: urls,
		Username:  user,
		Password:  pass,
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &esclient{
		client: es,
	}, nil
}

func (c *esclient) createIndex(ctx context.Context, index index) error {
	create := c.client.Indices.Create
	exists := c.client.Indices.Exists

	res, err := exists(
		[]string{index.Name},
		exists.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}
	// alreadey exists
	if res.StatusCode == 200 {
		return nil
	}

	f, err := os.Open(index.Mapping)
	if err != nil {
		return fmt.Errorf("open mapping file: %w", err)
	}

	res, err = create(
		index.Name,
		create.WithContext(ctx),
		create.WithBody(f),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to create index [index= %v, statusCode=%v]", index.Name, res.StatusCode)
		}
		return fmt.Errorf("failed to create index [index= %v, statusCode=%v, res=%v]", index.Name, res.StatusCode, string(body))
	}
	return nil
}

// CLOSE_ONLY: close index
// DELETE: delete index
// DELETE_LIFECYCLE_POLICY add :wq
func (c *esclient) deleteIndex(ctx context.Context, index index) error {
	return nil
}

func (c *esclient) syncIndex(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		c.createIndex(ctx, index)
	}
	return nil
}

func (c *esclient) syncAlias(ctx context.Context, conf config) error {
	for _, alias := range conf.Aliases {
		i := c.client.Indices

		query := aliasQuery(alias.Name, alias.Indices)

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			return fmt.Errorf("build query: %w", err)
		}

		res, err := i.UpdateAliases(&buf, i.UpdateAliases.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("upsert aliases: %w", err)
		}
		if res.StatusCode != 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("failed to sync alias [alias= %v, statusCode=%v]", alias.Name, res.StatusCode)
			}
			return fmt.Errorf("failed to sync alias [alias= %v, statusCode=%v, res=%v]", alias.Name, res.StatusCode, string(body))
		}
	}
	return nil
}

func aliasQuery(aliasName string, indices []string) map[string]interface{} {
	Actions := make([]map[string]interface{}, 0, len(indices)+1)

	Actions = append(Actions, map[string]interface{}{
		"remove": map[string]interface{}{
			"index": "*",
			"alias": aliasName,
		},
	})

	for _, index := range indices {
		Actions = append(Actions, map[string]interface{}{
			"add": map[string]interface{}{
				"index": index,
				"alias": aliasName,
			},
		})
	}

	return map[string]interface{}{
		"actions": Actions,
	}
}
