# eskeeper

<img src="https://img.shields.io/badge/go-v1.15-blue.svg"/> [![GoDoc](https://godoc.org/github.com/po3rin/eskeeper?status.svg)](https://godoc.org/github.com/po3rin/eskeeper) ![Go Test](https://github.com/po3rin/eskeeper/workflows/Go%20Test/badge.svg) [![Apache-2.0](https://img.shields.io/github/license/po3rin/eskeeper)](LICENSE)

eskeeper synchronizes index and alias with configuration files while ensuring idempotency. It still only supports WRITE. DELETE is not yet supported because the operation of deleting persistent data is dangerous and needs to be implemented carefully. 

## :muscle: Currently supports

### mode

- [x] CLI mode
- [x] Index status(open/close only)
- [ ] Agent mode

### sync 

* index
- [x] create
- [x] status (open or close)
- [x] reindex (only basic parameter)
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
  - name: test-v1 # index name
    mapping: testdata/test.json # index setting & mapping (json)

  - name: test-v2
    mapping: testdata/test.json

  - name: close-v1
    mapping: testdata/test.json
    status: close

  # reindex test-v1 -> reindex-v1	
  - name: reindex-v1
    mapping: testdata/test.json
    reindex:
        source: test-v1 
        slices: 3 # default=1
        waitForCompletion: true

        # 'on' field supports 2 hooks.
        # 'reindex': only when index is created for the first time.
        # 'always': always exec reindex.
        on: firstCreated


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
yellow close close-v1 xxxxxxxxxxxx 1 1
yellow open reindex-v1 ... 1 1 0 0 208b 208b

curl localhost:9200/_cat/aliases
alias2 test-v2 - - - -
alias1 test-v1 - - - -
alias2 test-v1 - - - -
```


## :triangular_ruler: Usage

eskeeper supports flag & environment value.

```bash
# use flags
eskeeper -u user -p pass -e=http://localhost:9200,http://localhost9300 < testdata/es.yaml

# use env
ESKEEPER_ES_USER=user ESKEEPER_ES_PASS=pass ESKEEPER_ES_URLS=http://localhost:9200 eskeeper < testdata/es.yaml
```

## :mag_right: Stages

eskeeper process is divided into four stages. verbose option lets you know eskeeper details.

```
$ eskeeper < es.yaml
loading config ...

=== validation stage ===
[pass] index: test-v1
[pass] index: test-v2
[pass] index: close-v1
[pass] alias: alias1
[pass] alias: alias2

=== pre-check stage ===
[pass] index: test-v1
[pass] index: test-v2
[pass] index: close-v1
[pass] alias: alias1
[pass] alias: alias2

=== sync stage ===
[synced] index: test-v1
[synced] index: test-v2
[synced] index: close-v1
[synced] alias: alias1
[synced] alias: alias2

=== post-check stage ===
[pass] index: test-v1
[pass] index: test-v2
[pass] index: close-v1
[pass] alias: alias1
[pass] alias: alias2

succeeded
```

#### validation stage
* Validates config yaml format

#### pre-check stage 

* Check if mapping file is valid format
* Check if there is an index for alias  

#### sync stage
* Sync indices and aliases with config

The order of synchronization is as follows.

```
create index
↓
open index
↓
update alias
↓
close index
```

Index close operation should be done after switching the alias.
Because there can be downtime before switching aliases.

#### post-check stage
* Check if indices & aliases has been created


## :triangular_flag_on_post: Contributing

Did you find something technically wrong, something to fix, or something? Please give me Issue or Pull Request !!

### Test

eskeeper's test uses [github.com/ory/dockertest](github.com/ory/dockertest). So you need docker to test.
