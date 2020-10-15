package eskeeper

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

type Config struct {
	Index map[string]string   `json:"index"`
	Alias map[string][]string `json:"alias"`
}

func Yaml2Conf(reader io.Reader) (Config, error) {
	var conf Config

	bBuf := new(bytes.Buffer)
	aBuf := io.TeeReader(reader, bBuf)

	indexPath, err := yaml.PathString("$.index")
	if err != nil {
		return conf, errors.Wrap(err, "parse path query")
	}

	var index map[string]string
	if err := indexPath.Read(aBuf, &index); err != nil {
		return conf, errors.Wrap(err, "read yaml")
	}
	conf.Index = index

	aliasPath, err := yaml.PathString("$.alias")
	if err != nil {
		return conf, errors.Wrap(err, "parse path query")
	}

	var alias map[string][]string
	if err := aliasPath.Read(bBuf, &alias); err != nil {
		return conf, errors.Wrap(err, "read yaml")
	}
	fmt.Println(alias)
	conf.Alias = alias

	return conf, nil
}

func Sync(ctx context.Context, reader io.Reader) error {
	conf, err := Yaml2Conf(reader)
	if err != nil {
		return errors.Wrap(err, "convert yaml to conf")
	}
}
