package eskeeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type config struct {
	Indices []index `json:"index"`
	Aliases []alias `json:"alias"` // supports close only
}

type index struct {
	Name    string `json:"name"`
	Mapping string `json:"mapping"`
	Status  string `json:"status"`
}

type alias struct {
	Name    string   `json:"name"`
	Indices []string `json:"index"`
}

func yaml2Conf(reader io.Reader) (config, error) {
	conf := config{}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(b, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}

func validateConfigFormat(c config) error {
	createIndices := make(map[string]struct{}, 0)

	for _, index := range c.Indices {
		_, exist := createIndices[index.Name]
		if exist {
			return fmt.Errorf("duplicated index name %v", index.Name)
		}

		createIndices[index.Name] = struct{}{}

		if index.Name == "" {
			return errors.New("index name is empty")
		}
		if index.Mapping != "" {
			m, err := ioutil.ReadFile(index.Mapping)
			if err != nil {
				return fmt.Errorf("read file %v: %w", index.Mapping, err)
			}
			// validate json format
			var jsonStr map[string]interface{}
			if err := json.Unmarshal(m, &jsonStr); err != nil {
				return fmt.Errorf("mapping json is invalid: %w", err)
			}
		}
	}

	for _, alias := range c.Aliases {
		if alias.Name == "" {
			return errors.New("alias name is empty")
		}
		if len(alias.Indices) == 0 {
			return fmt.Errorf("no indices in %v alias", alias.Name)
		}
		for _, index := range alias.Indices {
			if index == "" {
				return errors.New("index name is empty")
			}
		}
	}
	return nil
}
