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

	// index dose not exist.
	if !ok {
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

		// reindex --------

		if index.Reindex.Source == "" {
			return nil
		}

		ok, err := c.existIndex(ctx, index.Reindex.Source)
		if err != nil {
			return fmt.Errorf("check index exists for reindex process: %w", err)
		}
		if !ok {
			return fmt.Errorf("reindex (%s -> %s) conf is invalid. Make sure %s index exists", index.Reindex.Source, index.Reindex.Source, index.Name)
		}

		err = c.reindex(ctx, index.Name, index.Reindex)
		if err != nil {
			return fmt.Errorf("reindex (%s -> %s)", index.Reindex.Source, index.Name)
		}
		return nil
	}

	// index already exists.

	// reindex -------
	if index.Reindex.Source != "" && index.Reindex.On == "always" {
		ok, err = c.existIndex(ctx, index.Reindex.Source)
		if err != nil {
			return fmt.Errorf("check index exists for reindex process: %w", err)
		}
		if !ok {
			return fmt.Errorf("reindex (%s -> %s) conf is invalid. Make sure %s index exists", index.Reindex.Source, index.Reindex.Source, index.Name)
		}
		err = c.reindex(ctx, index.Name, index.Reindex)
		if err != nil {
			return fmt.Errorf("reindex (%s -> %s)", index.Reindex.Source, index.Name)
		}
		return nil
	}

	// Since downtime may occur when switching aliases, only open is processed before switching aliases.
	// TODO: refactoring.
	if index.Status == "close" {
		return nil
	}
	err = c.openIndex(ctx, index)
	if err != nil {
		return fmt.Errorf("open index: %w", err)
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

func (c *esclient) syncCloseStatus(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		if index.Status == "close" {
			err := c.closeIndex(ctx, index)
			if err != nil {
				return fmt.Errorf("crearted index status action: %w", err)
			}
		}
	}
	return nil
}

// func (c *esclient) indexStatusAction(ctx context.Context, index index) error {
// 	switch index.Status {
// 	case "close":
// 		err := c.closeIndex(ctx, index)
// 		if err != nil {
// 			return fmt.Errorf("close index: %w", err)
// 		}
// 	case "open":
// 		err := c.openIndex(ctx, index)
// 		if err != nil {
// 			return fmt.Errorf("open index: %w", err)
// 		}
// 	default:
// 		err := c.openIndex(ctx, index)
// 		if err != nil {
// 			return fmt.Errorf("open index: %w", err)
// 		}
// 	}
// 	return nil
// }

func (c *esclient) closeIndex(ctx context.Context, index index) error {
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
	return nil
}

func (c *esclient) openIndex(ctx context.Context, index index) error {
	open := c.client.Indices.Open
	res, err := open([]string{index.Name}, open.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("open index: %w", err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to open index [index= %v, statusCode=%v]: %w", index, res.StatusCode, err)
		}
		return fmt.Errorf("failed to open index [index= %v, statusCode=%v, res=%v]: %w", index, res.StatusCode, string(body), err)
	}
	return nil
}

func (c *esclient) syncIndices(ctx context.Context, conf config) error {
	for _, index := range conf.Indices {
		err := c.syncIndex(ctx, index)
		if err != nil {
			c.logf("[fail] index: %v\n", index.Name)
			return fmt.Errorf("sync index: %w", err)
		}
		c.logf("[synced] index: %v\n", index.Name)
	}
	return nil
}
