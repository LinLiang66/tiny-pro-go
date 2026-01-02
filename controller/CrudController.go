package controller

import (
	"context"
	"net/http"
	"strconv"
	"tiny-admin-api-serve/utils/elastic"

	"github.com/gin-gonic/gin"
)

// CrudController 通用CRUD控制器基础类
type CrudController[T any] struct {
	repository *elastic.BaseRepository[T]
}

// NewCrudController 创建新的CRUD控制器
func NewCrudController[T any]() (*CrudController[T], error) {
	// 初始化ES客户端
	if err := elastic.InitClient(); err != nil {
		return nil, err
	}

	// 创建仓库实例
	repo, err := elastic.NewBaseRepository[T]()
	if err != nil {
		return nil, err
	}

	return &CrudController[T]{
		repository: repo,
	}, nil
}

// Create 创建资源
func (c *CrudController[T]) Create(ctx *gin.Context) {
	ctxResponse := context.Background()

	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	id, err := c.repository.Insert(ctxResponse, &entity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id, "data": entity})
}

// GetById 根据ID获取资源
func (c *CrudController[T]) GetById(ctx *gin.Context) {
	ctxResponse := context.Background()

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	entity, err := c.repository.GetById(ctxResponse, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resource: " + err.Error()})
		return
	}

	if entity == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// Update 更新资源
func (c *CrudController[T]) Update(ctx *gin.Context) {
	ctxResponse := context.Background()

	var entity T
	if err := ctx.ShouldBindJSON(&entity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := c.repository.Update(ctxResponse, &entity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

// UpdateById 根据ID更新资源
func (c *CrudController[T]) UpdateById(ctx *gin.Context) {
	ctxResponse := context.Background()

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var updateData map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := c.repository.UpdateById(ctxResponse, id, updateData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Resource updated successfully"})
}

// Delete 根据ID删除资源
func (c *CrudController[T]) Delete(ctx *gin.Context) {
	ctxResponse := context.Background()

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := c.repository.DeleteById(ctxResponse, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// BatchDelete 批量删除资源
func (c *CrudController[T]) BatchDelete(ctx *gin.Context) {
	ctxResponse := context.Background()

	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := c.repository.DeleteBatch(ctxResponse, ids); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to batch delete resources: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Resources deleted successfully"})
}

// List 获取资源列表
func (c *CrudController[T]) List(ctx *gin.Context) {
	ctxResponse := context.Background()

	// 绑定查询条件
	var queryStruct interface{}
	if err := ctx.ShouldBindQuery(&queryStruct); err != nil {
		// 如果没有查询条件，创建一个空的查询结构体
		queryStruct = struct{}{}
	}

	entities, err := c.repository.ListByQueryStruct(ctxResponse, queryStruct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resources: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entities)
}

// Page 分页获取资源
func (c *CrudController[T]) Page(ctx *gin.Context) {
	ctxResponse := context.Background()

	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 绑定查询条件
	var queryStruct interface{}
	if err := ctx.ShouldBindQuery(&queryStruct); err != nil {
		// 如果没有查询条件，创建一个空的查询结构体
		queryStruct = struct{}{}
	}

	result, err := c.repository.PageByQueryStruct(ctxResponse, queryStruct, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resources: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Count 统计资源数量
func (c *CrudController[T]) Count(ctx *gin.Context) {
	ctxResponse := context.Background()

	// 绑定查询条件
	var queryStruct interface{}
	if err := ctx.ShouldBindQuery(&queryStruct); err != nil {
		// 如果没有查询条件，创建一个空的查询结构体
		queryStruct = struct{}{}
	}

	count, err := c.repository.CountByQueryStruct(ctxResponse, queryStruct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count resources: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}
