package lib

import (
	"fmt"

	"gorm.io/gorm"
)

type QueryFiltersEqualToString struct {
	Value string `json:"value" binding:"required"`
}

func (f *QueryFiltersEqualToString) Scope(fieldName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			fmt.Sprintf("%s = ?", fieldName),
			f.Value,
		)
		return db
	}
}

type QueryFiltersEqualToUint struct {
	Value uint `json:"value" binding:"required"`
}

func (f *QueryFiltersEqualToUint) Scope(fieldName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			fmt.Sprintf("%s = ?", fieldName),
			f.Value,
		)
		return db
	}
}

type QueryFiltersInStringArray struct {
	Values []string `json:"values" binding:"required"`
}

func (f *QueryFiltersInStringArray) Scope(fieldName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			fmt.Sprintf("%s IN ?", fieldName),
			f.Values,
		)
		return db
	}
}

type QueryFiltersInUintArray struct {
	Values []uint `json:"values" binding:"required"`
}

func (f *QueryFiltersInUintArray) Scope(fieldName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			fmt.Sprintf("%s IN ?", fieldName),
			f.Values,
		)
		return db
	}
}
