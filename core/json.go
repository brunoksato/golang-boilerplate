package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func ModelToJson(model interface{}) string {
	j, err := json.Marshal(model)
	if err != nil {
		panic(fmt.Sprintf("Error %v encoding JSON for %v", err, model))
	}

	jsonStr := string(j)
	v := reflect.Indirect(reflect.ValueOf(model))
	ot := v.Type()
	t := ot
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	} else if t.Kind() == reflect.Interface {
		t = v.Elem().Type()
	}

	return jsonStr
}

func ModelToJsonMap(modl interface{}) map[string]interface{} {
	jsonStr := ModelToJson(modl)
	m := JsonToMap(jsonStr)
	return m
}

func JsonToMap(jsonStr string) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		panic(fmt.Sprintf("Error %v unmarshaling JSON for %v", err, jsonStr))
	}

	return jsonMap
}

func JsonToMapArray(jsonStr string) []map[string]interface{} {
	var arr []map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &arr)

	if err != nil {
		panic(fmt.Sprintf("Error %v unmarshaling JSON for %v", err, jsonStr))
	}

	return arr
}

func JsonToModel(jsonStr string, item interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), &item)

	if err == nil {
		v := reflect.Indirect(reflect.ValueOf(item))
		ot := v.Type()
		t := ot
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			t = t.Elem()
		} else if t.Kind() == reflect.Interface {
			t = v.Elem().Type()
		}
	}
	return err
}

func SetByJsonTag(item interface{}, jsonKey string, newVal interface{}) DefaultError {
	data := map[string]interface{}{
		"type": reflect.TypeOf(item),
		"key":  jsonKey,
		"val":  newVal,
	}

	if jsonKey == "" || jsonKey == "-" {
		return NewBusinessError("Invalid JSON key", data)
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NewBusinessError("Cannot set value on nil item", data)
		}
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag
		fKey := JsonName(f)
		vField := v.Field(i)
		if fKey == jsonKey {
			if tag.Get("settable") == "false" {
				return NewPermissionError("field is unsettable", data)
			}
			destType := vField.Type()
			if destType.Kind() == reflect.Ptr {
				destType = destType.Elem()
			}
			SetValue(vField.Addr(), destType, newVal)
			return nil
		}
	}

	return NewNotFoundError("field not found", data)
}

func GetFieldByJsonTag(item interface{}, jsonKey string) (field *reflect.StructField, merr DefaultError) {
	data := map[string]interface{}{
		"type": reflect.TypeOf(item),
		"key":  jsonKey,
	}

	if jsonKey == "" || jsonKey == "-" {
		return nil, NewBusinessError("Invalid JSON key", data)
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, NewBusinessError("Cannot set value on nil item", data)
		}
		v = reflect.ValueOf(item).Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fKey := JsonName(f)
		if fKey == jsonKey {
			return &f, nil
		}
	}

	return nil, NewNotFoundError("field not found", data)
}

func JsonName(f reflect.StructField) string {
	tag := f.Tag
	jsonTag := tag.Get("json")
	if jsonTag == "" {
		return ""
	}

	vals := strings.Split(jsonTag, ",")
	return vals[0]
}

func SetValue(v reflect.Value, destType reflect.Type, newVal interface{}) {
	vSet := v.Elem()

	if vSet.Kind() == reflect.Ptr {
		if newVal == nil {
			vSet.Set(reflect.Zero(vSet.Type()))
			return
		} else {
			floatVal, err := InterfaceToFloat64(newVal)
			if err == nil {
				switch destType.Kind() {
				case reflect.Int:
					n := int(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Int32:
					n := int32(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Int64:
					n := int64(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Uint:
					n := uint(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Uint8:
					n := uint8(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Uint16:
					n := uint16(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Uint32:
					n := uint32(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Uint64:
					n := uint64(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Float32:
					n := float32(floatVal)
					vSet.Set(reflect.ValueOf(&n))
				case reflect.Float64:
					vSet.Set(reflect.ValueOf(&newVal))
				case reflect.Struct:
					if destType == reflect.TypeOf(NullableTimestamp{}) {
						ts := NewNullableTimestamp(int64(floatVal))
						vSet.Set(reflect.ValueOf(ts))
					} else if destType == reflect.TypeOf(Timestamp{}) {
						ts := NewTimestamp(int64(floatVal))
						vSet.Set(reflect.ValueOf(ts))
					}
				default:
					vSet.Set(reflect.ValueOf(&floatVal))
				}
			} else {
				if destType.Kind() == reflect.String {
					strVal := newVal.(string)
					vSet.Set(reflect.ValueOf(&strVal))
				} else if destType.Kind() == reflect.Bool {
					boolVal := newVal.(bool)
					vSet.Set(reflect.ValueOf(&boolVal))
				} else {
					vSet.Set(reflect.ValueOf(&newVal))
				}
			}
		}
	} else {
		floatVal, err := InterfaceToFloat64(newVal)
		if err == nil {
			switch destType.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				vSet.SetInt(int64(floatVal))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vSet.SetUint(uint64(floatVal))
			case reflect.Float32, reflect.Float64:
				vSet.SetFloat(floatVal)
			default:
				vSet.Set(reflect.ValueOf(floatVal))
			}
		} else {
			vSet.Set(reflect.ValueOf(newVal))
		}
	}
}

func InterfaceToFloat64(i interface{}) (float64, error) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int:
		return float64(i.(int)), nil
	case reflect.Uint:
		return float64(i.(uint)), nil
	case reflect.Uint64:
		return float64(i.(uint64)), nil
	case reflect.Int32:
		return float64(i.(int32)), nil
	case reflect.Int64:
		return float64(i.(int64)), nil
	case reflect.Float32:
		return float64(i.(float32)), nil
	case reflect.Float64:
		return i.(float64), nil
	default:
		return 0, fmt.Errorf("not implemented for type %v", v.Kind())
	}
}

func IsJsonEnabled(f reflect.StructField, apiType APIType) bool {
	enabled := false

	tag := f.Tag
	jsonField := tag.Get("json")
	if jsonField != "" {
		sensitive := false
		matches := false
		vals := strings.Split(jsonField, ",")
		for j := 1; j < len(vals); j++ {
			val := vals[j]
			if val == "user" || val == "admin" {
				sensitive = true
				matches = matches || val == "user" && apiType == USER_API
				matches = matches || val == "admin" && apiType == ADMIN_API
			}
		}
		enabled = matches || !sensitive
	}

	return enabled
}
