package user

import (
	"context"
	"errors"
	"mangahub/data/models"
	"mangahub/pkg/utils/database"
	gormHelper "mangahub/pkg/utils/gorm"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	DB database.DBConnector
}

func NewRepository(db database.DBConnector) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, hashedPassword string) (*models.User, error) {
	id := uuid.New().String()

	user := models.User{
		ID:           id,
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err := gormHelper.Create[models.User](r.DB.DB(), ctx, &user)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("duplicate uuid")
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := gormHelper.FindOne[models.User](r.DB.DB(), ctx, "username = ?", username)

	if err != nil {
		return nil, err
	}

	return &user, err
}
