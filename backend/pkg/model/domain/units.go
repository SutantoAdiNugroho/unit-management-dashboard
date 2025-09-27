package domain

import (
	"time"
	"unit-management-be/pkg/model/domain/enum"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Units struct {
	ID          uuid.UUID       `gorm:"type:varchar(36);primary_key" json:"id"`
	Name        string          `gorm:"type:varchar(255)" json:"name"`
	Type        enum.UnitType   `gorm:"type:enum('capsule', 'cabin')" json:"type"`
	Status      enum.UnitStatus `gorm:"type:enum('Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed')" json:"status"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
	LastUpdated time.Time       `gorm:"autoUpdateTime" json:"lastUpdated"`
}

func (u *Units) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

func (u *Units) TableName() string {
	return "units"
}
