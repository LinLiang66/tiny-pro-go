package controller

import (
	"strconv"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"
	"tiny-admin-api-serve/utils"

	"github.com/gin-gonic/gin"
)

type LangController struct {
	langImpl impl.LangImpl
}

func NewLangController() *LangController {
	return &LangController{
		langImpl: impl.Lang,
	}
}

// CreateLang 创建语言
func (lc *LangController) CreateLang(c *gin.Context) {
	var createLangDto dto.Lang
	if err := c.ShouldBindJSON(&createLangDto); err != nil {
		utils.Waring(c, err.Error())
		return
	}

	result, err := lc.langImpl.Create(createLangDto)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// FindAllLang 获取所有语言
func (lc *LangController) FindAllLang(c *gin.Context) {
	result, err := lc.langImpl.FindAll()
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// UpdateLang 更新语言
func (lc *LangController) UpdateLang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Waring(c, "invalid id parameter")
		return
	}

	var createLangDto dto.Lang
	if err := c.ShouldBindJSON(&createLangDto); err != nil {
		utils.Waring(c, err.Error())
		return
	}

	result, err := lc.langImpl.Update(id, createLangDto)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}

// RemoveLang 删除语言
func (lc *LangController) RemoveLang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Waring(c, "invalid id parameter")
		return
	}

	result, err := lc.langImpl.Remove(id)
	if err != nil {
		utils.Waring(c, err.Error())
		return
	}

	utils.SuccessData(c, result)
}
