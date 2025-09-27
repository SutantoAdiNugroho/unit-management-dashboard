package units

import (
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/dto/response"
)

type UnitRepository interface {
	Create(unit domain.Units) error
	GetByID(id string) (domain.Units, error)
	Delete(unit domain.Units) error
	FindAll(status, unitType, name string, page, size int) ([]response.UnitDetailResponse, int64, error)
	Update(unit domain.Units) error
}
