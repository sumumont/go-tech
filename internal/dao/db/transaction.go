package db

import (
	"context"
	"github.com/rs/xid"
)
import "gorm.io/gorm"

type ctxTransactionKey struct{}

func Transaction(c context.Context, fn func(runCtx context.Context) error) error {
	sessionCtx := context.WithValue(c, "dbSessionId", xid.New().String())
	return GetDBWith(sessionCtx).Transaction(func(tx *gorm.DB) error {
		runCtx := context.WithValue(sessionCtx, ctxTransactionKey{}, tx)
		return fn(runCtx)
	})
}
func GetDBWith(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return db.WithContext(ctx)
	}

	tx, ok := ctx.Value(ctxTransactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return db.WithContext(ctx)
}
