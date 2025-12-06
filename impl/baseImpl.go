package impl

import (
	"errors"
	"tiny-admin-api-serve/utils"
)

type BaseImpl struct {
}

// Create 新增
func (b BaseImpl) Create(model interface{}) (err error) {
	return utils.Db.DB.Create(model).Error
}

// Get 根据Id查询详情
func (b BaseImpl) Get(model interface{}, conds ...interface{}) error {
	if len(conds) < 1 {
		return errors.New("请输入id")
	}
	return utils.Db.DB.First(&model, conds).Error
}

// Update 更新
func (b BaseImpl) Update(value interface{}) (err error) {
	return utils.Db.DB.Save(value).Error
}

// Delete 删除
func (b BaseImpl) Delete(value interface{}) (err error) {
	return utils.Db.DB.Delete(value).Error
}
