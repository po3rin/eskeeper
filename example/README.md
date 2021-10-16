# A Tour of eskeeper

A Tour of eskeeper explains the basic usage of eskeeper.

The files to be used are prepared under the `example` directory.

```tree
example
├── README.md
├── docker-compose.yaml
├── es.yaml
├── README.md
└── test.json
```

## start Elasticsearch

First, start Elasticserarch.

```bash
$ docker-compose up -d --build
```

## create indices & aliases

Next, prepare index setting JSON file `test.json`. The content of this file is the same as the request body of the [Create index API](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html)

```json
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "test_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "title": {
        "type": "text"
      },
      "body": {
        "type": "text",
	"analyzer": "test_analyzer"
      }
    }
  }
}
```

create `es.yaml` to manage indices.

```yaml
index:
  - name: test-v1 # index name
    mapping: test.json # index setting & mapping (json)

  - name: test-v2
    mapping: test.json

  - name: close-v1
    mapping: test.json
    status: close

alias:
  - name: alias1
    index:
      - test-v1

  # multi indicies
  - name: alias2
    index:
      - test-v1
      - test-v2
```

You can see in this file that you create 3 indexes and 2 aliases. You can create an index with the same settings as the test.json created earlier. `close-v1` index use `status: close`. This allows closed indices to not have to maintain internal data structures for indexing or searching documents, resulting in a smaller overhead on the cluster.

`alias` field declares alias settings. You can redirection the index without downtime.

Now that we're ready, we can use eskeeper to create the index and alias.

```bash
eskeeper < es.yaml

curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b
yellow close close-v1 ... 1 1

curl localhost:9200/_cat/aliases
alias2 test-v2 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -
```

## change index setting and switch alias

Next, change mapping & config file.

```diff
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "test_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "title": {
        "type": "text"
      },
      "body": {
        "type": "text",
	"analyzer": "test_analyzer"
      },
+     "category": {
+	"type": "text"
+     }
  }
}
```

mapping update is not supported. so you should create new index.

```diff
index:
  - name: test-v1 # index name
    mapping: test.json # index setting & mapping (json)

  - name: test-v2
    mapping: test.json

  - name: close-v1
    mapping: test.json
    status: close

+ - name: test-v3
+   mapping: test.json

alias:
  - name: alias1
    index:
      - test-v1

  # multi indicies
  - name: alias2
    index:
      - test-v1
-     - test-v2
+     - test-v3
```

eskeeper creates new index ```test-v3``` with new mapping has `category` field. But, ```test-v1``` and ```test-v2``` setting & mapping is not changed (these index done not have `category` field).

```bash
eskeeper < es.yaml

curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b
yellow open test-v3 ... 1 1 0 0 208b 208b
yellow close close-v1 ... 1 1

curl localhost:9200/_cat/aliases
alias2 test-v3 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -
```

## reindex 

Finally, let's try reindex. put data to `test-v3` index.

```
curl -X POST localhost:9200/test-v3/_doc -H "Content-Type: application/json" -d '{
  "id": 1,
  "title": "how to use eskeeper",
  "body": "this is A Tour of eskeeper",
  "category": "Elasticsearch"
}'
```

checks `test-v3` has one doc.

```
curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b
yellow open test-v3 ... 1 1 1 0 4.1kb 4.1kb
yellow close close-v1 xxxxxxxxxxxx 1 1
```

To reidnex `test-v3` to new index `test-v4`, use reindex option.

```diff
index:
  - name: test-v1 # index name
    mapping: test.json # index setting & mapping (json)

  - name: test-v2
    mapping: test.json

  - name: close-v1
    mapping: test.json
    status: close

  - name: test-v3
    mapping: test.json

+ - name: test-v4
+   mapping: test.json
+   reindex:
+       source: test-v3
+       waitForCompletion: false
+       on: firstCreated

alias:
  - name: alias1
    index:
      - test-v1

  # multi indicies
  - name: alias2
    index:
      - test-v1
-     - test-v3
+     - test-v4
```

reindex field supports `source`, `slices`, `waitForComletion`, `on` options. `source`, `slices`, `waitForComletion` is same as [Reindex  API](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html) request body options. 

`on` field is eskeeper original option. If `on: firstCreated` is specified, reindex will only be executed the first time an index is created. If `on: always` is specified, a reindex will be run each time eskeeper is run.

```bash
eskeeper < es.yaml

curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b
yellow open test-v3 ... 1 1 1 0 4.1kb 4.1kb
yellow open test-v4 ... 1 1 1 0 4.1kb 4.1kb
yellow close close-v1 xxxxxxxxxxxx 1 1

curl localhost:9200/_cat/aliases
alias2 test-v4 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -
```

reindex is completed. check that `test-v4` has one doc.
