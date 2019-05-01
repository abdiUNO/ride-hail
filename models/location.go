package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Location struct {
	ID        string     `gorm:"primary_key;type:varchar(255);"`
	Latitude  float64    `json:"latitude";gorm:"default:0"`
	Longitude float64    `json:"longitude";gorm:"default:0"`
	Address   string     `gorm:"type:varchar(255);"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func (location *Location) BeforeCreate(scope *gorm.Scope) error {
	u1 := uuid.Must(uuid.NewV4())
	scope.SetColumn("ID", u1.String())
	return nil
}
