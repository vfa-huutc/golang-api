package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// CensorSensitiveData censors sensitive data in complex data structures
func CensorSensitiveData(data any, maskFields []string) any {
	// Handle nil input
	if data == nil {
		return nil
	}

	// Use reflection to handle more dynamic type checking
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Slice:
		return censorSlice(data, maskFields)
	case reflect.Map:
		return censorMap(data, maskFields)
	case reflect.Struct:
		return censorStruct(data, maskFields)
	case reflect.Ptr:
		// Dereference pointer and recursively censor
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

// censorSlice handles censoring slice types
func censorSlice(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	censoredSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Len())

	for i := range val.Len() {
		item := val.Index(i).Interface()
		censoredItem := CensorSensitiveData(item, maskFields)
		censoredSlice.Index(i).Set(reflect.ValueOf(censoredItem))
	}

	return censoredSlice.Interface()
}

// censorMap handles censoring map types
func censorMap(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	censoredMap := reflect.MakeMap(val.Type())

	iter := val.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// Check if the key (converted to string) is in maskFields
		keyStr := fmt.Sprintf("%v", key.Interface())

		var censoredValue reflect.Value
		if contains(maskFields, keyStr) {
			// Mask the entire value if the key matches
			censoredValue = reflect.ValueOf(maskValue(value.Interface()))
		} else {
			// Recursively censor nested structures
			censoredValue = reflect.ValueOf(CensorSensitiveData(value.Interface(), maskFields))
		}

		censoredMap.SetMapIndex(key, censoredValue)
	}

	return censoredMap.Interface()
}

// censorStruct handles censoring struct types
func censorStruct(data any, maskFields []string) any {
	val := reflect.ValueOf(data)
	typ := val.Type()

	censoredStruct := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if contains(maskFields, fieldType.Name) {
			// Trường cần mask
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					// Nếu con trỏ nil thì giữ nguyên nil
					censoredStruct.Field(i).Set(reflect.Zero(field.Type()))
				} else {
					// Mask giá trị bên trong con trỏ
					maskedVal := maskValue(field.Elem().Interface())
					maskedValReflect := reflect.ValueOf(maskedVal)

					// Tạo con trỏ mới cùng kiểu với trường
					ptr := reflect.New(fieldType.Type.Elem())
					ptr.Elem().Set(maskedValReflect)

					censoredStruct.Field(i).Set(ptr)
				}
			} else {
				// Trường không phải con trỏ thì mask trực tiếp
				censoredStruct.Field(i).Set(reflect.ValueOf(maskValue(field.Interface())))
			}
		} else {
			// Trường không cần mask, đệ quy censor nested
			censoredValue := CensorSensitiveData(field.Interface(), maskFields)
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					censoredStruct.Field(i).Set(reflect.Zero(field.Type()))
				} else {
					ptr := reflect.New(fieldType.Type.Elem())
					ptr.Elem().Set(reflect.ValueOf(censoredValue))
					censoredStruct.Field(i).Set(ptr)
				}
			} else {
				censoredStruct.Field(i).Set(reflect.ValueOf(censoredValue))
			}
		}
	}

	return censoredStruct.Interface()
}

// contains checks if a slice contains a given string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

// maskValue provides advanced masking for different value types
func maskValue(value any) any {
	switch v := value.(type) {
	case string:
		return maskString(v)
	case fmt.Stringer:
		return maskString(v.String())
	case []byte:
		masked := maskString(string(v))
		return []byte(masked)
	case nil:
		return nil
	default:
		return maskReflectedValue(value)
	}
}

// maskString provides sophisticated string masking
func maskString(s string) string {
	// Default masking for other strings
	if len(s) > 2 {
		maskLen := min(len(s)-2, 8)
		return string(s[0]) + strings.Repeat("*", maskLen) + string(s[len(s)-1])
	}
	return strings.Repeat("*", len(s))
}

// maskReflectedValue handles masking for complex types using reflection
func maskReflectedValue(value any) any {
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		// Create a masked slice of the same length
		maskedSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Len())
		for i := range val.Len() {
			maskedSlice.Index(i).Set(reflect.ValueOf("*****"))
		}
		return maskedSlice.Interface()
	case reflect.Struct:
		// Create a struct with all fields masked
		maskedStruct := reflect.New(val.Type()).Elem()
		for i := range val.NumField() {
			maskedStruct.Field(i).Set(reflect.ValueOf("*****"))
		}
		return maskedStruct.Interface()
	default:
		return "*****"
	}
}
