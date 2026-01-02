package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// BaseRepository ES 基础仓库实现
type BaseRepository[T any] struct {
	client    *elasticsearch.TypedClient
	indexName string
}

// NewBaseRepository 创建新的基础仓库
func NewBaseRepository[T any]() (*BaseRepository[T], error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}

	// 创建一个临时实例来获取索引名
	var model T
	indexName := GetIndexName(model)

	return &BaseRepository[T]{
		client:    client,
		indexName: indexName,
	}, nil
}

// Insert 插入单个文档
func (r *BaseRepository[T]) Insert(ctx context.Context, entity *T) (string, error) {
	// 设置创建时间和更新时间
	r.setTimestamps(entity)

	// 获取文档 ID
	id := GetDocumentID(entity)

	// 执行索引请求
	var resp *index.Response
	var err error
	if id != "" {
		resp, err = r.client.Index(r.indexName).Id(id).Document(entity).Do(ctx)
	} else {
		resp, err = r.client.Index(r.indexName).Document(entity).Do(ctx)
	}

	if err != nil {
		return "", err
	}

	// 设置文档 ID 到实体
	SetDocumentID(entity, resp.Id_)
	return resp.Id_, nil
}

// InsertBatch 批量插入文档
func (r *BaseRepository[T]) InsertBatch(ctx context.Context, entities []*T) ([]string, error) {
	// 准备批量请求
	bulkReq := r.client.Bulk()
	ids := make([]string, 0, len(entities))

	for _, entity := range entities {
		// 设置创建时间和更新时间
		r.setTimestamps(entity)

		// 获取文档 ID
		id := GetDocumentID(entity)

		// 构建索引请求
		indexOp := types.IndexOperation{
			Index_: &r.indexName,
		}
		if id != "" {
			indexOp.Id_ = &id
		}

		// 添加到批量请求
		if err := bulkReq.IndexOp(indexOp, entity); err != nil {
			return nil, err
		}
	}

	// 执行批量请求
	resp, err := bulkReq.Do(ctx)
	if err != nil {
		return nil, err
	}

	// 处理响应
	if resp.Errors {
		return nil, fmt.Errorf("bulk insert failed: %v", resp.Items)
	}

	// 获取生成的 ID
	for i, item := range resp.Items {
		for _, result := range item {
			if result.Id_ != nil {
				ids = append(ids, *result.Id_)
				// 设置文档 ID 到实体
				SetDocumentID(entities[i], *result.Id_)
			}
			break
		}
	}

	return ids, nil
}

// Update 更新文档
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	// 设置更新时间
	r.setUpdateTimestamp(entity)

	// 获取文档 ID
	id := GetDocumentID(entity)
	if id == "" {
		return fmt.Errorf("document ID is required for update")
	}

	// 执行更新请求
	_, err := r.client.Update(r.indexName, id).Doc(entity).Do(ctx)
	return err
}

// UpdateById 根据 ID 更新文档
func (r *BaseRepository[T]) UpdateById(ctx context.Context, id string, update map[string]interface{}) error {
	// 添加更新时间
	update["updated_at"] = time.Now().UTC()

	// 执行更新请求
	_, err := r.client.Update(r.indexName, id).Doc(update).Do(ctx)
	return err
}

// Delete 删除文档
func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T) error {
	// 获取文档 ID
	id := GetDocumentID(entity)
	if id == "" {
		return fmt.Errorf("document ID is required for delete")
	}

	// 执行删除请求
	_, err := r.client.Delete(r.indexName, id).Do(ctx)
	return err
}

// DeleteById 根据 ID 删除文档
func (r *BaseRepository[T]) DeleteById(ctx context.Context, id string) error {
	// 执行删除请求
	_, err := r.client.Delete(r.indexName, id).Do(ctx)
	return err
}

// DeleteBatch 批量删除文档
func (r *BaseRepository[T]) DeleteBatch(ctx context.Context, ids []string) error {
	// 准备批量请求
	bulkReq := r.client.Bulk()

	for _, id := range ids {
		// 构建删除请求
		deleteOp := types.DeleteOperation{
			Index_: &r.indexName,
			Id_:    &id,
		}

		// 添加到批量请求
		if err := bulkReq.DeleteOp(deleteOp); err != nil {
			return err
		}
	}

	// 执行批量请求
	resp, err := bulkReq.Do(ctx)
	if err != nil {
		return err
	}

	// 处理响应
	if resp.Errors {
		return fmt.Errorf("bulk delete failed: %v", resp.Items)
	}

	return nil
}

// GetById 根据 ID 获取文档
func (r *BaseRepository[T]) GetById(ctx context.Context, id string) (*T, error) {
	// 执行获取请求
	resp, err := r.client.Get(r.indexName, id).Do(ctx)
	if err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	// 解析文档
	return r.unmarshalEntity(resp.Source_, id)
}

// unmarshalEntity 反序列化实体，先移除ID字段避免类型不匹配
func (r *BaseRepository[T]) unmarshalEntity(source []byte, id string) (*T, error) {
	// 先将JSON解析为map，移除id字段后再反序列化到实体
	var docMap map[string]interface{}
	if err := json.Unmarshal(source, &docMap); err != nil {
		return nil, err
	}

	// 移除id字段，避免JSON反序列化时类型不匹配
	delete(docMap, "id")
	delete(docMap, "ID")

	// 将处理后的map转换回JSON
	cleanedJSON, err := json.Marshal(docMap)
	if err != nil {
		return nil, err
	}

	// 反序列化到实体
	var entity T
	if err := json.Unmarshal(cleanedJSON, &entity); err != nil {
		return nil, err
	}

	// 设置文档 ID
	SetDocumentID(&entity, id)
	return &entity, nil
}

// GetOne 获取单个文档（根据条件）
func (r *BaseRepository[T]) GetOne(ctx context.Context, wrapper QueryWrapper[T]) (*T, error) {
	// 构建搜索请求
	req := wrapper.BuildSearchRequest()

	// 设置只返回一个结果
	size := 1
	req.Size = &size

	// 执行搜索请求
	resp, err := r.client.Search().Index(r.indexName).Request(req).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 检查是否有结果
	if resp.Hits.Total.Value == 0 {
		return nil, nil
	}

	// 解析第一个结果
	hit := resp.Hits.Hits[0]
	id := ""
	if hit.Id_ != nil {
		id = *hit.Id_
	}

	entity, err := r.unmarshalEntity(hit.Source_, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// List 获取文档列表
func (r *BaseRepository[T]) List(ctx context.Context, wrapper QueryWrapper[T]) ([]*T, error) {
	// 构建搜索请求
	req := wrapper.BuildSearchRequest()

	// 执行搜索请求
	resp, err := r.client.Search().Index(r.indexName).Request(req).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	entities := make([]*T, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		id := ""
		if hit.Id_ != nil {
			id = *hit.Id_
		}
		entity, err := r.unmarshalEntity(hit.Source_, id)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// Page 分页查询
func (r *BaseRepository[T]) Page(ctx context.Context, wrapper QueryWrapper[T], page, size int) (*PageResult[T], error) {
	// 计算偏移量
	offset := (page - 1) * size

	// 构建搜索请求
	req := wrapper.Limit(size).Offset(offset).BuildSearchRequest()

	// 执行搜索请求
	resp, err := r.client.Search().Index(r.indexName).Request(req).Do(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	entities := make([]*T, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		id := ""
		if hit.Id_ != nil {
			id = *hit.Id_
		}
		entity, err := r.unmarshalEntity(hit.Source_, id)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	// 计算总页数
	total := resp.Hits.Total.Value
	pages := int((total + int64(size) - 1) / int64(size))

	return &PageResult[T]{
		Total:   total,
		Pages:   pages,
		Current: page,
		Size:    size,
		Records: entities,
		Index:   r.indexName,
	}, nil
}

// Count 统计文档数量
func (r *BaseRepository[T]) Count(ctx context.Context, wrapper QueryWrapper[T]) (int64, error) {
	// 构建查询
	query := wrapper.BuildQuery()

	// 执行计数请求
	resp, err := r.client.Count().Index(r.indexName).Query(&query).Do(ctx)
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

// Exists 检查文档是否存在
func (r *BaseRepository[T]) Exists(ctx context.Context, id string) (bool, error) {
	// 执行存在请求
	exists, err := r.client.Exists(r.indexName, id).Do(ctx)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetOneByQueryStruct 根据查询结构体获取单个文档
func (r *BaseRepository[T]) GetOneByQueryStruct(ctx context.Context, queryStruct interface{}) (*T, error) {
	// 从查询结构体创建查询包装器
	wrapper := FromQueryStruct[T](queryStruct)
	// 调用现有的GetOne方法
	return r.GetOne(ctx, wrapper)
}

// ListByQueryStruct 根据查询结构体获取文档列表
func (r *BaseRepository[T]) ListByQueryStruct(ctx context.Context, queryStruct interface{}) ([]*T, error) {
	// 从查询结构体创建查询包装器
	wrapper := FromQueryStruct[T](queryStruct)
	// 调用现有的List方法
	return r.List(ctx, wrapper)
}

// PageByQueryStruct 根据查询结构体进行分页查询
func (r *BaseRepository[T]) PageByQueryStruct(ctx context.Context, queryStruct interface{}, page, size int) (*PageResult[T], error) {
	// 从查询结构体创建查询包装器
	wrapper := FromQueryStruct[T](queryStruct)
	// 调用现有的Page方法
	return r.Page(ctx, wrapper, page, size)
}

// CountByQueryStruct 根据查询结构体统计文档数量
func (r *BaseRepository[T]) CountByQueryStruct(ctx context.Context, queryStruct interface{}) (int64, error) {
	// 从查询结构体创建查询包装器
	wrapper := FromQueryStruct[T](queryStruct)
	// 调用现有的Count方法
	return r.Count(ctx, wrapper)
}

// Search 原生搜索
func (r *BaseRepository[T]) Search(ctx context.Context, index string, req *search.Request) (*search.Response, error) {
	// 执行搜索请求
	return r.client.Search().Index(index).Request(req).Do(ctx)
}

// setTimestamps 设置创建时间和更新时间
func (r *BaseRepository[T]) setTimestamps(entity *T) {
	timestamp := time.Now().UTC()

	v := reflect.ValueOf(entity).Elem()

	// 设置创建时间
	if createdAt := v.FieldByName("CreatedAt"); createdAt.IsValid() && createdAt.CanSet() {
		if createdAt.Kind() == reflect.Struct && createdAt.Type().String() == "time.Time" {
			// 只有当创建时间为空时才设置
			if createdAt.MethodByName("IsZero").Call(nil)[0].Bool() {
				createdAt.Set(reflect.ValueOf(timestamp))
			}
		}
	}

	// 设置更新时间
	r.setUpdateTimestamp(entity)
}

// setUpdateTimestamp 设置更新时间
func (r *BaseRepository[T]) setUpdateTimestamp(entity *T) {
	timestamp := time.Now().UTC()

	v := reflect.ValueOf(entity).Elem()

	// 设置更新时间
	if updatedAt := v.FieldByName("UpdatedAt"); updatedAt.IsValid() && updatedAt.CanSet() {
		if updatedAt.Kind() == reflect.Struct && updatedAt.Type().String() == "time.Time" {
			updatedAt.Set(reflect.ValueOf(timestamp))
		}
	}
}
