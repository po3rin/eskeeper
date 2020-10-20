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

	// options
	user string
	pass string
}

// NewOption is optional func for eskeeper.New
type NewOption func(*Eskeeper)

// UserName is optional func for Elasticsearch user name.
func UserName(name string) NewOption {
	return func(e *Eskeeper) {
		e.user = name
	}
}

// Pass is optional func for Elasticsearch password.
func Pass(pass string) NewOption {
	return func(e *Eskeeper) {
		e.pass = pass
	}
}

// New inits Eskeeper.
func New(urls []string, opts ...NewOption) (*Eskeeper, error) {
	eskeeper := &Eskeeper{}

	for _, opt := range opts {
		opt(eskeeper)
	}

	es, err := newEsClient(urls, eskeeper.user, eskeeper.pass)
	if err != nil {
		return nil, err
	}

	eskeeper.client = es
	return eskeeper, nil
}

// Sync synchronizes config & Elasticsearch State.
func (e *Eskeeper) Sync(ctx context.Context, reader io.Reader) error {
	conf, err := yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}

	err = e.client.syncIndex(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync indices")
	}

	err = e.client.syncAlias(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync aliases")
	}
	return nil
}
