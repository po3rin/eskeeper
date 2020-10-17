package eskeeper

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type config struct {
	Indices []index `json:"index"`
	Aliases []alias `json:"alias"`
}

type index struct {
	Name    string `json:"name"`
	Mapping string `json:"mapping"`
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

// Eskeeper manages indices & aliases.
type Eskeeper struct {
	client *esclient
}

// New inits Eskeeper.
func New(urls []string, user, pass string) (*Eskeeper, error) {
	es, err := newEsClient(urls, user, pass)
	if err != nil {
		return nil, err
	}
	return &Eskeeper{
		client: es,
	}, nil
}

// Sync synchronizes config & Elasticsearch State.
func (e *Eskeeper) Sync(ctx context.Context, reader io.Reader) error {
	conf, err := yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}
	err = e.client.createIndex(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync indices")
	}
	err = e.client.syncAlias(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync aliases")
	}
	return nil
}
