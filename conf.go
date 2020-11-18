package eskeeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

var status = map[string]struct{}{
	"open":  struct{}{},
	"close": struct{}{},
	"":      struct{}{},
}

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

func (e *Eskeeper) validateConfigFormat(c config) error {
	createIndices := make(map[string]struct{}, 0)

	for _, index := range c.Indices {
		_, exist := createIndices[index.Name]
		if exist {
			e.logf("[fail] index: %v\n", index.Name)
			return fmt.Errorf("duplicated index name %v", index.Name)
		}

		createIndices[index.Name] = struct{}{}

		if index.Name == "" {
			e.logf("[fail] index: %v\n", index.Name)
			return errors.New("index name is empty")
		}
		if index.Mapping != "" {
			m, err := ioutil.ReadFile(index.Mapping)
			if err != nil {
				e.logf("[fail] index: %v\n", index.Name)
				return fmt.Errorf("read file %v: %w", index.Mapping, err)
			}
			// validate json format
			var jsonStr map[string]interface{}
			if err := json.Unmarshal(m, &jsonStr); err != nil {
				e.logf("[fail] index: %v\n", index.Name)
				return fmt.Errorf("mapping json is invalid: %w", err)
			}
		}
		_, ok := status[index.Status]
		if !ok {
			e.logf("[fail] index: %v\n", index.Name)
			return fmt.Errorf("unsupported status %v", index.Status)
		}
		e.logf("[pass] index: %v\n", index.Name)
	}

	for _, alias := range c.Aliases {
		if alias.Name == "" {
			e.logf("[fail] alias: %v\n", alias.Name)
			return errors.New("alias name is empty")
		}

		if len(alias.Indices) == 0 {
			e.logf("[fail] alias: %v\n", alias.Name)
			return fmt.Errorf("no indices in %v alias", alias.Name)
		}

		_, ok := createIndices[alias.Name]
		if ok {
			e.logf("[fail] alias: %v\n", alias.Name)
			return fmt.Errorf("alias name %v is a duplicate of an index name that already exists", alias.Name)
		}

		for _, index := range alias.Indices {
			if index == "" {
				e.logf("[fail] alias: %v\n", alias.Name)
				return errors.New("index name is empty")
			}
		}
		e.logf("[pass] alias: %v\n", alias.Name)
	}
	return nil
}
