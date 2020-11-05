package eskeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gofrs/uuid"
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

// TODO: alias pre-check
func (c *esclient) preCheck(ctx context.Context, conf config) error {
	for _, ix := range conf.Indices {
		// generate uuid for pre-check create index
		u2, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("generate UUID for pre-check: %w", err)
		}

		preIndex := index{
			Name:    fmt.Sprintf("eskeeper-%s", u2.String()),
			Mapping: ix.Mapping,
		}

		err = c.createIndex(ctx, preIndex)
		if err != nil {
			return fmt.Errorf("pre-check: pre create using random name index: %w", err)
		}

		err = c.deleteIndex(ctx, preIndex.Name)
		if err != nil {
			return fmt.Errorf("pre-check: delete pre-created index: %w", err)
		}
	}
	return nil
}

// postCheck checks created index & alias by name only.
func (c *esclient) postCheck(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		ok, err := c.existIndex(ctx, index.Name)
		if err != nil {
			return fmt.Errorf("post-check: check created index %v exist: %w", index.Name, err)
		}
		if !ok {
			return fmt.Errorf("post-check: created index %v is not found", index.Name)
		}
	}

	for _, alias := range conf.Aliases {
		ok, err := c.existAlias(ctx, alias.Name)
		if err != nil {
			return fmt.Errorf("post-check: check created alias %v exist: %w", alias.Name, err)
		}
		if !ok {
			return fmt.Errorf("post-check: created alias %v is not found", alias.Name)
		}
	}
	return nil
}

func (c *esclient) existIndex(ctx context.Context, index string) (bool, error) {
	exists := c.client.Indices.Exists

	res, err := exists(
		[]string{index},
		exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("check index exists: %w", err)
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func (c *esclient) existAlias(ctx context.Context, alias string) (bool, error) {
	exists := c.client.Indices.ExistsAlias

	res, err := exists(
		[]string{alias},
		exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("check alias exists: %w", err)
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func (c *esclient) createIndex(ctx context.Context, index index) error {
	create := c.client.Indices.Create
	ok, err := c.existIndex(ctx, index.Name)
	if err != nil {
		return err
	}
	// index already exists
	if ok {
		return nil
	}

	f, err := os.Open(index.Mapping)
	if err != nil {
		return fmt.Errorf("open mapping file: %w", err)
	}

	res, err := create(
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
			return fmt.Errorf("failed to create index [index=%v, statusCode=%v]", index.Name, res.StatusCode)
		}
		return fmt.Errorf("failed to create index [index=%v, statusCode=%v, res=%v]", index.Name, res.StatusCode, string(body))
	}
	return nil
}

func (c *esclient) deleteIndex(ctx context.Context, index string) error {
	delete := c.client.Indices.Delete
	res, err := delete([]string{index})
	if err != nil {
		return fmt.Errorf("delete index: %w", err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to delete index [index= %v, statusCode=%v]", index, res.StatusCode)
		}
		return fmt.Errorf("failed to delete index [index= %v, statusCode=%v, res=%v]", index, res.StatusCode, string(body))
	}
	return nil
}

func (c *esclient) syncIndex(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		err := c.createIndex(ctx, index)
		if err != nil {
			return fmt.Errorf("create index: %w", err)
		}
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
				return fmt.Errorf("failed to sync alias [alias=%v, statusCode=%v]", alias.Name, res.StatusCode)
			}
			return fmt.Errorf("failed to sync alias [alias=%v, statusCode=%v, res=%v]", alias.Name, res.StatusCode, string(body))
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
