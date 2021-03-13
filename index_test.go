package eskeeper

import (
	"context"
	"testing"
)

func TestSyncIndices(t *testing.T) {
	tests := []struct {
		name  string
		conf  config
		setup func(tb testing.TB)
	}{
		{
			name: "simple",
			conf: config{
				Indices: []index{
					{
						Name:    "create1",
						Mapping: "testdata/test.json",
					},
				},
			},
		},
		{
			name: "multi",
			conf: config{
				Indices: []index{
					{
						Name:    "create2",
						Mapping: "testdata/test.json",
					},
					{
						Name:    "create3",
						Mapping: "testdata/test.json",
					},
				},
			},
		},
		{
			name: "idempotence",
			conf: config{
				Indices: []index{
					{
						Name:    "idempotence",
						Mapping: "testdata/test.json",
					},
					{
						Name:    "idempotence",
						Mapping: "testdata/test.json",
					},
				},
			},
		},
		{
			name: "close",
			conf: config{
				Indices: []index{
					{
						Name:    "create-with-close-v1",
						Mapping: "testdata/test.json",
						Status:  "close",
					},
				},
			},
		},
		{
			name: "exists close",
			conf: config{
				Indices: []index{
					{
						Name:    "create-with-close-v2",
						Mapping: "testdata/test.json",
						Status:  "close",
					},
				},
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "create-with-close-v2")
			},
		},
		{
			name: "open",
			conf: config{
				Indices: []index{
					{
						Name:    "open-v1",
						Mapping: "testdata/test.json",
					},
				},
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "open-v1")
				closeIndexHelper(tb, "open-v1")
			},
		},
		{
			name: "open-already-open",
			conf: config{
				Indices: []index{
					{
						Name:    "open-already-open-v1",
						Mapping: "testdata/test.json",
					},
				},
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
			if tt.setup != nil {
				tt.setup(t)
			}
			err := es.syncIndices(ctx, tt.conf)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestExistIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   string
		setup   func(tb testing.TB)
		want    bool
		cleanup func(tb testing.TB)
	}{
		{
			name:  "simple",
			index: "exist-v1",
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "exist-v1")
			},
			want: true,
		},
		{
			name:  "not-found",
			index: "exist-v2",
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
			ok, err := es.existIndex(ctx, tt.index)
			if err != nil {
				t.Error(err)
			}
			if ok != tt.want {
				t.Errorf("want: %+v, got: %+v\n", tt.want, ok)
			}
		})
	}
}

func TestDeleteIndex(t *testing.T) {
	tests := []struct {
		name    string
		index   string
		setup   func(tb testing.TB)
		cleanup func(tb testing.TB)
	}{
		{
			name:  "simple",
			index: "delete-v1",
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "delete-v1")
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
			err := es.deleteIndex(ctx, tt.index)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSyncCloseStatus(t *testing.T) {
	tests := []struct {
		name  string
		conf  config
		setup func(tb testing.TB)
	}{
		{
			name: "exists close",
			conf: config{
				Indices: []index{
					{
						Name:    "sync-close-v1",
						Mapping: "testdata/test.json",
						Status:  "close",
					},
				},
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "sync-close-v1")
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
			if tt.setup != nil {
				tt.setup(t)
			}
			err := es.syncCloseStatus(ctx, tt.conf)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
