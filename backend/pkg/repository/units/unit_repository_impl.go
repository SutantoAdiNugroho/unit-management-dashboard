package units

import (
	"fmt"
	"unit-management-be/pkg/model/domain"

	"gorm.io/gorm"
)

type UnitRepositoryImpl struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &UnitRepositoryImpl{db: db}
}

func (u *UnitRepositoryImpl) Create(unit domain.Units) error {
	err := u.db.Create(unit).Error
	if err != nil {
		fmt.Printf("failed to create new unit: %v", err)
		return err
	}

	return nil
}
