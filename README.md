Apache Avro for Golang
=====================
(forked from `elodina/go-avro`) 

[![Build Status](https://travis-ci.org/go-avro/avro.svg?branch=master)](https://travis-ci.org/go-avro/avro)

About This fork
---------------

This fork has separated from elodina/go-avro in December 2016 because of the project not responding to PR's since around May 2016. Have tried to contact them to get maintainer access but the original maintainer no longer is able to make those changes, so I've forked currently. If elodina/go-avro returns, they are free to merge all the changes I've made back.


Documentation
-------------

Installation is as easy as follows:

`go get gopkg.in/avro.v0`

Some usage examples are located in [examples folder](https://github.com/go-avro/avro/tree/master/examples):

* [DataFileReader](https://github.com/go-avro/avro/blob/master/examples/data_file/data_file.go)
* [GenericDatumReader/Writer](https://github.com/go-avro/avro/blob/master/examples/generic_datum/generic_datum.go)
* [SpecificDatumReader/Writer](https://github.com/go-avro/avro/blob/master/examples/specific_datum/specific_datum.go)
* [Schema loading](https://github.com/go-avro/avro/blob/master/examples/load_schema/load_schema.go)
* Code gen support available in [codegen folder](https://github.com/go-avro/avro/tree/master/codegen)
