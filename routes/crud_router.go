package routers

import (
	"tiny-admin-api-serve/controller"
	"tiny-admin-api-serve/utils/elastic"

	"github.com/gin-gonic/gin"
)

// CrudRouterConfig CRUD路由配置
type CrudRouterConfig[T any] struct {
	Path        string                        // 路由路径前缀
	Apis        []elastic.Api                 // 支持的API操作
	Middlewares []gin.HandlerFunc             // 中间件列表
	Controller  *controller.CrudController[T] // 控制器实例
}

// RegisterCrudRoutes 注册CRUD路由
func RegisterCrudRoutes[T any](engine *gin.Engine, config CrudRouterConfig[T]) error {
	// 如果没有指定控制器，创建一个新的
	if config.Controller == nil {
		ctrl, err := controller.NewCrudController[T]()
		if err != nil {
			return err
		}
		config.Controller = ctrl
	}

	// 创建路由组
	group := engine.Group(config.Path)

	// 添加中间件
	for _, middleware := range config.Middlewares {
		group.Use(middleware)
	}

	// 注册路由
	if elastic.ApiCreate.Contains(config.Apis) {
		group.POST("", config.Controller.Create)
	}

	if elastic.ApiUpdate.Contains(config.Apis) {
		group.PUT("", config.Controller.Update)
	}

	if elastic.ApiBatchDelete.Contains(config.Apis) {
		group.DELETE("/batch", config.Controller.BatchDelete)
	}

	// 将带具体路径的GET请求放在前面注册
	if elastic.ApiList.Contains(config.Apis) {
		group.GET("", config.Controller.List)
	}

	if elastic.ApiPage.Contains(config.Apis) {
		group.GET("/page", config.Controller.Page)
	}

	if elastic.ApiCount.Contains(config.Apis) {
		group.GET("/count", config.Controller.Count)
	}
	if elastic.ApiExport.Contains(config.Apis) {
		// 导出功能需要额外实现
		group.GET("/export", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Export not implemented yet"})
		})
	}
	// 将带参数的路由放在最后注册，避免冲突
	if elastic.ApiGet.Contains(config.Apis) {
		group.GET("/:id", config.Controller.GetById)
	}

	if elastic.ApiUpdateById.Contains(config.Apis) {
		group.PATCH("/:id", config.Controller.UpdateById)
	}

	if elastic.ApiDelete.Contains(config.Apis) {
		group.DELETE("/:id", config.Controller.Delete)
	}

	return nil
}

// RegisterDefaultCrudRoutes 注册默认的CRUD路由（包含所有API操作）
func RegisterDefaultCrudRoutes[T any](engine *gin.Engine, path string, middlewares ...gin.HandlerFunc) error {
	return RegisterCrudRoutes(engine, CrudRouterConfig[T]{
		Path:        path,
		Apis:        elastic.AllApis,
		Middlewares: middlewares,
	})
}
