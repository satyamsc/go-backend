package repositories

import (
    "context"
    "go-backend/internal/models"
    "gorm.io/gorm"
)

type DeviceRepository struct{ db *gorm.DB }

func NewDeviceRepository(db *gorm.DB) *DeviceRepository { return &DeviceRepository{db: db} }

func (r *DeviceRepository) Create(ctx context.Context, d *models.Device) (int64, error) {
    if err := d.ValidateNew(); err != nil { return 0, err }
    if err := r.db.WithContext(ctx).Create(d).Error; err != nil { return 0, err }
    return d.ID, nil
}

func (r *DeviceRepository) Get(ctx context.Context, id int64) (*models.Device, error) {
    var d models.Device
    if err := r.db.WithContext(ctx).First(&d, id).Error; err != nil { return nil, err }
    return &d, nil
}

func (r *DeviceRepository) List(ctx context.Context, brand, state string) ([]models.Device, error) {
    var list []models.Device
    q := r.db.WithContext(ctx).Model(&models.Device{})
    if brand != "" { q = q.Where("brand = ?", brand) }
    if state != "" { q = q.Where("state = ?", state) }
    if err := q.Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func (r *DeviceRepository) Update(ctx context.Context, id int64, d *models.Device) error {
    return r.db.WithContext(ctx).Model(&models.Device{}).Where("id = ?", id).Updates(map[string]any{"name": d.Name, "brand": d.Brand, "state": d.State}).Error
}

func (r *DeviceRepository) Patch(ctx context.Context, id int64, fields map[string]any) error {
    return r.db.WithContext(ctx).Model(&models.Device{}).Where("id = ?", id).Updates(fields).Error
}

func (r *DeviceRepository) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).Delete(&models.Device{}, id).Error
}
