package eskeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Cside/jsondiff"
)

type indexConfig struct {
	Settings map[string]interface{} `json:"settings"`
	Mappings map[string]interface{} `json:"mappings"`
}

func (c *esclient) index(ctx context.Context, index index) ([]byte, error) {
	get := c.client.Indices.Get
	res, err := get(
		[]string{index.Name},
		get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
		}
		return nil, fmt.Errorf("get %v settings: %v", index.Name, string(body))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
	}
	return b, nil
}

func (c *esclient) settings(ctx context.Context, index index) ([]byte, error) {
	getSettings := c.client.Indices.GetSettings
	res, err := getSettings(
		getSettings.WithIndex(index.Name),
		getSettings.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
		}
		return nil, fmt.Errorf("get %v settings: %v", index.Name, string(body))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
	}
	return b, nil
}

func (c *esclient) mapping(ctx context.Context, index index) ([]byte, error) {
	getMapping := c.client.Indices.GetMapping
	res, err := getMapping(
		getMapping.WithIndex(index.Name),
		getMapping.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("get %v mapping: %w", index.Name, err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("get %v mapping: %w", index.Name, err)
		}
		return nil, fmt.Errorf("get %v mapping: %v", index.Name, string(body))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("get %v settings: %w", index.Name, err)
	}
	return b, nil
}

func (c *esclient) equalSettings(ctx context.Context, index index, settingJSON []byte) (bool, error) {
	m, err := c.settings(ctx, index)
	if err != nil {
		return false, fmt.Errorf("get setting: %w", err)
	}

	if diff := jsondiff.Diff(m, settingJSON); diff != "" {
		// fmt.Printf("detect settings diff:\n%s", diff)
		return false, nil
	}
	return true, nil
}

func (c *esclient) equalMapping(ctx context.Context, index index, mappingsJSON []byte) (bool, error) {
	m, err := c.mapping(ctx, index)
	if err != nil {
		return false, fmt.Errorf("get mappings: %w", err)
	}

	if diff := jsondiff.Diff(m, mappingsJSON); diff != "" {
		// fmt.Printf("detect mapping diff:\n%s", diff)
		return false, nil

	}
	return true, nil
}

func (c *esclient) updateIndex(ctx context.Context, index index) error {
	putMapping := c.client.Indices.PutMapping
	putSetting := c.client.Indices.PutSettings

	b, err := ioutil.ReadFile(index.Mapping)
	if err != nil {
		return fmt.Errorf("open mapping file: %w", err)
	}

	config := &indexConfig{}
	err = json.Unmarshal(b, config)
	if err != nil {
		return fmt.Errorf("unmarshal mapping json: %w", err)
	}

	// mapping --------
	if config.Mappings == nil {
		// dose not contain
		// return fmt.Errorf("get mappings from file: %w", err)
	} else {
		j, err := json.Marshal(config.Mappings)
		if err != nil {
			return fmt.Errorf("marshal mappings json: %w", err)
		}

		ok, err := c.equalMapping(ctx, index, j)
		if err != nil {
			return err
		}

		if !ok {
			res, err := putMapping(
				bytes.NewReader(j),
				putMapping.WithIndex(index.Name),
				putMapping.WithContext(ctx),
			)
			if err != nil {
				return fmt.Errorf("get %v mapping: %w", index, err)
			}
			if res.StatusCode != 200 {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return fmt.Errorf("update %v mapping: %w", index.Name, err)
				}
				return fmt.Errorf("update %v mapping: %v", index.Name, string(body))
			}
		}
	}

	// setting -------
	if config.Settings == nil {
		// dose not contain
		// return fmt.Errorf("get settings from file: %w", err)
		return nil
	}

	j, err := json.Marshal(config.Settings)
	if err != nil {
		return fmt.Errorf("marshal mappings json: %w", err)
	}

	ok, err := c.equalSettings(ctx, index, j)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	res, err := putSetting(
		bytes.NewReader(j),
		putSetting.WithIndex(index.Name),
		putSetting.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("update %v setting: %w", index, err)
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("update %v setting: %w", index.Name, err)
		}
		return fmt.Errorf("update %v setting: %v", index.Name, string(body))
	}

	return nil
}
