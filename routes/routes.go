package routers

import (
	"tiny-admin-api-serve/controller"
	"tiny-admin-api-serve/middleware"

	"github.com/gin-gonic/gin"
)

// RouterUser 用户路由
func RouterUser(engine *gin.Engine) {
	authController := controller.NewAuthController()
	// 公开路由 - 不需要鉴权
	openGroup := engine.Group("/auth")
	{
		openGroup.POST("/login", authController.Login)
		openGroup.POST("/register", authController.Register)
	}

	// 受保护路由 - 需要鉴权
	authGroup := engine.Group("/auth")
	authGroup.Use(middleware.Auth.AuthRequired())
	{
		authGroup.GET("/profile", authController.Profile)
		authGroup.POST("/logout", authController.Logout)
		authGroup.POST("/refresh", authController.Refresh)
	}
	userController := controller.NewUserController()
	userGroup := engine.Group("/user")
	userGroup.Use(middleware.Auth.AuthRequired())
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
	roleController := controller.NewRoleController()
	roleGroup := engine.Group("/role")
	roleGroup.Use(middleware.Auth.AuthRequired())
	{
		roleGroup.POST("", roleController.Create)
		roleGroup.GET("", roleController.GetAllRole)
		roleGroup.GET("/detail", roleController.GetAllRoleDetail)
		roleGroup.PATCH("", roleController.UpdateRole)
		roleGroup.DELETE("/:id", roleController.DeleteRole)
		roleGroup.GET("/info/:id", roleController.GetRoleInfo)
	}
	menuController := controller.NewMenuController()
	menuGroup := engine.Group("/menu")
	menuGroup.Use(middleware.Auth.AuthRequired())
	{
		menuGroup.GET("/role/:email", menuController.GetMenus)
		menuGroup.POST("", menuController.Create)
		menuGroup.GET("", menuController.GetAll)
		menuGroup.PATCH("", menuController.Update)
		menuGroup.DELETE("/:id", menuController.Delete)
	}
	permissionController := controller.NewPermissionController()
	permissionGroup := engine.Group("/permission")
	permissionGroup.Use(middleware.Auth.AuthRequired())
	{
		permissionGroup.POST("", permissionController.Create)
		permissionGroup.GET("", permissionController.GetAll)
		permissionGroup.PATCH("", permissionController.Update)
		permissionGroup.DELETE("/:id", permissionController.Delete)
	}
	i18Controller := controller.NewI18Controller()
	i18Group := engine.Group("/i18")
	i18Group.Use(middleware.Auth.AuthRequired())
	{
		i18Group.POST("", i18Controller.CreateI18Dto)
		i18Group.GET("/format", i18Controller.GetFormat)
		i18Group.GET("", i18Controller.FindAll)
		i18Group.GET("/:id", i18Controller.FindOne)
		i18Group.PATCH("/:id", i18Controller.Update)
		i18Group.DELETE("/:id", i18Controller.Remove)
		i18Group.POST("/batch", i18Controller.BatchRemove)
	}
	langController := controller.NewLangController()
	langGroup := engine.Group("/lang")
	langGroup.Use(middleware.Auth.AuthRequired())
	{
		langGroup.POST("", langController.CreateLang)
		langGroup.GET("", langController.FindAllLang)
		langGroup.PATCH("/:id", langController.UpdateLang)
		langGroup.DELETE("/:id", langController.RemoveLang)
	}

}
