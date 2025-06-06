package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// CensorSensitiveData censors sensitive data in complex data structures recursively.
func CensorSensitiveData(data any, maskFields []string) any {
	if data == nil {
		return nil
	}

	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		return censorSlice(data, maskFields)
	case reflect.Map:
		return censorMap(data, maskFields)
	case reflect.Struct:
		return censorStruct(data, maskFields)
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		return CensorSensitiveData(val.Elem().Interface(), maskFields)
	case reflect.String:
		return data
	default:
		return data
	}
}

// censorSlice recursively censors each element in a slice/array.
func censorSlice(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	censoredSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Len())

	for i := 0; i < val.Len(); i++ { // fix vòng lặp đúng cách
		item := val.Index(i).Interface()
		censoredItem := CensorSensitiveData(item, maskFields)
		censoredSlice.Index(i).Set(reflect.ValueOf(censoredItem))
	}

	return censoredSlice.Interface()
}

// censorMap recursively censors map entries based on keys.
func censorMap(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	censoredMap := reflect.MakeMap(val.Type())

	iter := val.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		keyStr := fmt.Sprintf("%v", key.Interface())

		var censoredValue reflect.Value
		if contains(maskFields, keyStr) {
			// Mask toàn bộ giá trị nếu key nhạy cảm
			censoredValue = reflect.ValueOf(maskValue(value.Interface()))
		} else {
			censoredValue = reflect.ValueOf(CensorSensitiveData(value.Interface(), maskFields))
		}

		censoredMap.SetMapIndex(key, censoredValue)
	}

	return censoredMap.Interface()
}

// censorStruct recursively censors struct fields based on field names.
func censorStruct(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	typ := val.Type()
	censoredStruct := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if contains(maskFields, fieldType.Name) {
			// Field need to be masked
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					censoredStruct.Field(i).Set(reflect.Zero(field.Type()))
				} else {
					maskedVal := maskValue(field.Elem().Interface())
					maskedValReflect := reflect.ValueOf(maskedVal)

					ptr := reflect.New(fieldType.Type.Elem())
					ptr.Elem().Set(matchedValOrZero(maskedValReflect, fieldType.Type.Elem()))
					censoredStruct.Field(i).Set(ptr)
				}
			} else {
				censoredStruct.Field(i).Set(matchedValOrZero(reflect.ValueOf(maskValue(field.Interface())), fieldType.Type))
			}
		} else {
			// Field does not need to be masked, process recursively
			censoredValue := CensorSensitiveData(field.Interface(), maskFields)
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					censoredStruct.Field(i).Set(reflect.Zero(field.Type()))
				} else {
					ptr := reflect.New(fieldType.Type.Elem())
					ptr.Elem().Set(matchedValOrZero(reflect.ValueOf(censoredValue), fieldType.Type.Elem()))
					censoredStruct.Field(i).Set(ptr)
				}
			} else {
				censoredStruct.Field(i).Set(matchedValOrZero(reflect.ValueOf(censoredValue), fieldType.Type))
			}
		}
	}

	return censoredStruct.Interface()
}

// matchedValOrZero tries set value if compatible, else zero value (tránh panic)
func matchedValOrZero(val reflect.Value, typ reflect.Type) reflect.Value {
	if val.Type().AssignableTo(typ) {
		return val
	}
	return reflect.Zero(typ)
}

// contains checks if a string is in a slice, case-insensitive
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

// maskValue masks sensitive values based on their type.
func maskValue(value any) any {
	switch v := value.(type) {
	case string:
		return maskString(v)
	case fmt.Stringer:
		return maskString(v.String())
	case []byte:
		return []byte(maskString(string(v)))
	case nil:
		return nil
	default:
		return maskReflectedValue(value)
	}
}

// maskString masks a string by replacing its middle characters with asterisks.
func maskString(s string) string {
	if len(s) > 2 {
		maskLen := min(len(s)-2, 8)
		return string(s[0]) + strings.Repeat("*", maskLen) + string(s[len(s)-1])
	}
	return strings.Repeat("*", len(s))
}

func maskReflectedValue(value any) any {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		maskedSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Len())
		for i := range val.Len() {
			maskedSlice.Index(i).Set(reflect.ValueOf("*****"))
		}
		return maskedSlice.Interface()
	case reflect.Struct:
		maskedStruct := reflect.New(val.Type()).Elem()
		for i := 0; i < val.NumField(); i++ {
			field := maskedStruct.Field(i)
			switch field.Kind() {
			case reflect.String:
				field.SetString("*****")
			case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
				field.SetInt(0)
			case reflect.Bool:
				field.SetBool(false)
			default:
				field.Set(reflect.Zero(field.Type()))
			}
		}
		return maskedStruct.Interface()
	default:
		return "*****"
	}
}

// Min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
