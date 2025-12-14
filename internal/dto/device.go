package dto

import "time"

type CreateDeviceRequest struct {
	Name  string `json:"name" binding:"required"`
	Brand string `json:"brand" binding:"required"`
	State string `json:"state" binding:"required,oneof=available in-use inactive"`
}

type UpdateDeviceRequest struct {
	Name      string     `json:"name" binding:"required"`
	Brand     string     `json:"brand" binding:"required"`
	State     string     `json:"state" binding:"required,oneof=available in-use inactive"`
	CreatedAt *time.Time `json:"created_at"`
}

type PatchDeviceRequest struct {
	Name  *string `json:"name" binding:"omitempty"`
	Brand *string `json:"brand" binding:"omitempty"`
	State *string `json:"state" binding:"omitempty,oneof=available in-use inactive"`
}
