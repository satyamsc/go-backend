package dto

import (
    "go-backend/internal/models"
)

type DeviceResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Brand     string `json:"brand"`
    State     string `json:"state"`
    CreatedAt string `json:"created_at"`
}

func FromModel(d *models.Device) DeviceResponse {
    return DeviceResponse{
        ID:        d.ID,
        Name:      d.Name,
        Brand:     d.Brand,
        State:     string(d.State),
        CreatedAt: d.CreatedAt.Time.UTC().Format(models.DbTimeLayout),
    }
}

func FromModels(list []models.Device) []DeviceResponse {
    out := make([]DeviceResponse, 0, len(list))
    for i := range list {
        out = append(out, FromModel(&list[i]))
    }
    return out
}
