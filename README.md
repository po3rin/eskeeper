# eskeeper

<img src="https://img.shields.io/badge/go-v1.15-blue.svg"/> [![GoDoc](https://godoc.org/github.com/po3rin/eskeeper?status.svg)](https://godoc.org/github.com/po3rin/eskeeper) [![Apache-2.0](https://img.shields.io/github/license/po3rin/eskeeper)](LICENSE)

*TODO: impliments*

eskeeper synchronizes index and alias with configuration files while ensuring idempotency. It still only supports WRITE. DELETE is not yet supported because the operation of deleting persistent data is dangerous and needs to be implemented carefully. 

## Quick Start

eskeeper recieves yaml format data form stdin.

```bash
$ go get -u github.com/po3rin/eskeeper/cmd/eskeeper
$ eskeeper < es.yaml
```

es.yaml is indices & aliases config file.

```yaml
index:
  - name: test-v1
    mapping: testdata/test.json

  - name: test-v2
    mapping: testdata/test.json

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

results

```bach
curl localhost:9200/_cat/aliases
alias2 test-v2 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -

curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b
```

## Features

* index
- [x] create
- [ ] update mapping
- [ ] delete

* alias
- [x] create
- [x] update (⚠️destructive change)
- [ ] delete
