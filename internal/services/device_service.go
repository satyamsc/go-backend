package services

import (
	"context"
	"errors"
	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type DeviceService struct {
	repo *repositories.DeviceRepository
}

func NewDeviceService(r *repositories.DeviceRepository) *DeviceService {
	return &DeviceService{repo: r}
}

func (s *DeviceService) Create(ctx context.Context, d *models.Device) (int64, error) {
	return s.repo.Create(ctx, d)
}
func (s *DeviceService) Get(ctx context.Context, id int64) (*models.Device, error) {
	return s.repo.Get(ctx, id)
}
func (s *DeviceService) List(ctx context.Context, brand, state string) ([]models.Device, error) {
	return s.repo.List(ctx, brand, state)
}

func (s *DeviceService) Update(ctx context.Context, id int64, incoming *models.Device) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !incoming.State.Valid() {
		return models.ErrInvalidState
	}
	if !incoming.CreatedAt.Equal(existing.CreatedAt) {
		return models.ErrCannotUpdateCreated
	}
	if existing.State == models.StateInUse && (incoming.Name != existing.Name || incoming.Brand != existing.Brand) {
		return models.ErrCannotUpdateFields
	}
	return s.repo.Update(ctx, id, incoming)
}

func (s *DeviceService) Patch(ctx context.Context, id int64, fields map[string]any) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if _, ok := fields["created_at"]; ok {
		return models.ErrCannotUpdateCreated
	}
	if existing.State == models.StateInUse {
		if _, ok := fields["name"]; ok {
			return models.ErrCannotUpdateFields
		}
		if _, ok := fields["brand"]; ok {
			return models.ErrCannotUpdateFields
		}
	}
	if v, ok := fields["state"]; ok {
		if str, ok2 := v.(string); ok2 {
			if !models.State(str).Valid() {
				return models.ErrInvalidState
			}
		} else {
			return errors.New("invalid state type")
		}
	}
	return s.repo.Patch(ctx, id, fields)
}

func (s *DeviceService) Delete(ctx context.Context, id int64) error {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if existing.State == models.StateInUse {
		return models.ErrCannotDeleteInUse
	}
	return s.repo.Delete(ctx, id)
}
