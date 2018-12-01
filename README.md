# tapper [![GoDoc](https://godoc.org/github.com/gmendonca/tapper?status.svg)](https://godoc.org/github.com/gmendonca/tapper) [![Go Report Card](https://goreportcard.com/badge/github.com/gmendonca/tapper)](https://goreportcard.com/report/github.com/gmendonca/tapper) [![Build Status](https://travis-ci.com/gmendonca/tapper.svg?branch=master)](https://travis-ci.com/gmendonca/tapper)

tapper is a Go program to process and analyse Hive and Presto query logs.

## Introduction

tapper uses the logs coming from an ES to process and analyse the query log.
[query-log-parser](https://github.com/gmendonca/query-log-parser) is an attempt to get the log
from Hive logs and Presto UI to make it eligible to be processed. We could read the [documentation](https://github.com/gmendonca/query-log-parser/blob/master/README.md)
for more information on the format of the queries. From that logs some metrics are posted to Datadog.

## Installation

```sh
$ go get github.com/gmendonca/tapper
$ which tapper
$GOPATH/bin/tapper
```

## Usage

```sh
$ tapper logs --config ./configs/config.json
```

