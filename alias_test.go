package eskeeper

import (
	"context"
	"testing"
	"time"
)

func TestExistAlias(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		setup   func(tb testing.TB)
		want    bool
		cleanup func(tb testing.TB)
	}{
		{
			name:  "simple",
			alias: "alias-exist-v1",
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "alias-exist-index-v1")
				createTmpIndexHelper(tb, "alias-exist-index-v2")
				time.Sleep(10 * time.Second)
				createTmpAliasHelper(tb, "alias-exist-v1", "alias-exist-index-v1")
				createTmpAliasHelper(tb, "alias-exist-v1", "alias-exist-index-v2")
			},
			want: true,
		},
		{
			name:  "not-found",
			alias: "not-found",
			setup: func(tb testing.TB) {
			},
			want: false,
		},
	}

	es, err := newEsClient([]string{url}, "", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.setup(t)
			ok, err := es.existAlias(ctx, tt.alias)
			if err != nil {
				t.Error(err)
			}
			if ok != tt.want {
				t.Errorf("want: %+v, got: %+v\n", tt.want, ok)
			}
		})
	}
}

func TestSyncAliases(t *testing.T) {
	tests := []struct {
		name    string
		conf    config
		setup   func(tb testing.TB)
		cleanup func(tb testing.TB)
	}{
		{
			name: "simple",
			conf: config{
				Aliases: []alias{
					{
						Name:    "test-sync-alias",
						Indices: []string{"test-v1", "test-v2"},
					},
				},
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "test-v1")
				createTmpIndexHelper(tb, "test-v2")
			},
		},
		{
			name: "switch",
			conf: config{
				Aliases: []alias{
					{
						Name:    "test-sync-alias",
						Indices: []string{"test-v3"},
					},
				},
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "test-v3")
			},
		},
	}

	es, err := newEsClient([]string{url}, "", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.setup(t)
			err := es.syncAliases(ctx, tt.conf)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
