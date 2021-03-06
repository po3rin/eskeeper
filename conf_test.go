package eskeeper

import (
	"os"
	"reflect"

	"testing"
)

func TestYaml2Conf(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want config
	}{
		{
			name: "simple",
			yaml: "testdata/es.yaml",
			want: config{
				Indices: []index{
					{
						Name:    "test-v1",
						Mapping: "testdata/test.json",
					},
					{
						Name:    "test-v2",
						Mapping: "testdata/test.json",
					},
					{
						Name:    "close-v1",
						Mapping: "testdata/test.json",
						Status:  "close",
					},
				},
				Aliases: []alias{
					{
						Name:    "alias1",
						Indices: []string{"test-v1"},
					},
					{
						Name:    "alias2",
						Indices: []string{"test-v1", "test-v2"},
					},
				},
			},
		},
		{
			name: "reindex",
			yaml: "testdata/es.reindex.yaml",
			want: config{
				Indices: []index{
					{
						Name:    "test-v1",
						Mapping: "testdata/test.json",
					},
					{
						Name:    "reindex-v1",
						Mapping: "testdata/test.json",
						Reindex: reindex{
							Source:            "test-v1",
							Slices:            3,
							WaitForCompletion: true,
							On:                "firstCreated",
						},
					},
				},
				Aliases: []alias{
					{
						Name:    "alias1",
						Indices: []string{"test-v2"},
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
			got, err := yaml2Conf(r)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nwant: %+v\ngot : %+v\n", tt.want, got)
			}
		})
	}
}
