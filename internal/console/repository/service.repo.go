package repository

import (
	"context"

	"github.com/vyolayer/vyolayer/internal/console/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	List(ctx context.Context) ([]*model.Service, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) List(ctx context.Context) ([]*model.Service, error) {

	var services []*model.Service

	err := r.db.WithContext(ctx).
		Find(&services).Error

	return services, err
}
