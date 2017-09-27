# Changelog

#### Version 0.2 (not yet released)

Intention: start making changes towards a 1.0 release.

API Changes:
 - The `BinaryEncoder` type is now a private type. `avro.NewBinaryEncoder()` 
   now returns a value of the `Encoder` interface.
 - The `BinaryDecoder` type is now also a private type. `avro.NewBinaryDecoder()`
   now returns a value of the `Decoder` interface.
 - Rename the `Writer` and `Reader` interfaces to `Marshaler` and `Unmarshaler` to
   be more like the JSON encoder and also use similar method names.
 - Rename error types `FooBar` to be `ErrFooBar`

Improvements:
 - Major improvement to docs and docs coverage


#### Version 0.1 (2017-08-23)

 - First version after forking from elodina.
 - Started a semver-considering API, using the gopkg.in interface, 
   and planning for a 1.x release.

Improvements:
 - Error reporting: specify which field is missing when throwing FieldDoesNotExist 
   [#5](https://github.com/go-avro/avro/pull/5)
 - Speedup encoding for strings and bools 
   [#6](https://github.com/go-avro/avro/pull/6)
 - Can prepare schemas which are self-recursive and co-recursive.

Bug Fixes:
 - Can decode maps of non-primitive types [#2](https://github.com/go-avro/avro/pull/2)
 - Fix encoding of 'fixed' type [#3](https://github.com/go-avro/avro/pull/3) [elodina/#78](https://github.com/elodina/go-avro/issues/78)
 - Fix encoding of boolean when used in a type union [#4](https://github.com/go-avro/avro/pull/4)