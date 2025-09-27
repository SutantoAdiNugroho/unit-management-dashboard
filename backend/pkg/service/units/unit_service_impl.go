package units

import (
	"net/http"
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/domain/enum"
	"unit-management-be/pkg/model/dto/request"
	unitrepository "unit-management-be/pkg/repository/units"
)

type UnitServiceImpl struct {
	unitRepository unitrepository.UnitRepository
}

func NewUnitService(unitRepository unitrepository.UnitRepository) UnitService {
	return &UnitServiceImpl{unitRepository: unitRepository}
}
func (u *UnitServiceImpl) CreateUnit(request request.CreateUnitDto) *handler.CustomError {
	status, isValidStatus := enum.ParseUnitStatus(request.Status)
	if !isValidStatus {
		return handler.NewError(http.StatusBadRequest, "invalid unit status, must be one of 'Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed'")
	}

	unitType, isValidUnitType := enum.ParseUnitType(request.Type)
	if !isValidUnitType {
		return handler.NewError(http.StatusBadRequest, "invalid unit type, must be 'cabin' or 'capsule'")
	}

	unit := domain.Units{
		Name:   request.Name,
		Status: status,
		Type:   unitType,
	}

	errSave := u.unitRepository.Create(unit)
	if errSave != nil {
		return handler.NewError(http.StatusInternalServerError, errSave.Error())
	}

	return nil
}
