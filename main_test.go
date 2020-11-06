package eskeeper

import (
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

func createTmpIndexHelper(tb testing.TB, name string) {
	tb.Helper()
	conf := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		tb.Fatal(err)
	}

	f, err := os.Open("testdata/test.json")
	if err != nil {
		tb.Fatal(err)
	}

	res, err := es.Indices.Create(
		name,
		es.Indices.Create.WithBody(f),
	)
	if err != nil {
		tb.Fatal(err)
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			tb.Fatal(err)
		}
		tb.Fatalf("failed to create index [index=%v, statusCode=%v, res=%v]", name, res.StatusCode, string(body))
	}
}

var testAliasQuery = `
{
  "actions" : [
    { "add" : { "index" : "%v", "alias" : "%v" } }
  ]
}`

func createTmpAliasHelper(tb testing.TB, name string, index string) {
	tb.Helper()
	conf := elasticsearch.Config{
		Addresses: []string{url},
	}
	es, err := elasticsearch.NewClient(conf)
	if err != nil {
		tb.Fatal(err)
	}

	res, err := es.Indices.UpdateAliases(
		strings.NewReader(fmt.Sprintf(testAliasQuery, index, name)),
	)
	if err != nil {
		tb.Fatal(err)
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			tb.Fatal(err)
		}
		tb.Fatalf("failed to create alias [index= %v, statusCode=%v, res=%v]", name, res.StatusCode, string(body))
	}
}
