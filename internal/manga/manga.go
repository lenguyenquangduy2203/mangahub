package manga

import (
	"context"
	"mangahub/data/models"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/database"
	gormHelper "mangahub/pkg/utils/gorm"
)

type Repository struct {
	DB database.DBConnector
}

func NewRepository(db database.DBConnector) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetMangaById(ctx context.Context, mangaID string) (*models.Manga, error) {
	manga, err := gormHelper.FindOne[models.Manga](r.DB.DB(), ctx, "id = ?", mangaID)

	if err != nil {
		return nil, err
	}

	return &manga, err
}

func (r *Repository) FindMangaByQuery(ctx context.Context, query dtos.MangaSearchQuery) ([]models.Manga, int64, error) {
	var mangas []models.Manga
	var total int64

	tx := r.DB.DB().WithContext(ctx).Model(&models.Manga{})

	if query.Title != "" {
		tx = tx.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if query.Author != "" {
		tx = tx.Where("author LIKE ?", "%"+query.Author+"%")
	}
	if query.Genre != "" {
		tx = tx.Where("genres LIKE ?", "%"+query.Genre+"%")
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}

	// Count total before pagination
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if query.Limit > 0 {
		tx = tx.Limit(query.Limit)
	}
	if query.Offset > 0 {
		tx = tx.Offset(query.Offset)
	}

	// Fetch page results
	if err := tx.Order("title ASC").Find(&mangas).Error; err != nil {
		return nil, 0, err
	}

	return mangas, total, nil
}
