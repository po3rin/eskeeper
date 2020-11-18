package eskeeper

import (
	"context"
	"testing"
)

func TestPreCheck(t *testing.T) {
	tests := []struct {
		name    string
		conf    config
		wantErr bool
	}{
		{
			name: "simple",
			conf: config{
				Indices: []index{
					{
						Name:    "precheck1",
						Mapping: "testdata/test.json",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid-field",
			conf: config{
				Indices: []index{
					{
						Name:    "precheck2",
						Mapping: "testdata/invalid.json",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "not-found-index",
			conf: config{
				Indices: []index{
					{
						Name:    "precheck3",
						Mapping: "testdata/invalid.json",
					},
				},
				Aliases: []alias{
					{
						Name: "precheck-alias",
						Indices: []string{
							"precheck3",
							"precheck4", // not found
						},
					},
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "duplicated name",
		// 	conf: config{
		// 		Indices: []index{
		// 			{
		// 				Name:    "duplicated-name",
		// 				Mapping: "testdata/test.json",
		// 			},
		// 		},
		// 		Aliases: []alias{
		// 			{
		// 				Name: "duplicated-name",
		// 				Indices: []string{
		// 					"duplicated-name",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: true,
		// },
	}

	es, err := newEsClient([]string{url}, "", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := es.preCheck(ctx, tt.conf)
			if tt.wantErr && err == nil {
				t.Error("expect error")
			}
			if !tt.wantErr && err != nil {
				t.Error(err)
			}
		})
	}
}
