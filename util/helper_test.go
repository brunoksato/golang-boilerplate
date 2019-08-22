package util

import (
	"reflect"
	"testing"
)

// Helper funcs

type TestStruct struct {
	ID          uint
	Title       string
	Description *string
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("assertEqual: Expected:\n%v\nbut got:\n%v\n (caller: %s)", expected, actual, CallerInfo())
	}
}

func AssertTrue(t *testing.T, actual bool) {
	if !actual {
		t.Errorf("assertTrue: Assertion is false (caller: %s)", CallerInfo())
	}
}

func AssertFalse(t *testing.T, actual bool) {
	if actual {
		t.Errorf("assertFalse: Assertion is true (caller: %s)", CallerInfo())
	}
}

// Tests

func TestAtou(t *testing.T) {
	var nilErr error

	val, err := Atou("1")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(1), val)

	val, err = Atou("0")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(0), val)

	val, err = Atou("1384234")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(1384234), val)

	val, err = Atou("-1")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(18446744073709551615), val)

	val, err = Atou("not a number")
	AssertEqual(t, "strconv.Atoi: parsing \"not a number\": invalid syntax", err.Error())
	AssertEqual(t, uint(0), val)

	val, err = Atou("#234")
	AssertEqual(t, "strconv.Atoi: parsing \"#234\": invalid syntax", err.Error())
	AssertEqual(t, uint(0), val)
}

func TestUtoa(t *testing.T) {
	AssertEqual(t, "1", Utoa(1))
	AssertEqual(t, "0", Utoa(0))
	AssertEqual(t, "1384234", Utoa(1384234))
}

func TestCaptialize(t *testing.T) {
	AssertEqual(t, "Testing", Capitalize("testing"))
	AssertEqual(t, "This is a longer string", Capitalize("this is a longer string"))
	AssertEqual(t, "Already Capitalized", Capitalize("Already Capitalized"))
	AssertEqual(t, "", Capitalize(""))
}

func TestSnakeToCamel(t *testing.T) {
	AssertEqual(t, "Testing", SnakeToCamel("testing"))
	AssertEqual(t, "This is a longer string", SnakeToCamel("this is a longer string"))
	AssertEqual(t, "Already Capitalized", SnakeToCamel("Already Capitalized"))
	AssertEqual(t, "ActualSnake", SnakeToCamel("actual_snake"))
	AssertEqual(t, "ThisIsAMultipleSnake", SnakeToCamel("this_is_a_multiple_snake"))
	AssertEqual(t, "", SnakeToCamel(""))
}

func TestCamelToSnake(t *testing.T) {
	AssertEqual(t, "testing", CamelToSnake("Testing"))
	AssertEqual(t, "testing", CamelToSnake("testing"))
	AssertEqual(t, "actual_camel", CamelToSnake("ActualCamel"))
	AssertEqual(t, "this_is_a_multiple_camel", CamelToSnake("ThisIsAMultipleCamel"))
	AssertEqual(t, "", CamelToSnake(""))
}

func TestRound(t *testing.T) {
	AssertEqual(t, 5, Round(4.5))
	AssertEqual(t, 4, Round(4.49))
	AssertEqual(t, 5, Round(4.7))
	AssertEqual(t, 5, Round(5.1))
	AssertEqual(t, 5, Round(5.0))
	AssertEqual(t, 0, Round(0.49))
	AssertEqual(t, 0, Round(-0.49))
	AssertEqual(t, -1, Round(-1.0))
	AssertEqual(t, -2, Round(-1.5))
	AssertEqual(t, -2, Round(-1.51))
}

func TestRoundFloat(t *testing.T) {
	AssertEqual(t, 5.0, RoundFloat(4.5))
	AssertEqual(t, 4.0, RoundFloat(4.49))
	AssertEqual(t, 5.0, RoundFloat(4.7))
	AssertEqual(t, 5.0, RoundFloat(5.1))
	AssertEqual(t, 5.0, RoundFloat(5.0))
	AssertEqual(t, 0.0, RoundFloat(0.49))
	AssertEqual(t, 0.0, RoundFloat(-0.49))
	AssertEqual(t, -1.0, RoundFloat(-1.0))
	AssertEqual(t, -1.0, RoundFloat(-1.5))
	AssertEqual(t, -2.0, RoundFloat(-1.51))
}

func TestRoundToDecimal(t *testing.T) {
	AssertEqual(t, 5.0, RoundToDecimal(4.5, 0))
	AssertEqual(t, 4.5, RoundToDecimal(4.5, 1))
	AssertEqual(t, 4.5, RoundToDecimal(4.5, 2))
	AssertEqual(t, 4.5, RoundToDecimal(4.49, 1))
	AssertEqual(t, 4.49, RoundToDecimal(4.49, 2))
	AssertEqual(t, 4.45, RoundToDecimal(4.449, 2))
	AssertEqual(t, 12433.2342, RoundToDecimal(12433.2341698323239823, 4))
}

func TestNanoToMicro(t *testing.T) {
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810000))
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810260))
	AssertEqual(t, int64(1452522568079811), NanoToMicro(1452522568079810500))
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810499))
	AssertEqual(t, int64(1452522568079811), NanoToMicro(1452522568079810599))
}

func TestMicroToNano(t *testing.T) {
	AssertEqual(t, int64(1452522568079810000), MicroToNano(1452522568079810))
	AssertEqual(t, int64(1452522568079811000), MicroToNano(1452522568079811))
}

func TestIsZeroStruct(t *testing.T) {
	AssertEqual(t, false, IsZeroStruct(nil))
	AssertEqual(t, true, IsZeroStruct(TestStruct{}))
	AssertEqual(t, true, IsZeroStruct(TestStruct{ID: 0}))
	AssertEqual(t, true, IsZeroStruct(TestStruct{ID: 0, Title: ""}))
	AssertEqual(t, false, IsZeroStruct(TestStruct{ID: 10}))
	AssertEqual(t, false, IsZeroStruct(TestStruct{Title: "Hello"}))
	s := "Hello, Test"
	AssertEqual(t, false, IsZeroStruct(TestStruct{Description: &s}))
	AssertEqual(t, false, IsZeroStruct(&TestStruct{}))
	var i interface{}
	AssertEqual(t, false, IsZeroStruct(i))
}

func TestIsEmptyValue(t *testing.T) {
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf(0)))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf(10)))
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf("")))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf("Yo mama")))
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf(uint(0))))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf(uint(10))))
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf(float64(0))))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf(float64(10))))
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf([]uint{})))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf([]uint{1, 5, 7})))
	var a, b *int
	bval := 10
	b = &bval
	AssertEqual(t, true, IsEmptyValue(reflect.ValueOf(a)))
	AssertEqual(t, false, IsEmptyValue(reflect.ValueOf(b)))
}

func TestPtrEquals(t *testing.T) {
	AssertFalse(t, PtrEquals(1, 1))
	AssertFalse(t, PtrEquals(1, 2))

	var s1, s2 string
	AssertFalse(t, PtrEquals(s1, s2))
	s1 = "Hello"
	s2 = "World"
	AssertFalse(t, PtrEquals(s1, s2))
	s2 = s1
	AssertFalse(t, PtrEquals(s1, s2))

	var sPtr1, sPtr2 *string
	AssertTrue(t, PtrEquals(sPtr1, sPtr2))
	sPtr1 = &s1
	AssertFalse(t, PtrEquals(sPtr1, sPtr2))
	sPtr1 = nil
	sPtr2 = &s2
	AssertFalse(t, PtrEquals(sPtr1, sPtr2))
	s2 = "World"
	sPtr1 = &s1
	sPtr2 = &s2
	AssertFalse(t, PtrEquals(sPtr1, sPtr2))
	sPtr2 = &s1
	AssertTrue(t, PtrEquals(sPtr1, sPtr2))
	s2 = "Hello"
	sPtr1 = &s1
	sPtr2 = &s2
	AssertTrue(t, PtrEquals(sPtr1, sPtr2))

	var uPtr1, uPtr2 *uint
	u1 := uint(123)
	u2 := uint(345)
	AssertTrue(t, PtrEquals(uPtr1, uPtr2))
	uPtr1 = &u1
	AssertFalse(t, PtrEquals(uPtr1, uPtr2))
	uPtr1 = nil
	uPtr2 = &u2
	AssertFalse(t, PtrEquals(uPtr1, uPtr2))
	uPtr1 = &u1
	uPtr2 = &u2
	AssertFalse(t, PtrEquals(uPtr1, uPtr2))
	uPtr2 = &u1
	AssertTrue(t, PtrEquals(uPtr1, uPtr2))
	u2 = uint(123)
	uPtr1 = &u1
	uPtr2 = &u2
	AssertTrue(t, PtrEquals(uPtr1, uPtr2))
}

func TestUnique(t *testing.T) {
	integers := []uint{1, 2, 3, 3, 4, 5}
	unique := Unique(integers)
	expected := []uint{1, 2, 3, 4, 5}
	AssertEqual(t, unique, expected)
}
