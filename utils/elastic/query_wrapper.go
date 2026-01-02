package elastic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

// 定义查询类型常量，模拟Java注解的查询类型
type QueryType string

const (
	// Eq 等于查询
	Eq QueryType = "EQ"
	// Neq 不等于查询
	Neq QueryType = "NEQ"
	// Gt 大于查询
	Gt QueryType = "GT"
	// Gte 大于等于查询
	Gte QueryType = "GTE"
	// Lt 小于查询
	Lt QueryType = "LT"
	// Lte 小于等于查询
	Lte QueryType = "LTE"
	// Like 模糊查询
	Like QueryType = "LIKE"
	// LikeLeft 左模糊查询
	LikeLeft QueryType = "LIKE_LEFT"
	// LikeRight 右模糊查询
	LikeRight QueryType = "LIKE_RIGHT"
	// NotLike 不包含模糊查询
	NotLike QueryType = "NOT_LIKE"
	// In 在数组中查询
	In QueryType = "IN"
	// NotIn 不在数组中查询
	NotIn QueryType = "NOT_IN"
	// IsNull 为空查询
	IsNull QueryType = "IS_NULL"
	// IsNotNull 不为空查询
	IsNotNull QueryType = "IS_NOT_NULL"
	// Between 范围查询
	Between QueryType = "BETWEEN"
)

// 定义查询标签的键名
const QueryTagKey = "query"

// QueryWrapperImpl 查询条件包装器实现
type QueryWrapperImpl[T any] struct {
	boolQuery     *types.BoolQuery
	sortFields    []types.SortCombinations
	limit         int
	offset        int
	groupByFields []string
	indexName     string
}

// NewQueryWrapper 创建新的查询包装器
func NewQueryWrapper[T any]() QueryWrapper[T] {
	return &QueryWrapperImpl[T]{
		boolQuery: &types.BoolQuery{},
		limit:     -1,
		offset:    0,
	}
}

// FromQueryStruct 从查询结构体创建查询包装器
// 支持的标签格式：`query:"type=EQ"` 或 `query:"type=LIKE,field=name"`
func FromQueryStruct[T any](queryStruct interface{}) QueryWrapper[T] {
	wrapper := NewQueryWrapper[T]()

	// 获取查询结构体的反射值
	v := reflect.ValueOf(queryStruct)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return wrapper
	}

	// 获取查询结构体的类型
	t := v.Type()

	// 遍历结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 跳过未导出的字段
		if !fieldValue.CanInterface() {
			continue
		}

		// 获取字段的值
		value := fieldValue.Interface()

		// 检查值是否为零值，如果是则跳过
		if isZeroValue(value) {
			continue
		}

		// 获取query标签
		queryTag := field.Tag.Get(QueryTagKey)
		if queryTag == "" {
			continue
		}

		// 解析标签内容
		queryConfig := parseQueryTag(queryTag)

		// 获取查询类型
		queryTypeStr, ok := queryConfig["type"]
		if !ok {
			continue
		}
		queryType := QueryType(queryTypeStr)

		// 获取字段名，如果标签中指定了field，则使用指定的字段名，否则使用结构体字段名
		fieldName := field.Name
		if customField, ok := queryConfig["field"]; ok {
			fieldName = customField
		}
		// 转换为蛇形命名法（可选，根据实际需求）
		// fieldName = toSnakeCase(fieldName)

		// 根据查询类型调用相应的方法
		applyQueryType(wrapper, queryType, fieldName, value)
	}

	return wrapper
}

// parseQueryTag 解析查询标签
// 标签格式：`query:"type=EQ,field=name"`
func parseQueryTag(tag string) map[string]string {
	config := make(map[string]string)

	// 分割标签中的键值对
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			config[key] = value
		}
	}

	return config
}

// applyQueryType 根据查询类型应用查询条件
func applyQueryType[T any](wrapper QueryWrapper[T], queryType QueryType, field string, value interface{}) {
	switch queryType {
	case Eq:
		wrapper.Eq(field, value)
	case Neq:
		wrapper.Neq(field, value)
	case Gt:
		wrapper.Gt(field, value)
	case Gte:
		wrapper.Gte(field, value)
	case Lt:
		wrapper.Lt(field, value)
	case Lte:
		wrapper.Lte(field, value)
	case Like:
		wrapper.Like(field, value)
	case LikeLeft:
		// 左模糊查询: %value
		wildcardValue := fmt.Sprintf("%v*", value)
		wrapper.Like(field, wildcardValue)
	case LikeRight:
		// 右模糊查询: value%
		wildcardValue := fmt.Sprintf("*%v", value)
		wrapper.Like(field, wildcardValue)
	case NotLike:
		wrapper.NotLike(field, value)
	case In:
		// 确保value是切片类型
		if sliceValue, ok := value.([]interface{}); ok {
			wrapper.In(field, sliceValue)
		} else if reflect.TypeOf(value).Kind() == reflect.Slice {
			// 转换为[]interface{}
			slice := reflect.ValueOf(value)
			interfaceSlice := make([]interface{}, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				interfaceSlice[i] = slice.Index(i).Interface()
			}
			wrapper.In(field, interfaceSlice)
		}
	case NotIn:
		// 确保value是切片类型
		if sliceValue, ok := value.([]interface{}); ok {
			wrapper.NotIn(field, sliceValue)
		} else if reflect.TypeOf(value).Kind() == reflect.Slice {
			// 转换为[]interface{}
			slice := reflect.ValueOf(value)
			interfaceSlice := make([]interface{}, slice.Len())
			for i := 0; i < slice.Len(); i++ {
				interfaceSlice[i] = slice.Index(i).Interface()
			}
			wrapper.NotIn(field, interfaceSlice)
		}
	case IsNull:
		wrapper.IsNull(field)
	case IsNotNull:
		wrapper.IsNotNull(field)
	// case Between: 范围查询需要特殊处理，需要两个值
	// 这里简化处理，假设value是包含两个元素的切片
	case Between:
		if sliceValue, ok := value.([]interface{}); ok && len(sliceValue) >= 2 {
			wrapper.Between(field, sliceValue[0], sliceValue[1])
		} else if reflect.TypeOf(value).Kind() == reflect.Slice {
			slice := reflect.ValueOf(value)
			if slice.Len() >= 2 {
				min := slice.Index(0).Interface()
				max := slice.Index(1).Interface()
				wrapper.Between(field, min, max)
			}
		}
	}
}

// isZeroValue 判断值是否为零值
func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array, reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// Eq 等于条件
func (w *QueryWrapperImpl[T]) Eq(field string, value interface{}) QueryWrapper[T] {
	query := types.Query{
		Term: map[string]types.TermQuery{
			field: {Value: value},
		},
	}
	w.addMustQuery(query)
	return w
}

// Neq 不等于条件
func (w *QueryWrapperImpl[T]) Neq(field string, value interface{}) QueryWrapper[T] {
	query := types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{
					Term: map[string]types.TermQuery{
						field: {Value: value},
					},
				},
			},
		},
	}
	w.addMustQuery(query)
	return w
}

// Gt 大于条件
func (w *QueryWrapperImpl[T]) Gt(field string, value interface{}) QueryWrapper[T] {
	return w.addRangeQuery(field, map[string]interface{}{"gt": value})
}

// Gte 大于等于条件
func (w *QueryWrapperImpl[T]) Gte(field string, value interface{}) QueryWrapper[T] {
	return w.addRangeQuery(field, map[string]interface{}{"gte": value})
}

// Lt 小于条件
func (w *QueryWrapperImpl[T]) Lt(field string, value interface{}) QueryWrapper[T] {
	return w.addRangeQuery(field, map[string]interface{}{"lt": value})
}

// Lte 小于等于条件
func (w *QueryWrapperImpl[T]) Lte(field string, value interface{}) QueryWrapper[T] {
	return w.addRangeQuery(field, map[string]interface{}{"lte": value})
}

// Like 模糊查询
func (w *QueryWrapperImpl[T]) Like(field string, value interface{}) QueryWrapper[T] {
	wildcardValue := fmt.Sprintf("*%v*", value)
	query := types.Query{
		Wildcard: map[string]types.WildcardQuery{
			field: {Value: &wildcardValue},
		},
	}
	w.addMustQuery(query)
	return w
}

// NotLike 不包含模糊查询
func (w *QueryWrapperImpl[T]) NotLike(field string, value interface{}) QueryWrapper[T] {
	wildcardValue := fmt.Sprintf("*%v*", value)
	query := types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{
					Wildcard: map[string]types.WildcardQuery{
						field: {Value: &wildcardValue},
					},
				},
			},
		},
	}
	w.addMustQuery(query)
	return w
}

// In 在数组中
func (w *QueryWrapperImpl[T]) In(field string, values []interface{}) QueryWrapper[T] {
	termsQuery := types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			field: values,
		},
	}
	query := types.Query{
		Terms: &termsQuery,
	}
	w.addMustQuery(query)
	return w
}

// NotIn 不在数组中
func (w *QueryWrapperImpl[T]) NotIn(field string, values []interface{}) QueryWrapper[T] {
	termsQuery := types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			field: values,
		},
	}
	query := types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{
					Terms: &termsQuery,
				},
			},
		},
	}
	w.addMustQuery(query)
	return w
}

// Between 范围查询
func (w *QueryWrapperImpl[T]) Between(field string, min, max interface{}) QueryWrapper[T] {
	return w.addRangeQuery(field, map[string]interface{}{"gte": min, "lte": max})
}

// NotBetween 不在范围内
func (w *QueryWrapperImpl[T]) NotBetween(field string, min, max interface{}) QueryWrapper[T] {
	// 将参数转换为 json.RawMessage
	minJson, _ := json.Marshal(min)
	maxJson, _ := json.Marshal(max)

	rangeQuery := types.UntypedRangeQuery{
		Gte: minJson,
		Lte: maxJson,
	}

	query := types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{
					Range: map[string]types.RangeQuery{
						field: &rangeQuery,
					},
				},
			},
		},
	}
	w.addMustQuery(query)
	return w
}

// IsNull 为空
func (w *QueryWrapperImpl[T]) IsNull(field string) QueryWrapper[T] {
	existsQuery := types.ExistsQuery{Field: field}
	query := types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{
					Exists: &existsQuery,
				},
			},
		},
	}
	w.addMustQuery(query)
	return w
}

// IsNotNull 不为空
func (w *QueryWrapperImpl[T]) IsNotNull(field string) QueryWrapper[T] {
	existsQuery := types.ExistsQuery{Field: field}
	query := types.Query{
		Exists: &existsQuery,
	}
	w.addMustQuery(query)
	return w
}

// OrderBy 排序
func (w *QueryWrapperImpl[T]) OrderBy(field string, asc bool) QueryWrapper[T] {
	var order *sortorder.SortOrder
	if asc {
		order = &sortorder.Asc
	} else {
		order = &sortorder.Desc
	}

	sortOpt := types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			field: {Order: order},
		},
	}

	w.sortFields = append(w.sortFields, sortOpt)
	return w
}

// GroupBy 分组
func (w *QueryWrapperImpl[T]) GroupBy(fields ...string) QueryWrapper[T] {
	w.groupByFields = append(w.groupByFields, fields...)
	return w
}

// Limit 限制数量
func (w *QueryWrapperImpl[T]) Limit(limit int) QueryWrapper[T] {
	w.limit = limit
	return w
}

// Offset 偏移量
func (w *QueryWrapperImpl[T]) Offset(offset int) QueryWrapper[T] {
	w.offset = offset
	return w
}

// And 与条件
func (w *QueryWrapperImpl[T]) And(conditions ...QueryWrapper[T]) QueryWrapper[T] {
	for _, condition := range conditions {
		if impl, ok := condition.(*QueryWrapperImpl[T]); ok {
			if impl.boolQuery.Must != nil {
				w.boolQuery.Must = append(w.boolQuery.Must, impl.boolQuery.Must...)
			}
			if impl.boolQuery.Filter != nil {
				w.boolQuery.Filter = append(w.boolQuery.Filter, impl.boolQuery.Filter...)
			}
		}
	}
	return w
}

// Or 或条件
func (w *QueryWrapperImpl[T]) Or(conditions ...QueryWrapper[T]) QueryWrapper[T] {
	if w.boolQuery.Should == nil {
		w.boolQuery.Should = []types.Query{}
	}

	for _, condition := range conditions {
		if impl, ok := condition.(*QueryWrapperImpl[T]); ok {
			w.boolQuery.Should = append(w.boolQuery.Should, types.Query{
				Bool: impl.boolQuery,
			})
		}
	}

	w.boolQuery.MinimumShouldMatch = "1"
	return w
}

// BuildQuery 构建查询
func (w *QueryWrapperImpl[T]) BuildQuery() types.Query {
	return types.Query{Bool: w.boolQuery}
}

// BuildSearchRequest 构建搜索请求
func (w *QueryWrapperImpl[T]) BuildSearchRequest() *search.Request {
	query := w.BuildQuery()
	req := &search.Request{
		Query: &query,
	}

	// 设置排序
	if len(w.sortFields) > 0 {
		req.Sort = w.sortFields
	}

	// 设置分页
	if w.limit > 0 {
		req.Size = &w.limit
	}
	if w.offset > 0 {
		req.From = &w.offset
	}

	// 设置分组
	if len(w.groupByFields) > 0 {
		req.Aggregations = make(map[string]types.Aggregations)
		for i, field := range w.groupByFields {
			aggName := fmt.Sprintf("group_by_%s_%d", field, i)
			fieldPtr := field
			req.Aggregations[aggName] = types.Aggregations{
				Terms: &types.TermsAggregation{
					Field: &fieldPtr,
				},
			}
		}
	}

	return req
}

// addMustQuery 添加必须查询条件
func (w *QueryWrapperImpl[T]) addMustQuery(query types.Query) {
	if w.boolQuery.Must == nil {
		w.boolQuery.Must = []types.Query{}
	}
	w.boolQuery.Must = append(w.boolQuery.Must, query)
}

// addRangeQuery 添加范围查询条件
func (w *QueryWrapperImpl[T]) addRangeQuery(field string, params map[string]interface{}) QueryWrapper[T] {
	rangeQuery := types.UntypedRangeQuery{}

	for k, v := range params {
		// 将参数转换为 json.RawMessage
		jsonValue, err := json.Marshal(v)
		if err != nil {
			continue
		}

		switch strings.ToLower(k) {
		case "gt":
			rangeQuery.Gt = jsonValue
		case "gte":
			rangeQuery.Gte = jsonValue
		case "lt":
			rangeQuery.Lt = jsonValue
		case "lte":
			rangeQuery.Lte = jsonValue
		}
	}

	query := types.Query{
		Range: map[string]types.RangeQuery{
			field: &rangeQuery,
		},
	}

	w.addMustQuery(query)
	return w
}
