# eskeeper

<img src="https://img.shields.io/badge/go-v1.15-blue.svg"/> [![GoDoc](https://godoc.org/github.com/po3rin/eskeeper?status.svg)](https://godoc.org/github.com/po3rin/eskeeper)

eskeeper synchronizes index and alias with configuration files while ensuring idempotency.

## Quick Start

eskeeper recieves yaml format data form stdin.

```bash
$ go get -u github.com/po3rin/eskeeper/cmd/eskeeper
$ eskeeper < es.yaml
```

## Features

* index
- [x] create
- [ ] update mapping
- [ ] delete

* alias
- [x] create
- [x] update
- [ ] delete
