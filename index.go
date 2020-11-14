package eskeeper

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
)

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

func (c *esclient) syncIndex(ctx context.Context, index index) error {
	create := c.client.Indices.Create
	ok, err := c.existIndex(ctx, index.Name)
	if err != nil {
		return err
	}
	// index already exists
	if ok {
		err := c.indexStatusAction(ctx, index)
		if err != nil {
			return fmt.Errorf("index status action: %w", err)
		}
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
	err = c.indexStatusAction(ctx, index)
	if err != nil {
		return fmt.Errorf("crearted index status action: %w", err)
	}
	return nil
}

func (c *esclient) deleteIndex(ctx context.Context, index string) error {
	delete := c.client.Indices.Delete
	res, err := delete([]string{index}, delete.WithContext(ctx))
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

func (c *esclient) indexStatusAction(ctx context.Context, index index) error {
	switch index.Status {
	case "close":
		close := c.client.Indices.Close
		res, err := close([]string{index.Name}, close.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("close index: %w", err)
		}
		if res.StatusCode != 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("failed to close index [index= %v, statusCode=%v]", index, res.StatusCode)
			}
			return fmt.Errorf("failed to close index [index= %v, statusCode=%v, res=%v]", index, res.StatusCode, string(body))
		}
	default:
		open := c.client.Indices.Open
		res, err := open([]string{index.Name}, open.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("open index: %w", err)
		}
		if res.StatusCode != 200 {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("failed to open index [index= %v, statusCode=%v]", index, res.StatusCode)
			}
			return fmt.Errorf("failed to open index [index= %v, statusCode=%v, res=%v]", index, res.StatusCode, string(body))
		}
	}
	return nil
}

func (c *esclient) syncIndices(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		err := c.syncIndex(ctx, index)
		if err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}
	return nil
}
