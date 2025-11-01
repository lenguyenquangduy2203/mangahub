package gorm

import (
	"context"

	"gorm.io/gorm"
)

func Create[T any](db *gorm.DB, ctx context.Context, dataPtr *T) error {
	return gorm.G[T](db).Create(ctx, dataPtr)
}

func FindOne[T any](db *gorm.DB, ctx context.Context, sqlTerm string, value any) (T, error) {
	return gorm.G[T](db).Where(sqlTerm, value).First(ctx)
}
