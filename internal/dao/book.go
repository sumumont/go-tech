package dao

import (
	"context"
	"go-tech/internal/dao/db"
	"go-tech/internal/dao/model"
	"go-tech/internal/logging"
)

type BookDao[T model.Book] struct{}

func (b BookDao[T]) FindAll(ctx context.Context, param model.Book, page *db.BaseList) ([]T, int64, error) {
	var result []T
	var total int64
	var err error
	tx := db.GetDBWith(ctx).Model(T{})
	if len(param.Name) != 0 {
		tx = tx.Where("name ~ ?", param.Name)
	}
	if len(param.Author) != 0 {
		tx = tx.Where("author ~ ?", param.Author)
	}
	tx, err = db.PageList(tx, page)
	if err != nil {
		return result, 0, err
	}
	tx.Count(&total)
	err = tx.Find(&result).Error
	if err != nil {
		logging.ErrorStack(ctx, err).Send()
		return result, 0, err
	}
	return result, total, err
}
func (b BookDao[T]) Save(ctx context.Context, param *T) error {
	err := db.GetDBWith(ctx).Save(param).Error
	if err != nil {
		logging.ErrorStack(ctx, err).Send()
	}
	return err
}
