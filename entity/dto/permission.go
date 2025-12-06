package dto

type Permission struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`
	Desc string `json:"desc" gorm:"column:desc"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permission"
}
