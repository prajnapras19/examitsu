package lib

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
}
