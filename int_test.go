package null

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"
)

var (
	intJSON       = []byte(`12345`)
	intStringJSON = []byte(`"12345"`)
	nullIntJSON   = []byte(`{"Int64":12345,"Valid":true}`)
)

func TestIntFrom(t *testing.T) {
	i := IntFrom(12345)
	assertInt(t, i, "IntFrom()")

	zero := IntFrom(0)
	if !zero.Valid {
		t.Error("IntFrom(0)", "is invalid, but should be valid")
	}
}

func TestIntFromPtr(t *testing.T) {
	n := int64(12345)
	iptr := &n
	i := IntFromPtr(iptr)
	assertInt(t, i, "IntFromPtr()")

	null := IntFromPtr(nil)
	assertNullInt(t, null, "IntFromPtr(nil)")
}

func TestUnmarshalInt(t *testing.T) {
	var i Int
	err := json.Unmarshal(intJSON, &i)
	maybePanic(err)
	assertInt(t, i, "int json")

	var si Int
	err = json.Unmarshal(intStringJSON, &si)
	maybePanic(err)
	assertInt(t, si, "int string json")

	var ni Int
	err = json.Unmarshal(nullIntJSON, &ni)
	maybePanic(err)
	assertInt(t, ni, "sql.NullInt64 json")

	var bi Int
	err = json.Unmarshal(floatBlankJSON, &bi)
	maybePanic(err)
	assertNullInt(t, bi, "blank json string")

	var null Int
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullInt(t, null, "null json")

	var badType Int
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullInt(t, badType, "wrong type json")

	var invalid Int
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullInt(t, invalid, "invalid json")
}

func TestUnmarshalNonIntegerNumber(t *testing.T) {
	var i Int
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to int")
	}
}

func TestUnmarshalInt64Overflow(t *testing.T) {
	int64Overflow := uint64(math.MaxInt64)

	// Max int64 should decode successfully
	var i Int
	err := json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	int64Overflow++
	err = json.Unmarshal([]byte(strconv.FormatUint(int64Overflow, 10)), &i)
	if err == nil {
		panic("err should be present; decoded value overflows int64")
	}
}

func TestTextUnmarshalInt(t *testing.T) {
	var i Int
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertInt(t, i, "UnmarshalText() int")

	var blank Int
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullInt(t, blank, "UnmarshalText() empty int")

	var null Int
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullInt(t, null, `UnmarshalText() "null"`)
}

func TestMarshalInt(t *testing.T) {
	i := IntFrom(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalIntText(t *testing.T) {
	i := IntFrom(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewInt(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestIntPointer(t *testing.T) {
	i := IntFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s int: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewInt(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s int: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestIntIsZero(t *testing.T) {
	i := IntFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewInt(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewInt(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestIntSetValid(t *testing.T) {
	change := NewInt(0, false)
	assertNullInt(t, change, "SetValid()")
	change.SetValid(12345)
	assertInt(t, change, "SetValid()")
}

func TestIntScan(t *testing.T) {
	var i Int
	err := i.Scan(12345)
	maybePanic(err)
	assertInt(t, i, "scanned int")

	var null Int
	err = null.Scan(nil)
	maybePanic(err)
	assertNullInt(t, null, "scanned null")
}

func TestIntValueOrZero(t *testing.T) {
	valid := NewInt(12345, true)
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewInt(12345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestIntPlus(t *testing.T) {
	a := NewInt(1, true)
	b := NewInt(5, true)
	c := NewInt(0, false)

	if a.Plus(b).ValueOrZero() != 6 {
		t.Error("unexpected ValueOrZero", a.Plus(b).ValueOrZero())
	}
	if b.Plus(a).ValueOrZero() != 6 {
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

func TestIntMinus(t *testing.T) {
	a := NewInt(1, true)
	b := NewInt(5, true)
	c := NewInt(0, false)

	if a.Minus(b).ValueOrZero() != -4 {
		t.Error("unexpected ValueOrZero", a.Minus(b).ValueOrZero())
	}
	if b.Minus(a).ValueOrZero() != 4 {
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

func TestIntDividedBy(t *testing.T) {
	a := NewInt(1, true)
	b := NewInt(5, true)
	c := NewInt(0, false)
	d := NewInt(0, true)

	if a.DividedBy(b).ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", a.DividedBy(b).ValueOrZero())
	}
	if b.DividedBy(a).ValueOrZero() != 5 {
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
	if a.DividedBy(d).Valid {
		t.Error("unexpected Valid", a.DividedBy(d))
	}
}

func TestIntTimes(t *testing.T) {
	a := NewInt(2, true)
	b := NewInt(5, true)
	c := NewInt(0, false)
	d := NewInt(0, true)

	if a.Times(b).ValueOrZero() != 10 {
		t.Error("unexpected ValueOrZero", a.Times(b).ValueOrZero())
	}
	if b.Times(a).ValueOrZero() != 10 {
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
	v := a.Times(d)
	if !v.Valid {
		t.Error("unexpected Valid", v.Valid)
	}
	if v.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", v.ValueOrZero())
	}
}

func assertInt(t *testing.T, i Int, from string) {
	if i.Int64 != 12345 {
		t.Errorf("bad %s int: %d ≠ %d\n", from, i.Int64, 12345)
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullInt(t *testing.T, i Int, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}
