package util

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func NewSliceForType(t reflect.Type) interface{} {
	items := reflect.New(reflect.SliceOf(t)).Interface()
	return items
}

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func IntPtr(i int) *int {
	return &i
}

func Float32Ptr(f float32) *float32 {
	return &f
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func BoolPtr(b bool) *bool {
	return &b
}

func AUintToAInt(au []uint) []int {
	ai := make([]int, len(au))
	for i := 0; i < len(au); i++ {
		ai[i] = int(au[i])
	}
	return ai
}

func Atou(a string) (uint, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func Utoa(u uint) string {
	return fmt.Sprintf("%d", u)
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}

// SnakeToCamel returns a string converted from snake case to uppercase
func SnakeToCamel(s string) string {
	if s == "" {
		return ""
	}
	var result string

	words := strings.Split(s, "_")

	for _, word := range words {
		w := []rune(word)
		w[0] = unicode.ToUpper(w[0])
		result += string(w)
	}

	return result
}

func NanoToMicro(n int64) int64 {
	// Conversion to float64 was leading to imperfect results!
	//return int64(Round(float64(n) / 1000.0))
	return (n + 500) / 1000
}

func MicroToNano(ms int64) int64 {
	return ms * 1000
}

func DurationInMs(start, end time.Time) int64 {
	return (end.UnixNano() - start.UnixNano()) / 1000000
}

func IsZeroStruct(item interface{}) bool {
	zero := false

	if item != nil {
		v := reflect.ValueOf(item)
		t := v.Type()
		k := t.Kind()

		if k == reflect.Struct {
			z := reflect.Zero(t)
			if reflect.DeepEqual(v.Interface(), z.Interface()) {
				zero = true
			}
		}
	}

	return zero
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func RoundToDecimal(f float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	res := f * factor
	res = RoundFloat(res)
	res = res / factor
	return res
}

func RoundFloat(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return RoundFloat(f*shift) / shift
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func CallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func ParentCallerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func ItemsOrEmptySlice(t reflect.Type, items interface{}) interface{} {
	if reflect.ValueOf(items).IsNil() {
		items = reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	}
	return items
}

func ConvertQueryTermToOrderTerm(term string) (fields []string, ascending []bool) {
	if strings.HasPrefix(term, "[") && strings.HasSuffix(term, "]") {
		term = term[1 : len(term)-1]
	}

	if term == "" {
		return
	}

	fieldConfigs := strings.Split(term, ",")

	for _, config := range fieldConfigs {
		if strings.HasSuffix(config, "-asc") {
			fields = append(fields, config[0:len(config)-4])
			ascending = append(ascending, true)
		} else if strings.HasSuffix(config, "-desc") {
			fields = append(fields, config[0:len(config)-5])
			ascending = append(ascending, false)
		} else {
			fields = append(fields, config)
			ascending = append(ascending, true)
		}
	}

	fmt.Println(ascending)

	return
}

func PtrEquals(ptr1 interface{}, ptr2 interface{}) bool {
	if reflect.TypeOf(ptr1).Kind() != reflect.Ptr {
		return false
	}
	if reflect.TypeOf(ptr2).Kind() != reflect.Ptr {
		return false
	}
	vptr1 := reflect.ValueOf(ptr1)
	vptr2 := reflect.ValueOf(ptr2)

	if vptr1.IsNil() {
		if vptr2.IsNil() {
			return true
		}
		return false
	} else if vptr2.IsNil() {
		return false
	}

	v1 := reflect.Indirect(vptr1)
	v2 := reflect.Indirect(vptr2)

	return v1.Interface() == v2.Interface()
}

func Unique(input []uint) []uint {
	u := make([]uint, 0, len(input))
	m := make(map[uint]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
