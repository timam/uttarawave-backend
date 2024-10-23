package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type PackageRepository interface {
	CreatePackage(ctx context.Context, pkg *models.Package) error
	GetPackageByID(ctx context.Context, id string) (*models.Package, error)
	GetAllPackages(ctx context.Context, packageType string, page, pageSize int) ([]models.Package, int64, error)
	DeletePackage(ctx context.Context, id string) error
}

type GormPackageRepository struct{}

func NewGormPackageRepository() *GormPackageRepository {
	return &GormPackageRepository{}
}

func (r *GormPackageRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	return db.DB.WithContext(ctx).Create(pkg).Error
}

func (r *GormPackageRepository) GetPackageByID(ctx context.Context, id string) (*models.Package, error) {
	var pkg models.Package
	err := db.DB.WithContext(ctx).First(&pkg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (r *GormPackageRepository) GetAllPackages(ctx context.Context, packageType string, page, pageSize int) ([]models.Package, int64, error) {
	var packages []models.Package
	var total int64
	query := db.DB.WithContext(ctx)

	if packageType != "" {
		query = query.Where("type = ?", packageType)
	}

	err := query.Model(&models.Package{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&packages).Error
	if err != nil {
		return nil, 0, err
	}

	return packages, total, nil
}

func (r *GormPackageRepository) DeletePackage(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.Package{}, "id = ?", id).Error
}
