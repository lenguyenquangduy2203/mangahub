package gorm

import (
	"context"

	"gorm.io/gorm"
)

func Create[T any](db *gorm.DB, ctx context.Context, dataPtr *T) error {
	return gorm.G[T](db).Create(ctx, dataPtr)
}

func FindOne[T any](db *gorm.DB, ctx context.Context, arg any, values ...any) (T, error) {
	if _, isString := arg.(string); isString {
		return gorm.G[T](db).Where(arg, values...).First(ctx)
	}
	return gorm.G[T](db).Where(arg).First(ctx)
}

func Find[T any](db *gorm.DB, ctx context.Context, arg any, values ...any) ([]T, error) {
	if _, isString := arg.(string); isString {
		return gorm.G[T](db).Where(arg, values...).Find(ctx)
	}
	return gorm.G[T](db).Where(arg).Find(ctx)
}
