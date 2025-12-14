package models

import (
    "errors"
)

type State string

const (
    StateAvailable State = "available"
    StateInUse     State = "in-use"
    StateInactive  State = "inactive"
)

func (s State) Valid() bool {
    switch s {
    case StateAvailable, StateInUse, StateInactive:
        return true
    default:
        return false
    }
}

type Device struct {
    ID        int64     `json:"id" gorm:"primaryKey;column:id"`
    Name      string    `json:"name" gorm:"column:name;index:idx_devices_brand"`
    Brand     string    `json:"brand" gorm:"column:brand;index:idx_devices_brand"`
    State     State     `json:"state" gorm:"column:state;index:idx_devices_state"`
    CreatedAt FormattedTime `json:"created_at" gorm:"column:created_at;type:text"`
}

var (
    ErrInvalidState        = errors.New("invalid state")
    ErrCannotUpdateCreated = errors.New("creation time cannot be updated")
    ErrCannotUpdateFields  = errors.New("name/brand cannot be updated while in use")
    ErrCannotDeleteInUse   = errors.New("in-use devices cannot be deleted")
)

func (d *Device) ValidateNew() error {
    if d.Name == "" || d.Brand == "" {
        return errors.New("name and brand are required")
    }
    if !d.State.Valid() {
        return ErrInvalidState
    }
    if d.CreatedAt.IsZero() {
        d.CreatedAt = NowFormattedTime()
    }
    return nil
}
