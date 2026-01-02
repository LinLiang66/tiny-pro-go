package elastic

// PageRequest 分页请求结构体
type PageRequest struct {
	Page int `form:"page" json:"page" binding:"min=1"`         // 当前页码，从1开始
	Size int `form:"size" json:"size" binding:"min=1,max=100"` // 每页数量，默认10，最大100
}

// GetOffset 获取偏移量
func (pr *PageRequest) GetOffset() int {
	if pr.Page <= 0 {
		pr.Page = 1
	}
	if pr.Size <= 0 {
		pr.Size = 10
	}
	return (pr.Page - 1) * pr.Size
}

// GetSize 获取每页数量
func (pr *PageRequest) GetSize() int {
	if pr.Size <= 0 {
		return 10
	}
	if pr.Size > 100 {
		return 100
	}
	return pr.Size
}

// NewPageResult 创建分页响应
func NewPageResult[T any](records []*T, total int64, current, size int, index string) *PageResult[T] {
	// 计算总页数
	pages := int(total) / size
	if int(total)%size > 0 {
		pages++
	}

	return &PageResult[T]{
		Total:   total,
		Pages:   pages,
		Current: current,
		Size:    size,
		Records: records,
		Index:   index,
	}
}
