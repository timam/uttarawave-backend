package repositories

import (
	"context"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/db"
)

type PackageRepository interface {
	CreateInternetPackage(ctx context.Context, pkg *models.InternetPackage) error
	UpdateInternetPackage(ctx context.Context, pkg *models.InternetPackage) error
	DeleteInternetPackage(ctx context.Context, id string) error
	GetInternetPackageByID(ctx context.Context, id string) (*models.InternetPackage, error)
	GetAllInternetPackages(ctx context.Context) ([]models.InternetPackage, error)

	CreateCableTVPackage(ctx context.Context, pkg *models.CableTVPackage) error
	UpdateCableTVPackage(ctx context.Context, pkg *models.CableTVPackage) error
	DeleteCableTVPackage(ctx context.Context, id string) error
	GetCableTVPackageByID(ctx context.Context, id string) (*models.CableTVPackage, error)
	GetAllCableTVPackages(ctx context.Context) ([]models.CableTVPackage, error)
}

type GormPackageRepository struct{}

func NewGormPackageRepository() *GormPackageRepository {
	return &GormPackageRepository{}
}

func (r *GormPackageRepository) CreateInternetPackage(ctx context.Context, pkg *models.InternetPackage) error {
	return db.DB.WithContext(ctx).Create(pkg).Error
}

func (r *GormPackageRepository) UpdateInternetPackage(ctx context.Context, pkg *models.InternetPackage) error {
	return db.DB.WithContext(ctx).Save(pkg).Error
}

func (r *GormPackageRepository) DeleteInternetPackage(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.InternetPackage{}, "id = ?", id).Error
}

func (r *GormPackageRepository) GetInternetPackageByID(ctx context.Context, id string) (*models.InternetPackage, error) {
	var pkg models.InternetPackage
	err := db.DB.WithContext(ctx).First(&pkg, "id = ?", id).Error
	return &pkg, err
}

func (r *GormPackageRepository) GetAllInternetPackages(ctx context.Context) ([]models.InternetPackage, error) {
	var packages []models.InternetPackage
	err := db.DB.WithContext(ctx).Find(&packages).Error
	return packages, err
}

func (r *GormPackageRepository) CreateCableTVPackage(ctx context.Context, pkg *models.CableTVPackage) error {
	return db.DB.WithContext(ctx).Create(pkg).Error
}

func (r *GormPackageRepository) UpdateCableTVPackage(ctx context.Context, pkg *models.CableTVPackage) error {
	return db.DB.WithContext(ctx).Save(pkg).Error
}

func (r *GormPackageRepository) DeleteCableTVPackage(ctx context.Context, id string) error {
	return db.DB.WithContext(ctx).Delete(&models.CableTVPackage{}, "id = ?", id).Error
}

func (r *GormPackageRepository) GetCableTVPackageByID(ctx context.Context, id string) (*models.CableTVPackage, error) {
	var pkg models.CableTVPackage
	err := db.DB.WithContext(ctx).First(&pkg, "id = ?", id).Error
	return &pkg, err
}

func (r *GormPackageRepository) GetAllCableTVPackages(ctx context.Context) ([]models.CableTVPackage, error) {
	var packages []models.CableTVPackage
	err := db.DB.WithContext(ctx).Find(&packages).Error
	return packages, err
}
