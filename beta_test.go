package eskeeper

// func TestUpdateIndex(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		index   index
// 		setup   func(tb testing.TB)
// 		wantErr bool
// 	}{
// 		{
// 			name: "simple",
// 			index: index{
// 				Name:    "update-test-v1",
// 				Mapping: "testdata/updateIndex.json",
// 			},
// 			setup: func(tb testing.TB) {
// 				createTmpIndexHelper(tb, "update-test-v1")
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "same",
// 			index: index{
// 				Name:    "update-test-v2",
// 				Mapping: "testdata/test.json",
// 			},
// 			setup: func(tb testing.TB) {
// 				createTmpIndexHelper(tb, "update-test-v2")
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "invalid-update",
// 			index: index{
// 				Name:    "update-test-v3",
// 				Mapping: "testdata/invalidUpdateIndex.json",
// 			},
// 			setup: func(tb testing.TB) {
// 				createTmpIndexHelper(tb, "update-test-v3")
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	es, err := newEsClient([]string{url}, "", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx := context.Background()
// 			tt.setup(t)
// 			err := es.updateIndex(ctx, tt.index)
// 			if tt.wantErr && err == nil {
// 				t.Error("expect error")
// 			}
// 			if !tt.wantErr && err != nil {
// 				t.Error(err)
// 			}
// 		})
// 	}
// }
