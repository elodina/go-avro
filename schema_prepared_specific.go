package avro

import "reflect"

func specificDecoder(entry *structFieldPlan) preparedDecoder {
	switch entry.schema.Type() {
	case Record:
		return recordDec(entry.schema)
	case Enum:
		return enumDec(entry.schema.(*EnumSchema))
	default:
		// Generic decoders get less drastic speedups, but we can add more later.
		return genericDec(entry.schema)
	}
}

// structFieldPlan is a plan that assists in decoding
type structFieldPlan struct {
	name   string
	index  []int
	schema Schema
	dec    preparedDecoder
}

type preparedDecoder func(reflectField reflect.Value, dec Decoder) (reflect.Value, error)

func genericDec(schema Schema) preparedDecoder {
	return func(reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
		return sdr.readValue(schema, reflectField, dec)
	}
}

func enumDec(schema *EnumSchema) preparedDecoder {
	symbolsToIndex := NewGenericEnum(schema.Symbols).symbolsToIndex
	return func(reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
		enumIndex, err := dec.ReadEnum()
		if err != nil {
			return reflect.ValueOf(enumIndex), err
		}
		enum := &GenericEnum{
			Symbols:        schema.Symbols,
			symbolsToIndex: symbolsToIndex,
			index:          enumIndex,
		}
		return reflect.ValueOf(enum), nil
	}
}

func recordDec(schema Schema) preparedDecoder {
	return func(reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
		return sdr.mapRecord(schema, reflectField, dec)
	}
}
