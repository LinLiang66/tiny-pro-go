package controller

import (
	"net/http"
	"strconv"
	"strings"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"
	"tiny-admin-api-serve/utils"

	"github.com/gin-gonic/gin"
)

type I18Controller struct {
	i18n impl.I18Impl
}

func NewI18Controller() *I18Controller {
	return &I18Controller{
		i18n: impl.I18,
	}
}

// CreateI18Dto 创建国际化条目
func (ic I18Controller) CreateI18Dto(c *gin.Context) {
	var createI18Dto dto.CreateI18Dto
	if err := c.ShouldBindJSON(&createI18Dto); err != nil {
		utils.Waring(c, err.Error())
		return
	}

	result, err := ic.i18n.Create(createI18Dto)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// GetFormat 获取格式化的国际化数据
func (ic I18Controller) GetFormat(c *gin.Context) {
	langHeader := c.GetHeader("x-lang")
	if langHeader == "" {
		utils.Waring(c, "missing x-lang header")
		return
	}

	result, err := ic.i18n.GetFormat(langHeader)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindAll 查询所有国际化条目（分页）
func (ic I18Controller) FindAll(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	allParam := c.Query("all")

	// 处理 all 参数
	allBool := true
	if allParam != "" {
		allValue, _ := strconv.Atoi(allParam)
		allBool = !(allValue != 0)
	}

	// 解析 lang 参数
	var langIds []int64
	langParam := c.Query("lang")
	if langParam != "" {
		langStrs := strings.Split(langParam, ",")
		for _, langStr := range langStrs {
			if langId, err := strconv.ParseInt(strings.TrimSpace(langStr), 10, 64); err == nil {
				langIds = append(langIds, langId)
			}
		}
	}

	key := c.Query("key")
	content := c.Query("content")

	result, err := ic.i18n.FindAll(page, limit, allBool, langIds, key, content)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// FindOne 根据ID获取单个国际化条目
func (ic I18Controller) FindOne(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Waring(c, "invalid id parameter")
		return
	}

	result, err := impl.I18.GetById(id)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// Update 根据ID更新国际化条目
func (ic I18Controller) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Waring(c, "invalid id parameter")
		return
	}

	var updateDto dto.CreateI18Dto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		utils.Waring(c, err.Error())
		return
	}

	result, err := ic.i18n.UpdateById(id, updateDto)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// Remove 根据ID删除国际化条目
func (ic I18Controller) Remove(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Waring(c, "invalid id parameter")
		return
	}

	result, err := ic.i18n.RemoveById(id)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// BatchRemove 批量删除国际化条目
func (ic I18Controller) BatchRemove(c *gin.Context) {
	var ids []int64
	if err := c.ShouldBindJSON(&ids); err != nil {
		utils.Waring(c, err.Error())
		return
	}

	result, err := impl.I18.BatchDelete(ids)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}
