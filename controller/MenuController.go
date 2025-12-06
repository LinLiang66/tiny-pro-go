package controller

import (
	"net/http"
	"strconv"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService impl.MenuImpl
}

func NewMenuController() *MenuController {
	return &MenuController{
		menuService: impl.Menu,
	}
}

// GetMenus 根据邮箱获取菜单
func (mc *MenuController) GetMenus(c *gin.Context) {
	email := c.Param("email")

	menus, err := mc.menuService.GetMenubyEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// GetAll 获取所有菜单
func (mc *MenuController) GetAll(c *gin.Context) {
	menus, err := mc.menuService.FindAllMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

// Create  创建菜单
func (mc *MenuController) Create(c *gin.Context) {
	var createMenuDto dto.Menu
	if err := c.ShouldBindJSON(&createMenuDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu, err := mc.menuService.CreateMenu(createMenuDto, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// Update 更新菜单
func (mc *MenuController) Update(c *gin.Context) {
	var updateMenuDto dto.Menu
	if err := c.ShouldBindJSON(&updateMenuDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := mc.menuService.UpdateMenu(updateMenuDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, success)
}

// Delete 删除菜单
func (mc *MenuController) Delete(c *gin.Context) {
	idStr := c.Query("id")
	parentIdStr := c.Query("parentId")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id parameter"})
		return
	}

	parentId, err := strconv.ParseInt(parentIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parentId parameter"})
		return
	}

	menu, err := mc.menuService.DeleteMenu(id, parentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}
