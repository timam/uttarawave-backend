package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type BuildingRepository interface {
	CreateBuilding(ctx context.Context, building *models.Building) error
	DeleteBuilding(id string) error
	UpdateBuilding(ctx context.Context, id string, updates map[string]interface{}) error
	GetBuildingByID(ctx context.Context, id string) (*models.Building, error)
	GetAllBuildings(ctx context.Context) ([]models.Building, error)
}

type GormBuildingRepository struct{}

func NewGormBuildingRepository() *GormBuildingRepository {
	return &GormBuildingRepository{}
}

func (r *GormBuildingRepository) CreateBuilding(ctx context.Context, building *models.Building) error {
	return db.DB.WithContext(ctx).Create(building).Error
}

func (r *GormBuildingRepository) DeleteBuilding(id string) error {
	return db.DB.Delete(&models.Building{}, "id = ?", id).Error
}

func (r *GormBuildingRepository) UpdateBuilding(ctx context.Context, id string, updates map[string]interface{}) error {
	return db.DB.WithContext(ctx).Model(&models.Building{}).Where("id = ?", id).Updates(updates).Error
}

func (r *GormBuildingRepository) GetBuildingByID(ctx context.Context, id string) (*models.Building, error) {
	var building models.Building
	if err := db.DB.WithContext(ctx).First(&building, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &building, nil
}

func (r *GormBuildingRepository) GetAllBuildings(ctx context.Context) ([]models.Building, error) {
	var buildings []models.Building
	if err := db.DB.WithContext(ctx).Find(&buildings).Error; err != nil {
		return nil, err
	}
	return buildings, nil
}
