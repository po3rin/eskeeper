package eskeeper

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Eskeeper manages indices & aliases.
type Eskeeper struct {
	client *esclient

	// options
	user    string
	pass    string
	verbose bool
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

// Verbose is optional func for verbose option.
func Verbose(v bool) NewOption {
	return func(e *Eskeeper) {
		e.verbose = v
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

	es.verbose = eskeeper.verbose
	eskeeper.client = es

	return eskeeper, nil
}

// Sync synchronizes config & Elasticsearch State.
func (e *Eskeeper) Sync(ctx context.Context, reader io.Reader) error {
	e.log("loading config ...")
	conf, err := yaml2Conf(reader)
	if err != nil {
		return err
	}

	e.log("\n=== validation stage ===")
	err = e.validateConfigFormat(conf)
	if err != nil {
		return err
	}

	e.log("\n=== pre-check stage ===")
	err = e.client.preCheck(ctx, conf)
	if err != nil {
		return err
	}

	e.log("\n=== sync stage ===")
	err = e.client.syncIndices(ctx, conf)
	if err != nil {
		return err
	}

	err = e.client.syncAliases(ctx, conf)
	if err != nil {
		return err
	}

	err = e.client.syncCloseStatus(ctx, conf)
	if err != nil {
		return err
	}

	e.log("\n=== post-check stage ===")
	err = e.client.postCheck(ctx, conf)
	if err != nil {
		return err
	}

	e.log("\nsucceeded")
	return nil
}

// Validate validates cofig.
func (e *Eskeeper) Validate(ctx context.Context, reader io.Reader) error {
	conf, err := yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}
	err = e.validateConfigFormat(conf)
	if err != nil {
		return errors.Wrap(err, "validate config")
	}
	return nil
}

func (e *Eskeeper) log(msg string) {
	if e.verbose {
		fmt.Printf("\x1b[34m%s\x1b[0m\n", msg)
	}
}

func (e *Eskeeper) logf(format string, a ...interface{}) {
	if e.verbose {
		fmt.Printf(format, a...)
	}
}
