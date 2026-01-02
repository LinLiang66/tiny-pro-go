package elastic

// Api 定义支持的API操作类型
type Api string

const (
	// ApiCreate 创建操作
	ApiCreate Api = "CREATE"
	// ApiGet 根据ID获取操作
	ApiGet Api = "GET"
	// ApiUpdate 更新操作
	ApiUpdate Api = "UPDATE"
	// ApiUpdateById 根据ID更新操作
	ApiUpdateById Api = "UPDATE_BY_ID"
	// ApiDelete 删除操作
	ApiDelete Api = "DELETE"
	// ApiBatchDelete 批量删除操作
	ApiBatchDelete Api = "BATCH_DELETE"
	// ApiList 获取列表操作
	ApiList Api = "LIST"
	// ApiPage 分页获取操作
	ApiPage Api = "PAGE"
	// ApiCount 统计操作
	ApiCount Api = "COUNT"
	// ApiExport 导出操作
	ApiExport Api = "EXPORT"
)

// AllApis 所有支持的API操作
var AllApis = []Api{
	ApiCreate,
	ApiGet,
	ApiUpdate,
	ApiUpdateById,
	ApiDelete,
	ApiBatchDelete,
	ApiList,
	ApiPage,
	ApiCount,
	ApiExport,
}

// Contains 检查是否包含指定的API操作
func (a Api) Contains(apis []Api) bool {
	for _, api := range apis {
		if a == api {
			return true
		}
	}
	return false
}
