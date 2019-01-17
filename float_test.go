package null

import (
	"encoding/json"
	"math"
	"testing"
)

var (
	floatJSON       = []byte(`1.2345`)
	floatStringJSON = []byte(`"1.2345"`)
	floatBlankJSON  = []byte(`""`)
	nullFloatJSON   = []byte(`{"Float64":1.2345,"Valid":true}`)
)

func TestFloatFrom(t *testing.T) {
	f := FloatFrom(1.2345)
	assertFloat(t, f, "FloatFrom()")

	zero := FloatFrom(0)
	if !zero.Valid {
		t.Error("FloatFrom(0)", "is invalid, but should be valid")
	}
}

func TestFloatFromPtr(t *testing.T) {
	n := float64(1.2345)
	iptr := &n
	f := FloatFromPtr(iptr)
	assertFloat(t, f, "FloatFromPtr()")

	null := FloatFromPtr(nil)
	assertNullFloat(t, null, "FloatFromPtr(nil)")
}

func TestUnmarshalFloat(t *testing.T) {
	var f Float
	err := json.Unmarshal(floatJSON, &f)
	maybePanic(err)
	assertFloat(t, f, "float json")

	var sf Float
	err = json.Unmarshal(floatStringJSON, &sf)
	maybePanic(err)
	assertFloat(t, sf, "string float json")

	var nf Float
	err = json.Unmarshal(nullFloatJSON, &nf)
	maybePanic(err)
	assertFloat(t, nf, "sql.NullFloat64 json")

	var null Float
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullFloat(t, null, "null json")

	var blank Float
	err = json.Unmarshal(floatBlankJSON, &blank)
	maybePanic(err)
	assertNullFloat(t, blank, "null blank string json")

	var badType Float
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullFloat(t, badType, "wrong type json")

	var invalid Float
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
}

func TestTextUnmarshalFloat(t *testing.T) {
	var f Float
	err := f.UnmarshalText([]byte("1.2345"))
	maybePanic(err)
	assertFloat(t, f, "UnmarshalText() float")

	var blank Float
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullFloat(t, blank, "UnmarshalText() empty float")

	var null Float
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullFloat(t, null, `UnmarshalText() "null"`)
}

func TestMarshalFloat(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := json.Marshal(f)
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalFloatText(t *testing.T) {
	f := FloatFrom(1.2345)
	data, err := f.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "1.2345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewFloat(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestFloatPointer(t *testing.T) {
	f := FloatFrom(1.2345)
	ptr := f.Ptr()
	if *ptr != 1.2345 {
		t.Errorf("bad %s float: %#v ≠ %v\n", "pointer", ptr, 1.2345)
	}

	null := NewFloat(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s float: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestFloatIsZero(t *testing.T) {
	f := FloatFrom(1.2345)
	if f.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewFloat(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewFloat(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestFloatSetValid(t *testing.T) {
	change := NewFloat(0, false)
	assertNullFloat(t, change, "SetValid()")
	change.SetValid(1.2345)
	assertFloat(t, change, "SetValid()")
}

func TestFloatScan(t *testing.T) {
	var f Float
	err := f.Scan(1.2345)
	maybePanic(err)
	assertFloat(t, f, "scanned float")

	var sf Float
	err = sf.Scan("1.2345")
	maybePanic(err)
	assertFloat(t, sf, "scanned string float")

	var null Float
	err = null.Scan(nil)
	maybePanic(err)
	assertNullFloat(t, null, "scanned null")
}

func TestFloatInfNaN(t *testing.T) {
	nan := NewFloat(math.NaN(), true)
	_, err := nan.MarshalJSON()
	if err == nil {
		t.Error("expected error for NaN, got nil")
	}

	inf := NewFloat(math.Inf(1), true)
	_, err = inf.MarshalJSON()
	if err == nil {
		t.Error("expected error for Inf, got nil")
	}
}

func TestFloatValueOrZero(t *testing.T) {
	valid := NewFloat(1.2345, true)
	if valid.ValueOrZero() != 1.2345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewFloat(1.2345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestFloatPlus(t *testing.T) {
	a := NewFloat(1.2345, true)
	b := NewFloat(5.2, true)
	c := NewFloat(0, false)

	if a.Plus(b).ValueOrZero() != 6.4345 {
		t.Error("unexpected ValueOrZero", a.Plus(b).ValueOrZero())
	}
	if b.Plus(a).ValueOrZero() != 6.4345 {
		t.Error("unexpected ValueOrZero", b.Plus(a).ValueOrZero())
	}
	if a.Plus(c).Valid {
		t.Error("unexpected Valid", a.Plus(c).Valid)
	}
	if c.Plus(a).Valid {
		t.Error("unexpected Valid", a.Plus(c).Valid)
	}
	if c.Plus(c).Valid {
		t.Error("unexpected Valid", a.Plus(c).Valid)
	}
}

func TestFloatMinus(t *testing.T) {
	a := NewFloat(1.2, true)
	b := NewFloat(5.2, true)
	c := NewFloat(0, false)

	if a.Minus(b).ValueOrZero() != 1.2 - 5.2 {
		t.Error("unexpected ValueOrZero", a.Minus(b).ValueOrZero())
	}
	if b.Minus(a).ValueOrZero() != 5.2 - 1.2 {
		t.Error("unexpected ValueOrZero", b.Minus(a).ValueOrZero())
	}
	if a.Minus(c).Valid {
		t.Error("unexpected Valid", a.Minus(c).Valid)
	}
	if c.Minus(a).Valid {
		t.Error("unexpected Valid", a.Plus(c).Valid)
	}
	if c.Minus(c).Valid {
		t.Error("unexpected Valid", a.Plus(c).Valid)
	}
}

func TestFloatDividedBy(t *testing.T) {
	a := NewFloat(10, true)
	b := NewFloat(4, true)
	c := NewFloat(0, false)
	d := NewFloat(0, true)

	if a.DividedBy(b).ValueOrZero() != 2.5 {
		t.Error("unexpected ValueOrZero", a.DividedBy(b).ValueOrZero())
	}
	if b.DividedBy(a).ValueOrZero() != 0.4 {
		t.Error("unexpected ValueOrZero", b.DividedBy(a).ValueOrZero())
	}
	if a.DividedBy(c).Valid {
		t.Error("unexpected Valid", a.DividedBy(c).Valid)
	}
	if c.DividedBy(c).Valid {
		t.Error("unexpected Valid", c.DividedBy(c).Valid)
	}
	if c.DividedBy(a).Valid {
		t.Error("unexpected Valid", c.DividedBy(a).Valid)
	}
	if a.DividedBy(d) != FloatFrom(math.Inf(1)) {
		t.Error("unexpected Value", a.DividedBy(d))
	}
}

func TestFloatTimes(t *testing.T) {
	a := NewFloat(1.2, true)
	b := NewFloat(5.2, true)
	c := NewFloat(0, false)

	if a.Times(b).ValueOrZero() != 6.24 {
		t.Error("unexpected ValueOrZero", a.Times(b).ValueOrZero())
	}
	if b.Times(a).ValueOrZero() != 6.24 {
		t.Error("unexpected ValueOrZero", b.Times(a).ValueOrZero())
	}
	if a.Times(c).Valid {
		t.Error("unexpected Valid", a.Times(c).Valid)
	}
	if c.Times(a).Valid {
		t.Error("unexpected Valid", a.Times(c).Valid)
	}
	if c.Times(c).Valid {
		t.Error("unexpected Valid", a.Times(c).Valid)
	}
}

func assertFloat(t *testing.T, f Float, from string) {
	if f.Float64 != 1.2345 {
		t.Errorf("bad %s float: %f ≠ %f\n", from, f.Float64, 1.2345)
	}
	if !f.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullFloat(t *testing.T, f Float, from string) {
	if f.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
