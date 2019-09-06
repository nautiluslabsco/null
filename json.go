package null

import (
	"math"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

func JSONFloatEncoderExtension() jsoniter.Extension {
	return jsoniter.EncoderExtension(map[reflect2.Type]jsoniter.ValEncoder{
		reflect2.DefaultTypeOfKind(reflect.Float64): floatEncDec{},
		reflect2.TypeByName("null.Float"):           nullFloatEncDec{},
	})
}

func JSONFloatDecoderExtension() jsoniter.Extension {
	return jsoniter.DecoderExtension(map[reflect2.Type]jsoniter.ValDecoder{
		reflect2.DefaultTypeOfKind(reflect.Float64): floatEncDec{},
		reflect2.TypeByName("null.Float"):           nullFloatEncDec{},
	})
}

func decode(iter *jsoniter.Iterator) Float {
	valueType := iter.WhatIsNext()
	switch valueType {
	case jsoniter.NilValue:
		iter.ReadNil()
		return Float{}
	case jsoniter.NumberValue:
		return FloatFrom(iter.ReadFloat64())
	case jsoniter.StringValue:
		switch str := iter.ReadString(); str {
		case "+Inf":
			return FloatFrom(math.Inf(+1))
		case "-Inf":
			return FloatFrom(math.Inf(-1))
		case "NaN":
			return FloatFrom(math.NaN())
		default:
			iter.ReportError("fuzzyFloat64Decoder", "not number \"NaN\", \"+Inf\", or \"-Inf\"")
		}
	default:
		iter.ReportError("fuzzyFloat64Decoder", "not number or string")
	}
	return Float{}
}

func encode(val Float, stream *jsoniter.Stream) {
	if val.Valid {
		fval := val.Float64
		if math.IsNaN(fval) {
			stream.WriteString("NaN")
		} else if math.IsInf(fval, 1) {
			stream.WriteString("+Inf")
		} else if math.IsInf(fval, -1) {
			stream.WriteString("-Inf")
		} else {
			stream.WriteFloat64(fval)
		}
	} else {
		stream.WriteNil()
	}

}

type nullFloatEncDec struct {
}

func (e nullFloatEncDec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	val := *((*Float)(ptr))
	encode(val, stream)
}

func (e nullFloatEncDec) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (d nullFloatEncDec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	val := decode(iter)
	*((*Float)(ptr)) = val
}

type floatEncDec struct {
}

func (e floatEncDec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	fval := *((*float64)(ptr))
	encode(FloatFrom(fval), stream)
}

func (e floatEncDec) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (d floatEncDec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	val := decode(iter)
	*((*float64)(ptr)) = val.Float64
}
