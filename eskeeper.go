package eskeeper

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

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

	err = validateConfigFormat(conf)
	if err != nil {
		return errors.Wrap(err, "validate config")
	}

	err = e.client.preCheck(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "pre-check")
	}

	err = e.client.syncIndices(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync indices")
	}

	err = e.client.syncAliases(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync aliases")
	}

	err = e.client.syncCloseStatus(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "sync aliases")
	}

	err = e.client.postCheck(ctx, conf)
	if err != nil {
		return errors.Wrap(err, "post-check")
	}
	return nil
}

// Validate validates cofig.
func (e *Eskeeper) Validate(ctx context.Context, reader io.Reader) error {
	conf, err := yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}
	err = validateConfigFormat(conf)
	if err != nil {
		return errors.Wrap(err, "validate config")
	}
	return nil
}
