package eskeeper

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
)

func (c *esclient) reindex(ctx context.Context, dest string, reindex reindex) error {
	ri := c.client.Reindex
	body := strings.NewReader(
		fmt.Sprintf(`
{
  "source": {
    "index": "%s"
  },
  "dest": {
    "index": "%s"
  }
}`,
			reindex.Source, dest,
		),
	)

	slices := reindex.Slices
	if slices == 0 {
		slices = 1
	}

	res, err := ri(
		body,
		ri.WithContext(ctx),
		ri.WithSlices(slices),
		ri.WithWaitForCompletion(reindex.WaitForCompletion),
	)
	if err != nil {
		return fmt.Errorf("reindex: %w", err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to reindex [index=%v, statusCode=%v, res=%v]", reindex.Source, res.StatusCode, string(body))
	}
	return nil
}
