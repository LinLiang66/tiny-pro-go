package dto

type PaginationQueryDto struct {
	Page  int `json:"page" form:"page" binding:"min=1,max=500"`
	Limit int `json:"limit" form:"limit" binding:"min=1,max=500"`
}

// NewPaginationQueryDto 设置默认值
func NewPaginationQueryDto() PaginationQueryDto {
	return PaginationQueryDto{
		Page:  1,
		Limit: 10,
	}
}

// PageWrapper 分页包装器
type PageWrapper[T any] struct {
	Items []T  `json:"items"` // 当前页数据
	Meta  Meta `json:"meta"`
}

// Meta 分页元数据
type Meta struct {
	TotalElements    int64 `json:"totalElements"`
	NumberOfElements int   `json:"numberOfElements"`
	Size             int   `json:"size"`
	TotalPages       int   `json:"totalPages"`
	Number           int   `json:"number"` // 从 1 开始
}

// NewPageWrapper 创建新的分页包装器
func NewPageWrapper[T any](content []T, totalElements int64, numberOfElements int, size int, totalPages int, number int) *PageWrapper[T] {
	return &PageWrapper[T]{
		Items: content,
		Meta: Meta{
			TotalElements:    totalElements,
			NumberOfElements: numberOfElements,
			Size:             size,
			TotalPages:       totalPages,
			Number:           number,
		},
	}
}
