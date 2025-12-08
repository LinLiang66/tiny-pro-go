package dto

type Lang struct {
	ID    int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name  string `json:"name" gorm:"column:name"`
	I18ns []I18  `json:"i18ns,omitempty" gorm:"foreignKey:LangID"`
}

// TableName 指定表名
func (Lang) TableName() string {
	return "lang"
}

type I18 struct {
	ID      int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Key     string `json:"key" gorm:"column:key"`
	Content string `json:"content" gorm:"column:content;type:text"`
	LangID  int64  `json:"langId" gorm:"column:lang_id"`
	Lang    Lang   `json:"lang,omitempty" gorm:"foreignKey:LangID"`
}

// TableName 指定表名
func (I18) TableName() string {
	return "i18"
}

type CreateI18Dto struct {
	Lang    string `json:"lang" binding:"required"`
	Key     string `json:"key" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type I18Vo struct {
	ID      int64  `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
	Lang    Lang   `json:"lang"`
}
