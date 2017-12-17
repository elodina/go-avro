# Apache Avro for Golang

[![Build Status](https://travis-ci.org/go-avro/avro.svg?branch=master)](https://travis-ci.org/go-avro/avro) [![GoDoc](https://godoc.org/gopkg.in/avro.v0?status.svg)](https://godoc.org/gopkg.in/avro.v0)


Support for decoding/encoding avro using both map-style access (GenericRecord) and to/from arbitrary Go structs (SpecificRecord).

This library started as a fork of `elodina/go-avro` but has now proceeded to become a maintained library.

## Installation

Installation via go get:

    go get gopkg.in/avro.v0


## Documentation

 * [Read API/usage docs on Godoc](https://godoc.org/gopkg.in/avro.v0)
 * [Changelog](CHANGELOG.md)

Some usage examples are located in [examples folder](https://github.com/go-avro/avro/tree/master/examples):

* [DataFileReader](https://github.com/go-avro/avro/blob/master/examples/data_file/data_file.go)
* [GenericDatumReader/Writer](https://github.com/go-avro/avro/blob/master/examples/generic_datum/generic_datum.go)
* [SpecificDatumReader/Writer](https://github.com/go-avro/avro/blob/master/examples/specific_datum/specific_datum.go)
* [Schema loading](https://github.com/go-avro/avro/blob/master/examples/load_schema/load_schema.go)
* Code gen support available in [codegen folder](https://github.com/go-avro/avro/tree/master/codegen)


## About This fork

This fork separated from elodina/go-avro in December 2016 because of the
project not responding to PR's since around May 2016. Had tried to contact them
to get maintainer access but the original maintainer no longer is able to make
those changes.

Originally, we were waiting in hope the elodina maintainer would return, but it
hasn't happened, so the plan now is to proceed with this as its own library and
take PRs, push for feature additions and version bumps.
