package elastic

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// BaseModel ES 基础实体模型
type BaseModel struct {
	ID        string     `json:"id,omitempty" es:"_id"`
	Index     string     `json:"_index,omitempty" es:"_index"`
	Score     float64    `json:"_score,omitempty" es:"_score"`
	CreatedAt time.Time  `json:"created_at,omitempty" es:"created_at"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" es:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" es:"deleted_at"`
}

// EsIndex 定义 ES 索引信息
type EsIndex struct {
	IndexName string
	Type      string
	Shards    int
	Replicas  int
}

// GetIndexName 获取实体对应的索引名
func GetIndexName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 默认使用类型名的小写形式作为索引名
	indexName := strings.ToLower(t.Name())

	// 检查是否有 es 标签指定索引名
	if field, ok := t.FieldByName("Index"); ok {
		if tag := field.Tag.Get("es"); tag != "" && tag != "_index" {
			indexName = tag
		}
	}

	return indexName
}

// GetDocumentID 获取实体的文档 ID
func GetDocumentID(model interface{}) string {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 检查是否有 ID 字段
	if field := v.FieldByName("ID"); field.IsValid() {
		idStr, ok := convertToIDString(field.Interface())
		if ok && idStr != "" {
			return idStr
		}
	}

	// 检查是否有 es 标签为 _id 的字段
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("es"); tag == "_id" {
			if idField := v.Field(i); idField.IsValid() {
				idStr, ok := convertToIDString(idField.Interface())
				if ok && idStr != "" {
					return idStr
				}
			}
			break
		}
	}

	return ""
}

// convertToIDString 将各种类型的 ID 转换为字符串
func convertToIDString(id interface{}) (string, bool) {
	switch v := id.(type) {
	case string:
		return v, true
	case int:
		return fmt.Sprintf("%d", v), true
	case int64:
		return fmt.Sprintf("%d", v), true
	case int32:
		return fmt.Sprintf("%d", v), true
	case uint:
		return fmt.Sprintf("%d", v), true
	case uint64:
		return fmt.Sprintf("%d", v), true
	case uint32:
		return fmt.Sprintf("%d", v), true
	case *string:
		if v != nil {
			return *v, true
		}
	case *int:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	case *int64:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	case *int32:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	case *uint:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	case *uint64:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	case *uint32:
		if v != nil {
			return fmt.Sprintf("%d", *v), true
		}
	}
	return "", false
}

// SetDocumentID 设置实体的文档 ID
func SetDocumentID(model interface{}, id string) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return
	}

	v = v.Elem()

	// 设置ID字段的函数
	setID := func(field reflect.Value) bool {
		if !field.IsValid() || !field.CanSet() {
			return false
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(id)
			return true
		case reflect.Int, reflect.Int64:
			if idInt, err := strconv.ParseInt(id, 10, 64); err == nil {
				field.SetInt(idInt)
				return true
			}
		case reflect.Uint, reflect.Uint64:
			if idUint, err := strconv.ParseUint(id, 10, 64); err == nil {
				field.SetUint(idUint)
				return true
			}
		}
		return false
	}

	// 检查是否有 ID 字段
	if setID(v.FieldByName("ID")) {
		return
	}

	// 检查是否有 es 标签为 _id 的字段
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("es"); tag == "_id" {
			if setID(v.Field(i)) {
				break
			}
		}
	}
}
