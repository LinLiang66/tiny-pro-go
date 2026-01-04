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

// IndexConfig 索引配置信息
type IndexConfig struct {
	BaseIndexName    string
	PartitionType    string // 分区类型："monthly"（按月）、"yearly"（按年）、"none"（不分）
	TimeSuffixFormat string // 时间后缀格式，如 "200601"（按月）、"2006"（按年）
}

// EsIndex 定义 ES 索引信息
type EsIndex struct {
	IndexName string
	Type      string
	Shards    int
	Replicas  int
}

// parseIndexTag 解析索引标签，支持格式："index_name,partition=monthly/yearly/none"
func parseIndexTag(tag string) IndexConfig {
	config := IndexConfig{
		PartitionType:    "", // 默认使用全局配置
		TimeSuffixFormat: "",
	}

	// 拆分标签中的各个部分
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		config.BaseIndexName = parts[0]
	}

	// 解析分区策略
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if strings.HasPrefix(part, "partition=") {
			partitionType := strings.TrimPrefix(part, "partition=")
			config.PartitionType = partitionType

			// 根据分区类型设置时间后缀格式
			switch partitionType {
			case "monthly":
				config.TimeSuffixFormat = "200601"
			case "yearly":
				config.TimeSuffixFormat = "2006"
			default:
				config.PartitionType = "none"
				config.TimeSuffixFormat = ""
			}
			break
		}
	}

	return config
}

// GetBaseIndexName 获取实体对应的基础索引名（不带时间后缀）
func GetBaseIndexName(model interface{}) string {
	return GetIndexConfig(model).BaseIndexName
}

// GetIndexConfig 获取实体对应的索引配置
func GetIndexConfig(model interface{}) IndexConfig {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 默认配置
	defaultIndexName := strings.ToLower(t.Name())
	config := IndexConfig{
		BaseIndexName:    defaultIndexName,
		PartitionType:    "monthly", // 默认按月分索引
		TimeSuffixFormat: "200601",
	}

	// 检查是否有 es 标签指定索引配置
	if field, ok := t.FieldByName("Index"); ok {
		if tag := field.Tag.Get("es"); tag != "" && tag != "_index" {
			// 解析标签获取配置
			tagConfig := parseIndexTag(tag)
			// 如果标签中指定了基础索引名，则使用它
			if tagConfig.BaseIndexName != "" {
				config.BaseIndexName = tagConfig.BaseIndexName
			}
			// 如果标签中指定了分区策略，则使用它
			if tagConfig.PartitionType != "" {
				config.PartitionType = tagConfig.PartitionType
				config.TimeSuffixFormat = tagConfig.TimeSuffixFormat
			}
		}
	}

	return config
}

// GetIndexNameWithTimeSuffix 获取带时间后缀的索引名
func GetIndexNameWithTimeSuffix(baseIndexName string, timeSuffixFormat string) string {
	if timeSuffixFormat == "" {
		timeSuffixFormat = "200601" // 默认使用年月格式
	}
	timeSuffix := time.Now().Format(timeSuffixFormat)
	return fmt.Sprintf("%s_%s", baseIndexName, timeSuffix)
}

// GetIndexName 获取实体对应的索引名，根据配置决定是否带时间后缀
func GetIndexName(model interface{}) string {
	config := GetIndexConfig(model)

	// 如果不分区，则直接返回基础索引名
	if config.PartitionType == "none" {
		return config.BaseIndexName
	}

	// 否则返回带时间后缀的索引名
	return GetIndexNameWithTimeSuffix(config.BaseIndexName, config.TimeSuffixFormat)
}

// GetIndexPattern 获取实体对应的索引模式（通配符模式），用于跨索引查询
func GetIndexPattern(model interface{}) string {
	config := GetIndexConfig(model)

	// 如果不分区，则直接返回基础索引名
	if config.PartitionType == "none" {
		return config.BaseIndexName
	}

	// 否则返回带通配符的索引模式
	return fmt.Sprintf("%s_*", config.BaseIndexName)
}

// GetIndexNamesByDateRange 根据日期范围生成需要查询的索引名列表
func GetIndexNamesByDateRange(model interface{}, startTime, endTime time.Time, timeSuffixFormat string) []string {
	config := GetIndexConfig(model)

	// 如果不分区，则直接返回基础索引名
	if config.PartitionType == "none" {
		return []string{config.BaseIndexName}
	}

	// 如果没有指定时间后缀格式，则使用配置中的格式
	if timeSuffixFormat == "" {
		timeSuffixFormat = config.TimeSuffixFormat
	}

	baseIndexName := config.BaseIndexName
	indexNames := make([]string, 0)

	// 确保startTime <= endTime
	if startTime.After(endTime) {
		startTime, endTime = endTime, startTime
	}

	// 判断时间格式是否只包含年份
	isYearOnlyFormat := timeSuffixFormat == "2006"

	// 初始化当前时间
	var currentTime, adjustedEndTime time.Time
	if isYearOnlyFormat {
		// 按年分索引的情况
		currentTime = time.Date(startTime.Year(), 1, 1, 0, 0, 0, 0, startTime.Location())
		adjustedEndTime = time.Date(endTime.Year(), 1, 1, 0, 0, 0, 0, endTime.Location())
	} else {
		// 按月分索引的情况
		currentTime = time.Date(startTime.Year(), startTime.Month(), 1, 0, 0, 0, 0, startTime.Location())
		adjustedEndTime = time.Date(endTime.Year(), endTime.Month(), 1, 0, 0, 0, 0, endTime.Location())
	}

	for {
		timeSuffix := currentTime.Format(timeSuffixFormat)
		indexName := fmt.Sprintf("%s_%s", baseIndexName, timeSuffix)

		// 只添加唯一的索引名
		if len(indexNames) == 0 || indexNames[len(indexNames)-1] != indexName {
			indexNames = append(indexNames, indexName)
		}

		// 到达结束时间，退出循环
		if isYearOnlyFormat {
			if currentTime.Year() == adjustedEndTime.Year() {
				break
			}
			// 增加一年
			currentTime = currentTime.AddDate(1, 0, 0)
		} else {
			if currentTime.Year() == adjustedEndTime.Year() && currentTime.Month() == adjustedEndTime.Month() {
				break
			}
			// 增加一个月
			currentTime = currentTime.AddDate(0, 1, 0)
		}
	}

	return indexNames
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
