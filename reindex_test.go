package eskeeper

import (
	"context"
	"testing"
)

func TestReindex(t *testing.T) {
	tests := []struct {
		name    string
		dest    string
		reindex reindex
		setup   func(tb testing.TB)
		cleanup func(tb testing.TB)
	}{
		{
			name: "reindex",
			dest: "reindex-dest",
			reindex: reindex{
				Source:            "reindex-src",
				Slices:            3,
				WaitForCompletion: true,
				On:                "firstCreated",
			},
			setup: func(tb testing.TB) {
				createTmpIndexHelper(tb, "reindex-src")
				createTmpIndexHelper(tb, "reindex-dest")
				postDocHelper(tb, "reindex-src")
			},
			cleanup: func(tb testing.TB) {
				deleteIndexHelper(tb, []string{"reindex-src", "reindex-dest"})
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
			err := es.reindex(ctx, tt.dest, tt.reindex)
			if err != nil {
				t.Error(err)
			}
			if tt.cleanup != nil {
				tt.cleanup(t)
			}
		})
	}
}
