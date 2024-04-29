package funcs

import (
	"fmt"
	"log/slog"
	"rest/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDataWithLogWord[T any](c *gin.Context, keyWord string, model *T, scopes ...db.ScopeFunc) error {
	err := GetData(c, model, scopes...)
	if err != nil {
		slog.Warn(fmt.Sprintf("get %s info failed", keyWord),
			slog.String("current_user", c.GetString("current_user")),
			slog.String("err", err.Error()))
	}
	return nil
}

func GetData[T any](c *gin.Context, model *T, scopes ...db.ScopeFunc) error {
	var err error
	if err = c.BindUri(&model); err != nil {
		return err
	}
	tx := db.DB.Unscoped()
	return GetDataWithTransaction(c, tx, model, scopes...)
}

func GetBatchData[T any](c *gin.Context, models *[]T, scopes ...db.ScopeFunc) error {
	var err error
	if err = c.BindUri(&models); err != nil {
		return err
	}
	tx := db.DB.Unscoped()
	return GetBatchDataWithTransaction(c, tx, models, scopes...)
}

func GetDataWithTransaction[T any](c *gin.Context, tx *gorm.DB, model *T, scopes ...db.ScopeFunc) error {
	var err error
	c.BindUri(&model)
	tx = tx.Unscoped()
	for _, f := range scopes {
		tx = f(tx)
	}
	tx = tx.First(&model)
	err = tx.Error
	return err
}

func GetBatchDataWithTransaction[T any](c *gin.Context, tx *gorm.DB, models *[]T, scopes ...db.ScopeFunc) error {
	var err error
	c.BindUri(&models)
	tx = tx.Unscoped()
	for _, f := range scopes {
		tx = f(tx)
	}
	tx = tx.Find(&models)
	err = tx.Error
	return err
}
