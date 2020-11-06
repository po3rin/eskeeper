package eskeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

func (c *esclient) syncAliases(ctx context.Context, conf config) error {
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
