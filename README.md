# eskeeper

<img src="https://img.shields.io/badge/go-v1.15-blue.svg"/> [![GoDoc](https://godoc.org/github.com/po3rin/eskeeper?status.svg)](https://godoc.org/github.com/po3rin/eskeeper) ![Go Test](https://github.com/po3rin/eskeeper/workflows/Go%20Test/badge.svg) [![Apache-2.0](https://img.shields.io/github/license/po3rin/eskeeper)](LICENSE)

eskeeper synchronizes index and alias with configuration files while ensuring idempotency. It still only supports WRITE. DELETE is not yet supported because the operation of deleting persistent data is dangerous and needs to be implemented carefully. 

## :muscle: Currently supports

### mode

- [x] CLI mode
- [ ] Agent mode

### sync 

* index
- [x] create
- [ ] update mapping
- [ ] delete

* alias
- [x] create
- [x] update
- [ ] delete


## :four_leaf_clover: Quick Start

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
curl localhost:9200/_cat/indices
yellow open test-v1 ... 1 1 0 0 208b 208b
yellow open test-v2 ... 1 1 0 0 208b 208b

curl localhost:9200/_cat/aliases
alias2 test-v2 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -
```

## :triangular_ruler: Settings

eskeeper supports flag & environment value.

```bash
# use flags
eskeeper -u user -p pass -e=http://localhost:9200,http://localhost9300 < testdata/es.yaml

# use env
ESKEEPER_ES_USER=user ESKEEPER_ES_PASS=pass ESKEEPER_ES_URLS=http://localhost:9200 eskeeper < testdata/es.yaml
```

## :mag_right: Stages

eskeeper process is divided into four stages.

#### validation config stage
* Validates config yaml format

#### pre-check stage 

* Check if mapping file is valid format
* Check if there is an index for alias  

#### sync stage
* Sync indices and aliases with config

#### post-check stage
* Check if indices & aliases has been created


## :triangular_flag_on_post: Contributing

Did you find something technically wrong, something to fix, or something? Please give me Issue or Pull Request !!

### Test

eskeeper's test uses [github.com/ory/dockertest](github.com/ory/dockertest). So you need docker to test.
