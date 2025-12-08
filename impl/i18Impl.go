package impl

import (
	"errors"
	"strconv"
	"strings"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"
)

type I18Impl struct {
	BaseImpl
}

var I18 = I18Impl{}

// Create 创建国际化条目
func (i I18Impl) Create(createI18Dto dto.CreateI18Dto) (*dto.I18, error) {
	// 查找语言
	var lang dto.Lang
	langId, _ := strconv.ParseInt(createI18Dto.Lang, 10, 64)
	err := utils.Db.DB.Where("id = ?", langId).First(&lang).Error
	if err != nil {
		return nil, errors.New("language not found")
	}

	// 校验 key + lang 是否已存在
	var existingI18 dto.I18
	err = utils.Db.DB.Where("key = ? AND lang_id = ?", createI18Dto.Key, langId).First(&existingI18).Error
	if err == nil {
		return nil, errors.New("i18n entry already exists")
	}

	// 创建新的国际化条目
	newI18 := dto.I18{
		Key:     createI18Dto.Key,
		Content: createI18Dto.Content,
		LangID:  langId,
	}

	result := utils.Db.DB.Create(&newI18)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newI18, nil
}

// GetFormat 获取格式化的国际化数据
func (i I18Impl) GetFormat(langHeader string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	// 查找语言
	var lang dto.Lang
	err := utils.Db.DB.Where("name = ?", langHeader).First(&lang).Error
	if err != nil {
		return nil, errors.New("language not found")
	}

	// 查询该语言的所有国际化条目
	var i18List []dto.I18
	err = utils.Db.DB.Where("lang_id = ?", lang.ID).Find(&i18List).Error
	if err != nil {
		return nil, err
	}

	i18map := make(map[string]string)
	for _, item := range i18List {
		i18map[item.Key] = item.Content
	}

	result[langHeader] = i18map
	return result, nil
}

// FindAll 查询所有国际化条目（分页）
func (i I18Impl) FindAll(page, limit int, allBool bool, langIds []int64, key, content string) (*dto.PageWrapper[dto.I18Vo], error) {
	var i18List []dto.I18
	var total int64

	// 构建查询
	query := utils.Db.DB.Model(&dto.I18{})

	// 按 lang 过滤
	if len(langIds) > 0 {
		query = query.Where("lang_id IN ?", langIds)
	}

	// 按 content 过滤
	if content != "" {
		if strings.Contains(content, "%") {
			query = query.Where("content LIKE ?", content)
		} else {
			query = query.Where("content = ?", content)
		}
	}

	// 按 key 过滤
	if key != "" {
		if strings.Contains(key, "%") {
			query = query.Where("key LIKE ?", key)
		} else {
			query = query.Where("key = ?", key)
		}
	}

	// 获取总数
	query.Count(&total)

	// 处理分页
	if allBool && page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	// 执行查询
	result := query.Find(&i18List)
	if result.Error != nil {
		return nil, result.Error
	}

	// 转换为 VO 对象
	var voList []dto.I18Vo
	for _, item := range i18List {
		// 获取语言信息
		var lang dto.Lang
		utils.Db.DB.Where("id = ?", item.LangID).First(&lang)

		vo := dto.I18Vo{
			ID:      item.ID,
			Key:     item.Key,
			Content: item.Content,
			Lang: dto.Lang{
				ID:   lang.ID,
				Name: lang.Name,
			},
		}
		voList = append(voList, vo)
	}

	// 计算分页信息
	totalPages := 1
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	pageWrapper := dto.NewPageWrapper[dto.I18Vo](
		voList,
		total,
		len(voList),
		limit,
		totalPages,
		page,
	)

	return pageWrapper, nil
}

// UpdateById 根据ID更新国际化条目
func (i I18Impl) UpdateById(id int64, updateDto dto.CreateI18Dto) (*dto.I18Vo, error) {
	var i18 dto.I18
	err := utils.Db.DB.Where("id = ?", id).First(&i18).Error
	if err != nil {
		return nil, errors.New("i18n entry not found")
	}

	// 更新字段
	if updateDto.Key != "" {
		i18.Key = updateDto.Key
	}
	if updateDto.Content != "" {
		i18.Content = updateDto.Content
	}

	if updateDto.Lang != "" {
		langId, err := strconv.ParseInt(updateDto.Lang, 10, 64)
		if err != nil {
			return nil, errors.New("invalid language id")
		}

		var lang dto.Lang
		err = utils.Db.DB.Where("id = ?", langId).First(&lang).Error
		if err != nil {
			return nil, errors.New("language not found")
		}

		i18.LangID = langId
	}

	// 保存更新
	result := utils.Db.DB.Save(&i18)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取语言信息
	var lang dto.Lang
	utils.Db.DB.Where("id = ?", i18.LangID).First(&lang)

	// 构造返回对象
	i18Vo := dto.I18Vo{
		ID:      i18.ID,
		Key:     i18.Key,
		Content: i18.Content,
		Lang: dto.Lang{
			ID:   lang.ID,
			Name: lang.Name,
		},
	}

	return &i18Vo, nil
}

// GetById 根据ID获取国际化条目
func (i I18Impl) GetById(id int64) (*dto.I18Vo, error) {
	var i18 dto.I18
	err := utils.Db.DB.Where("id = ?", id).First(&i18).Error
	if err != nil {
		return nil, errors.New("i18n entry not found")
	}

	// 获取语言信息
	var lang dto.Lang
	utils.Db.DB.Where("id = ?", i18.LangID).First(&lang)

	i18Vo := dto.I18Vo{
		ID:      i18.ID,
		Key:     i18.Key,
		Content: i18.Content,
		Lang: dto.Lang{
			ID:   lang.ID,
			Name: lang.Name,
		},
	}

	return &i18Vo, nil
}

// RemoveById 根据ID删除国际化条目
func (i I18Impl) RemoveById(id int64) (*dto.I18, error) {
	var i18 dto.I18
	err := utils.Db.DB.Where("id = ?", id).First(&i18).Error
	if err != nil {
		return nil, errors.New("i18n entry not found")
	}

	result := utils.Db.DB.Delete(&i18)
	if result.Error != nil {
		return nil, result.Error
	}

	return &i18, nil
}

// BatchDelete 批量删除国际化条目
func (i I18Impl) BatchDelete(ids []int64) ([]dto.I18, error) {
	var i18List []dto.I18
	err := utils.Db.DB.Where("id IN ?", ids).Find(&i18List).Error
	if err != nil {
		return nil, errors.New("failed to find i18n entries")
	}

	result := utils.Db.DB.Where("id IN ?", ids).Delete(&dto.I18{})
	if result.Error != nil {
		return nil, errors.New("failed to delete i18n entries")
	}

	return i18List, nil
}
