package impl

import (
	"errors"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"
)

type LangImpl struct {
	BaseImpl
}

var Lang = LangImpl{}

// Create 创建语言
func (l LangImpl) Create(createLangDto dto.Lang) (*dto.Lang, error) {
	// 检查语言是否已存在
	var existingLang dto.Lang
	err := utils.Db.DB.Where("name = ?", createLangDto.Name).First(&existingLang).Error
	if err == nil {
		return nil, errors.New("language already exists")
	}

	// 创建新语言
	newLang := dto.Lang{
		Name: createLangDto.Name,
	}

	result := utils.Db.DB.Create(&newLang)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newLang, nil
}

// FindAll 获取所有语言
func (l LangImpl) FindAll() ([]dto.Lang, error) {
	var langs []dto.Lang
	result := utils.Db.DB.Find(&langs)
	if result.Error != nil {
		return nil, result.Error
	}
	return langs, nil
}

// Update 更新语言
func (l LangImpl) Update(id int, createLangDto dto.Lang) (*dto.Lang, error) {
	var lang dto.Lang
	err := utils.Db.DB.Where("id = ?", id).First(&lang).Error
	if err != nil {
		return nil, errors.New("language not found")
	}

	lang.Name = createLangDto.Name
	result := utils.Db.DB.Save(&lang)
	if result.Error != nil {
		return nil, result.Error
	}

	return &lang, nil
}

// Remove 删除语言
func (l LangImpl) Remove(id int) (*dto.Lang, error) {
	var lang dto.Lang
	err := utils.Db.DB.Where("id = ?", id).First(&lang).Error
	if err != nil {
		return nil, errors.New("language not found")
	}

	// 删除关联的国际化条目
	utils.Db.DB.Where("lang_id = ?", lang.ID).Delete(&dto.I18{})

	// 删除语言
	result := utils.Db.DB.Delete(&lang)
	if result.Error != nil {
		return nil, result.Error
	}

	return &lang, nil
}
