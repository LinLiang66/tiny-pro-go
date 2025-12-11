package routers

import (
	"tiny-admin-api-serve/controller"
	"tiny-admin-api-serve/middleware"

	"github.com/gin-gonic/gin"
)

// RouterUser 用户路由
func RouterUser(engine *gin.Engine) {
	// 认证相关路由
	authController := controller.NewAuthController()
	authGroup := engine.Group("/auth")
	{
		authGroup.POST("/login", middleware.IsPublic(), authController.Login)
		authGroup.POST("/register", middleware.IsPublic(), authController.Register)
		authGroup.GET("/profile", authController.Profile)
		authGroup.POST("/logout", authController.Logout)
		authGroup.POST("/refresh", authController.Refresh)
	}
	userController := controller.NewUserController()
	// 用户相关路由
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/reg", userController.Register)
		userGroup.GET("/info/:email", userController.GetUserInfo)
		userGroup.GET("/info", userController.GetUserInfo)
		userGroup.GET("/info/", userController.GetUserInfo)
		userGroup.GET("/info/:email/", userController.GetUserInfo)
		userGroup.DELETE("/:email", userController.DelUser)
		userGroup.PATCH("/update", userController.UpdateUser)
		userGroup.GET("", userController.GetAllUser)
		userGroup.PATCH("/admin/updatePwd", userController.UpdatePwdAdmin)
		userGroup.PATCH("/updatePwd", userController.UpdatePwdUser)
		userGroup.POST("/batch", userController.BatchRemoveUser)
	}

	// 角色相关路由
	roleController := controller.NewRoleController()
	roleGroup := engine.Group("/role")
	{
		roleGroup.POST("", roleController.Create)
		roleGroup.GET("", roleController.GetAllRole)
		roleGroup.GET("/detail", roleController.GetAllRoleDetail)
		roleGroup.PATCH("", roleController.UpdateRole)
		roleGroup.DELETE("/:id", roleController.DeleteRole)
		roleGroup.GET("/info/:id", roleController.GetRoleInfo)
	}

	// 菜单相关路由
	menuController := controller.NewMenuController()
	menuGroup := engine.Group("/menu")
	{
		menuGroup.GET("/role/:email", menuController.GetMenus)
		menuGroup.POST("", menuController.Create)
		menuGroup.GET("", menuController.GetAll)
		menuGroup.PATCH("", menuController.Update)
		menuGroup.DELETE("/:id", menuController.Delete)
	}

	// 权限相关路由
	permissionController := controller.NewPermissionController()
	permissionGroup := engine.Group("/permission")
	{
		permissionGroup.POST("", permissionController.Create)
		permissionGroup.GET("", permissionController.GetAll)
		permissionGroup.PATCH("", permissionController.Update)
		permissionGroup.DELETE("/:id", permissionController.Delete)
	}

	// i18n相关路由
	i18Controller := controller.NewI18Controller()
	i18Group := engine.Group("/i18")
	{
		i18Group.POST("", i18Controller.CreateI18Dto)
		i18Group.GET("/format", i18Controller.GetFormat)
		i18Group.GET("", i18Controller.FindAll)
		i18Group.GET("/:id", i18Controller.FindOne)
		i18Group.PATCH("/:id", i18Controller.Update)
		i18Group.DELETE("/:id", i18Controller.Remove)
		i18Group.POST("/batch", i18Controller.BatchRemove)
	}

	// 语言相关路由
	langController := controller.NewLangController()
	langGroup := engine.Group("/lang")
	{
		langGroup.POST("", langController.CreateLang)
		langGroup.GET("", langController.FindAllLang)
		langGroup.PATCH("/:id", langController.UpdateLang)
		langGroup.DELETE("/:id", langController.RemoveLang)
	}

}
