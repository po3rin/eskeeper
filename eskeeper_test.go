package eskeeper_test

import (
	"os"
	"reflect"

	"testing"

	"github.com/po3rin/eskeeper"
)

func TestYaml2Conf(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want eskeeper.Config
	}{
		{
			name: "simple",
			yaml: "testdata/es.yaml",
			want: eskeeper.Config{
				Index: []eskeeper.Index{
					eskeeper.Index{
						Name:    "test-v1",
						Mapping: "test.json",
					},
					eskeeper.Index{
						Name:    "test-v2",
						Mapping: "test.json",
					},
				},
				Alias: []eskeeper.Alias{
					eskeeper.Alias{
						Name:  "alias1",
						Index: []string{"test-v1"},
					},
					eskeeper.Alias{
						Name:  "alias2",
						Index: []string{"test-v1", "test-v2"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := os.Open(tt.yaml)
			if err != nil {
				t.Fatal(err)
			}
			got, err := eskeeper.Yaml2Conf(r)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("\nwant: %+v\ngot : %+v\n", tt.want, got)
			}
		})
	}
}
