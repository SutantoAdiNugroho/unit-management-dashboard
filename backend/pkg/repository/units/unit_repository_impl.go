package units

import (
	"fmt"
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/dto/response"
	"unit-management-be/pkg/utils"

	"gorm.io/gorm"
)

type UnitRepositoryImpl struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &UnitRepositoryImpl{db: db}
}

func (u *UnitRepositoryImpl) Create(unit domain.Units) (domain.Units, error) {
	err := u.db.Create(&unit).Error
	if err != nil {
		fmt.Printf("failed to create new unit: %v", err)
		return unit, err
	}

	return unit, nil
}

func (u *UnitRepositoryImpl) GetByID(id string) (domain.Units, error) {
	response := domain.Units{}
	if err := u.db.Table("units").Where("id = ? AND deleted_at IS NULL", id).First(&response).Error; err != nil {
		fmt.Printf("failed to get unit by id: %v", err)
		return response, err
	}
	return response, nil
}

func (u *UnitRepositoryImpl) Delete(unit domain.Units) error {
	if err := u.db.Delete(&unit).Error; err != nil {
		fmt.Printf("failed to delete unit by id: %v", err)
		return err
	}

	return nil
}

func (u *UnitRepositoryImpl) FindAll(status, unitType, name string, page, size int) ([]response.UnitDetailResponse, int64, error) {
	units := make([]response.UnitDetailResponse, 0)

	selectStatement := "units.id AS ID, units.name AS Name, units.type AS Type, units.status AS Status"
	baseQuery := u.db.Table("units").Select(selectStatement).Where("units.deleted_at IS NULL")

	if !utils.IsEmptyString(status) {
		baseQuery = baseQuery.Where("units.status = ?", status)
	}

	if !utils.IsEmptyString(unitType) {
		baseQuery = baseQuery.Where("units.type = ?", unitType)
	}

	if !utils.IsEmptyString(name) {
		baseQuery = baseQuery.Where("LOWER(units.name) LIKE LOWER(?)", "%"+name+"%")
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		fmt.Printf("failed to count units: %v", err)
		return units, total, err
	}

	offset := (page - 1) * size
	paginateQuery := baseQuery.Limit(size).Offset(offset).Order("units.name ASC")
	if err := paginateQuery.Scan(&units).Error; err != nil {
		fmt.Printf("failed to scan units: %v", err)
		return units, total, err
	}

	return units, total, nil
}

func (u *UnitRepositoryImpl) Update(unit domain.Units) error {
	if err := u.db.Save(unit).Error; err != nil {
		fmt.Printf("failed to save unit: %v", err)
		return err
	}

	return nil
}
