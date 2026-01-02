package elastic

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Repository ES 仓库接口，定义 CRUD 操作
type Repository[T any] interface {
	// Insert 插入单个文档
	Insert(ctx context.Context, entity *T) (string, error)

	// InsertBatch 批量插入文档
	InsertBatch(ctx context.Context, entities []*T) ([]string, error)

	// Update 更新文档
	Update(ctx context.Context, entity *T) error

	// UpdateById 根据 ID 更新文档
	UpdateById(ctx context.Context, id string, update map[string]interface{}) error

	// Delete 删除文档
	Delete(ctx context.Context, entity *T) error

	// DeleteById 根据 ID 删除文档
	DeleteById(ctx context.Context, id string) error

	// DeleteBatch 批量删除文档
	DeleteBatch(ctx context.Context, ids []string) error

	// GetById 根据 ID 获取文档
	GetById(ctx context.Context, id string) (*T, error)

	// GetOne 获取单个文档（根据条件）
	GetOne(ctx context.Context, wrapper QueryWrapper[T]) (*T, error)

	// List 获取文档列表
	List(ctx context.Context, wrapper QueryWrapper[T]) ([]*T, error)

	// Page 分页查询
	Page(ctx context.Context, wrapper QueryWrapper[T], page, size int) (*PageResult[T], error)

	// Count 统计文档数量
	Count(ctx context.Context, wrapper QueryWrapper[T]) (int64, error)

	// Exists 检查文档是否存在
	Exists(ctx context.Context, id string) (bool, error)

	// Search 原生搜索
	Search(ctx context.Context, index string, req *search.Request) (*search.Response, error)
}

// PageResult 分页结果
type PageResult[T any] struct {
	Total    int64  `json:"total"`
	Pages    int    `json:"pages"`
	Current  int    `json:"current"`
	Size     int    `json:"size"`
	Records  []*T   `json:"records"`
	Index    string `json:"index"`
}

// QueryWrapper 查询条件包装器
type QueryWrapper[T any] interface {
	// Eq 等于条件
	Eq(field string, value interface{}) QueryWrapper[T]

	// Neq 不等于条件
	Neq(field string, value interface{}) QueryWrapper[T]

	// Gt 大于条件
	Gt(field string, value interface{}) QueryWrapper[T]

	// Gte 大于等于条件
	Gte(field string, value interface{}) QueryWrapper[T]

	// Lt 小于条件
	Lt(field string, value interface{}) QueryWrapper[T]

	// Lte 小于等于条件
	Lte(field string, value interface{}) QueryWrapper[T]

	// Like 模糊查询
	Like(field string, value interface{}) QueryWrapper[T]

	// NotLike 不包含模糊查询
	NotLike(field string, value interface{}) QueryWrapper[T]

	// In 在数组中
	In(field string, values []interface{}) QueryWrapper[T]

	// NotIn 不在数组中
	NotIn(field string, values []interface{}) QueryWrapper[T]

	// Between 范围查询
	Between(field string, min, max interface{}) QueryWrapper[T]

	// NotBetween 不在范围内
	NotBetween(field string, min, max interface{}) QueryWrapper[T]

	// IsNull 为空
	IsNull(field string) QueryWrapper[T]

	// IsNotNull 不为空
	IsNotNull(field string) QueryWrapper[T]

	// OrderBy 排序
	OrderBy(field string, asc bool) QueryWrapper[T]

	// GroupBy 分组
	GroupBy(fields ...string) QueryWrapper[T]

	// Limit 限制数量
	Limit(limit int) QueryWrapper[T]

	// Offset 偏移量
	Offset(offset int) QueryWrapper[T]

	// And 与条件
	And(conditions ...QueryWrapper[T]) QueryWrapper[T]

	// Or 或条件
	Or(conditions ...QueryWrapper[T]) QueryWrapper[T]

	// BuildQuery 构建查询
	BuildQuery() types.Query

	// BuildSearchRequest 构建搜索请求
	BuildSearchRequest() *search.Request
}
