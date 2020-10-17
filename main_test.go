package eskeeper

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch"
	"github.com/ory/dockertest"
)

var url string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run(
		"docker.elastic.co/elasticsearch/elasticsearch",
		"7.9.2",
		[]string{
			"ES_JAVA_OPTS=-Xms512m -Xmx512m",
			"discovery.type=single-node",
			"node.name=es01",
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort("9200/tcp")
	url = fmt.Sprintf("http://localhost:%s", port)

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		res, err := es.Info()
		if err != nil {
			log.Println("waiting to be ready...")
			return err
		}
		defer res.Body.Close()
		return nil
	}); err != nil {
		log.Fatalf("could not retry to connect : %s\n", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreateIndex(t *testing.T) {
	tests := []struct {
		name string
		conf Config
	}{
		{
			name: "simple",
			conf: Config{
				Index: []Index{
					{
						Name:    "create1",
						Mapping: "testdata/test.json",
					},
				},
			},
		},
		{
			name: "multi",
			conf: Config{
				Index: []Index{
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
			conf: Config{
				Index: []Index{
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
	}

	es, err := newEsClient([]string{url}, "", "")
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := es.createIndex(ctx, tt.conf)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestSyncAlias(t *testing.T) {
	tests := []struct {
		name    string
		conf    Config
		setup   func(tb testing.TB)
		cleanup func(tb testing.TB)
	}{
		{
			name: "simple",
			conf: Config{
				Alias: []Alias{
					{
						Name:  "test-alias",
						Index: []string{"test-v1", "test-v2"},
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
			conf: Config{
				Alias: []Alias{
					{
						Name:  "test-alias",
						Index: []string{"test-v3"},
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
			err := es.syncAlias(ctx, tt.conf)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

var testMapping = `{
    "mappings": {
        "properties": {
            "id": {
                "type": "long",
                "index": true
            },
            "title": {
                "type": "text"
            },
            "body": {
                "type": "text"
            }
        }
    }
}`

func createTmpIndexHelper(tb testing.TB, name string) {
	tb.Helper()
	conf := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		tb.Fatal(err)
	}

	res, err := es.Indices.Create(
		name,
		es.Indices.Create.WithBody(strings.NewReader(testMapping)),
	)
	if err != nil {
		tb.Fatal(err)
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			tb.Fatal(err)
		}
		tb.Fatalf("failed to create index [index= %v, statusCode=%v, res=%v]", name, res.StatusCode, string(body))
	}

}
