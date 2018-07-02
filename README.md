# prmdg

[![Build Status](https://travis-ci.org/achiku/prmdg.svg?branch=master)](https://travis-ci.org/achiku/prmdg)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/achiku/prmdg/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/achiku/prmdg)](https://goreportcard.com/report/github.com/achiku/prmdg)

prmd style JSON Hyper Schema to Go structs, and validators


## Why created

## Prerequisite

```
go get -u golang.org/x/tools/cmd/goimports
```

`prmdg` applies `goimports` to the ourput file.

## Installation

```
go get -u github.com/achiku/prmdg
```

If you want to use `github.com/gureg/null` in Go struct by adding `--nullable` option, you need to install `github.com/gureg/null` first.

## Usage

```
usage: prmdg --file=FILE [<flags>] <command> [<args> ...]

prmd generated JSON Hyper Schema to Go

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -p, --package="main"  package name for Go file
  -f, --file=FILE       path JSON Schema
  -o, --output=OUTPUT   path to Go output file

Commands:
  help [<command>...]
    Show help.

  struct [<flags>]
    generate struct file

  jsval
    generate validator file using github.com/lestrrat/go-jsval

  validator
    generate validator file using github.com/go-playground/validator

```

## Generating struct from JSON Hyper Schema

```
usage: prmdg struct [<flags>]

generate struct file

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -p, --package="main"  package name for Go file
  -f, --file=FILE       path JSON Schema
  -o, --output=OUTPUT   path to Go output file
      --validate-tag    add `validate` tag to struct
      --use-title       use title tag in request/response struct name
      --nullable        use github.com/guregu/null for null value
```


## Generating validator from JSON Hyper Schema

```
usage: prmdg jsval

generate validator file using github.com/lestrrat/go-jsval

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -p, --package="main"  package name for Go file
  -f, --file=FILE       path JSON Schema
  -o, --output=OUTPUT   path to Go output file

```


```
usage: prmdg validator

generate validator file using github.com/go-playground/validator

Flags:
      --help            Show context-sensitive help (also try --help-long and --help-man).
  -p, --package="main"  package name for Go file
  -f, --file=FILE       path JSON Schema
  -o, --output=OUTPUT   path to Go output file

```
