package lib

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedBy string
	UpdatedBy string
	DeletedBy string
}

func (m *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	return tx.Model(tx.Statement.Model).Update("deleted_by", m.DeletedBy).Error
}
