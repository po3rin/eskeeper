package eskeeper

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
)

// TODO: alias pre-check
func (c *esclient) preCheck(ctx context.Context, conf config) error {
	// use alias pre-check
	createIndices := make(map[string]struct{}, 0)

	for _, ix := range conf.Indices {
		createIndices[ix.Name] = struct{}{}

		// generate uuid for pre-check create index
		u2, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("generate UUID for pre-check: %w", err)
		}

		preIndex := index{
			Name:    fmt.Sprintf("eskeeper-%s", u2.String()),
			Mapping: ix.Mapping,
		}

		err = c.syncIndex(ctx, preIndex)
		if err != nil {
			return fmt.Errorf("pre-check: pre create using random name index: %w", err)
		}

		err = c.deleteIndex(ctx, preIndex.Name)
		if err != nil {
			return fmt.Errorf("pre-check: delete pre-created index: %w", err)
		}
	}

	// check target index exists
	for _, alias := range conf.Aliases {
		for _, index := range alias.Indices {
			if _, ok := createIndices[index]; ok {
				continue
			}
			ok, err := c.existIndex(ctx, index)
			if err != nil {
				return fmt.Errorf("pre-check: check index %v exists for alias %v: %w", index, alias.Name, err)
			}
			if !ok {
				return fmt.Errorf("pre-check: index %v for alias %v is not found", index, alias.Name)
			}
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
