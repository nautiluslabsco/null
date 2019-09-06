package null

import (
	"math"
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

type floatData struct {
	Float            float64
	FloatNaN         float64
	FloatPtr         *float64
	FloatPtrNaN      *float64
	FloatPtrNil      *float64
	NullFloat        Float
	NullFloatNaN     Float
	NullFloatInf     Float
	NullFloatInvalid Float
}

func TestJSONEncodeDecode(t *testing.T) {
	json := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
	json.RegisterExtension(JSONFloatDecoderExtension())
	json.RegisterExtension(JSONFloatEncoderExtension())

	data := floatData{
		10.0,
		math.NaN(),
		FloatFrom(11.0).Ptr(),
		FloatFrom(math.NaN()).Ptr(),
		nil,
		FloatFrom(10.0),
		FloatFrom(math.NaN()),
		FloatFrom(math.Inf(1)),
		Float{},
	}
	bs, err := json.Marshal(data)
	assert.NoError(t, err)

	var decodedData floatData
	err = json.Unmarshal(bs, &decodedData)
	assert.NoError(t, err)

	assert.Equal(t, math.IsNaN(*decodedData.FloatPtrNaN), math.IsNaN(*data.FloatPtrNaN))
	decodedData.FloatPtrNaN = nil
	data.FloatPtrNaN = nil

	assert.Equal(t, math.IsNaN(decodedData.NullFloatNaN.Float64), math.IsNaN(data.NullFloatNaN.Float64))
	decodedData.NullFloatNaN = Float{}
	data.NullFloatNaN = Float{}

	assert.Equal(t, math.IsNaN(decodedData.NullFloatNaN.Float64), math.IsNaN(data.NullFloatNaN.Float64))
	decodedData.FloatNaN = 0
	data.FloatNaN = 0

	assert.True(t, reflect.DeepEqual(decodedData, data))
}
