package db

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(db *gorm.DB) error {
	if DB != nil {
		return errors.New("db already set")
	}
	DB = db
	return nil
}

type ScopeFunc = func(db *gorm.DB) *gorm.DB
type Condition struct {
	Condition string
	Value     any
}

func Joins(query string, args ...any) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Joins(query, args...)
		return db
	}
}

func Where(query any, args ...any) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(query, args...)
		return db
	}
}

func Or(query any, args ...any) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Or(query, args...)
		return db
	}
}

func Preloads(preloads ...string) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
		return db
	}
}

func Preload(query string, args ...any) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

func PreloadWithSearchCondition(query string, searchCondition *Condition) ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, searchCondition.Condition, searchCondition.Value)
	}
}

func PreloadWithSearchConditions(query string, searchConditions *[]*Condition) ScopeFunc {
	var (
		condition  string
		conditions []string
		values     []any
	)
	for _, searchCondition := range *searchConditions {
		conditions = append(conditions, searchCondition.Condition)
		values = append(values, searchCondition.Value)
	}
	condition = strings.Join(conditions, " AND ")
	params := []any{condition}
	params = append(params, values...)

	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, params...)
	}
}

func Unscoped() ScopeFunc {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Unscoped()
		return db
	}
}
