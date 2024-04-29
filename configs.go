package rest

import (
	"github.com/moonlit0114/rest/db"
	"gorm.io/gorm"
)

func InitDB(DB *gorm.DB) error {
	return db.InitDB(DB)
}
