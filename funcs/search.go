package funcs

import (
	"github.com/moonlit0114/rest/db"
	"gorm.io/gorm"
)

type SearchParams = map[string]string

func GetSearchContent[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...db.ScopeFunc) (total int64, result *[]T, err error) {
	var (
		models = []T{}
	)
	tx = tx.Model(&models)
	// 处理搜索条件
	if searchPageParams.ParseCondition(); err != nil {
		return 0, nil, err
	}
	for _, condition := range searchPageParams.SearchConditions {
		tx = tx.Where(condition.Condition, condition.Value)
	}

	// 处理翻页参数
	if searchPageParams.OrderBy == "" {
		searchPageParams.OrderBy = "ID"
	}
	tx = tx.Count(&total).Order(searchPageParams.OrderBy)

	// 处理传入的scopes函数
	for _, f := range scopes {
		tx = f(tx)
	}

	// 获取查询结果及错误信息
	if err = tx.Find(&models).Error; err != nil {
		return 0, nil, err
	}
	return total, &models, err
}

func GetSearchData[T any](tx *gorm.DB,
	searchPageParams *SearchPageParams,
	scopes ...db.ScopeFunc) (*PageResult, error) {
	var (
		models *[]T
		total  int64
		err    error
	)
	if total, models, err = GetSearchContent[T](tx, searchPageParams, scopes...); err != nil {
		return nil, err
	}
	result := NewPageResult(searchPageParams.PageParams, len(*models), int(total), models)
	return result, err
}
