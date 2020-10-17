package eskeeper

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type Config struct {
	Index []Index `json:"index"`
	Alias []Alias `json:"alias"`
}

type Index struct {
	Name    string `json:"name"`
	Mapping string `json:"mapping"`
}

type Alias struct {
	Name  string   `json:"name"`
	Index []string `json:"index"`
}

func Yaml2Conf(reader io.Reader) (*Config, error) {
	conf := &Config{}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func Sync(ctx context.Context, reader io.Reader) error {
	conf, err := Yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}
	_ = conf
	return nil
}
