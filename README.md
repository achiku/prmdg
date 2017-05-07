# prmdg

[![Build Status](https://travis-ci.org/achiku/prmdg.svg?branch=master)](https://travis-ci.org/achiku/prmdg)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/achiku/prmdg/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/achiku/prmdg)](https://goreportcard.com/report/github.com/achiku/prmdg)

prmd format JSON Hyper Schema to Go structs, and validators


## Why created

## Installation

```
go get -u github.com/achiku/prmdg
```

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
