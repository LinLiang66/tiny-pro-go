package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

// Marshal 自定义JSON序列化方法
func Marshal(v interface{}) ([]byte, error) {
	return customMarshal(v)
}

// customMarshal 递归处理自定义序列化
func customMarshal(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)

	// 处理指针
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return []byte("null"), nil
		}
		return customMarshal(val.Elem().Interface())
	}

	// 处理切片和数组
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		if val.Len() == 0 {
			return []byte("[]"), nil
		}

		var buf bytes.Buffer
		buf.WriteByte('[')

		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			bytes, err := customMarshal(item)
			if err != nil {
				return nil, err
			}
			buf.Write(bytes)
			if i < val.Len()-1 {
				buf.WriteByte(',')
			}
		}

		buf.WriteByte(']')
		return buf.Bytes(), nil
	}

	// 处理结构体
	if val.Kind() == reflect.Struct {
		// 处理时间类型
		if val.Type() == reflect.TypeOf(time.Time{}) {
			t := val.Interface().(time.Time)
			return []byte(`"` + formatTime(t) + `"`), nil
		}

		// 处理普通结构体
		var buf bytes.Buffer
		buf.WriteByte('{')

		typ := val.Type()
		first := true

		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			// 解析json tag
			name := field.Name
			if jsonTag != "" {
				// 处理json tag中的选项，如omitempty
				if idx := len(jsonTag); idx > 0 && jsonTag[idx-1] == '}' {
					if start := bytes.Index([]byte(jsonTag), []byte("omitempty")); start > 0 {
						name = jsonTag[:start-1]
					} else {
						name = jsonTag
					}
				} else {
					name = jsonTag
				}
			}

			fieldVal := val.Field(i)
			if !fieldVal.CanInterface() {
				continue
			}

			value, err := customMarshal(fieldVal.Interface())
			if err != nil {
				return nil, err
			}

			if !first {
				buf.WriteByte(',')
			}
			buf.WriteString(`"` + name + `":`)
			buf.Write(value)
			first = false
		}

		buf.WriteByte('}')
		return buf.Bytes(), nil
	}

	// 处理整数类型
	if val.Kind() == reflect.Int || val.Kind() == reflect.Int64 {
		intVal := val.Int()
		// JavaScript安全整数范围：-(2^53-1) 到 2^53-1
		if intVal < -9007199254740991 || intVal > 9007199254740991 {
			// 超出安全范围，序列化为字符串
			return []byte(`"` + strconv.FormatInt(intVal, 10) + `"`), nil
		}
		return []byte(strconv.FormatInt(intVal, 10)), nil
	}

	if val.Kind() == reflect.Uint || val.Kind() == reflect.Uint64 {
		uintVal := val.Uint()
		if uintVal > 9007199254740991 {
			// 超出安全范围，序列化为字符串
			return []byte(`"` + strconv.FormatUint(uintVal, 10) + `"`), nil
		}
		return []byte(strconv.FormatUint(uintVal, 10)), nil
	}

	// 处理浮点数类型（JSON反序列化为interface{}时，所有数字都会变成float64）
	if val.Kind() == reflect.Float64 {
		floatVal := val.Float()
		// 检查是否是整数
		if floatVal == float64(int64(floatVal)) {
			// 转换为int64
			intVal := int64(floatVal)
			// 检查是否超出JavaScript安全整数范围
			if intVal < -9007199254740991 || intVal > 9007199254740991 {
				// 超出安全范围，序列化为字符串
				return []byte(`"` + strconv.FormatInt(intVal, 10) + `"`), nil
			}
			// 作为普通整数处理
			return []byte(strconv.FormatInt(intVal, 10)), nil
		} else if floatVal == float64(uint64(floatVal)) {
			// 转换为uint64
			uintVal := uint64(floatVal)
			// 检查是否超出JavaScript安全整数范围
			if uintVal > 9007199254740991 {
				// 超出安全范围，序列化为字符串
				return []byte(`"` + strconv.FormatUint(uintVal, 10) + `"`), nil
			}
			// 作为普通整数处理
			return []byte(strconv.FormatUint(uintVal, 10)), nil
		}
		// 普通浮点数，使用默认序列化
		return []byte(strconv.FormatFloat(floatVal, 'f', -1, 64)), nil
	}

	// 处理map类型
	if val.Kind() == reflect.Map {
		var buf bytes.Buffer
		buf.WriteByte('{')

		iter := val.MapRange()
		first := true

		for iter.Next() {
			key := iter.Key()
			value := iter.Value()

			// 处理key为string类型
			keyStr, ok := key.Interface().(string)
			if !ok {
				continue
			}

			// 处理value
			valBytes, err := customMarshal(value.Interface())
			if err != nil {
				return nil, err
			}

			if !first {
				buf.WriteByte(',')
			}
			buf.WriteString(`"` + keyStr + `":`)
			buf.Write(valBytes)
			first = false
		}

		buf.WriteByte('}')
		return buf.Bytes(), nil
	}

	// 其他类型使用默认序列化
	return json.Marshal(v)
}

// formatTime 格式化时间为指定格式
func formatTime(t time.Time) string {
	// 如果时间只有日期部分（时分秒都为0），格式化为2025-12-10
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format("2006-01-02")
	}
	// 否则格式化为2025-12-10 16:55:10
	return t.Format("2006-01-02 15:04:05")
}
