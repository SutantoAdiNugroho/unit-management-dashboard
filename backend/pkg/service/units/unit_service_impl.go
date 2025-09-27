package units

import (
	"net/http"
	"time"
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/domain/enum"
	"unit-management-be/pkg/model/dto"
	"unit-management-be/pkg/model/dto/request"
	"unit-management-be/pkg/model/dto/response"
	unitrepository "unit-management-be/pkg/repository/units"

	"gorm.io/gorm"
)

type UnitServiceImpl struct {
	unitRepository unitrepository.UnitRepository
}

func NewUnitService(unitRepository unitrepository.UnitRepository) UnitService {
	return &UnitServiceImpl{unitRepository: unitRepository}
}
func (u *UnitServiceImpl) CreateUnit(request request.CreateUnitDto) (*domain.Units, *handler.CustomError) {
	status, isValidStatus := enum.ParseUnitStatus(request.Status)
	if !isValidStatus {
		return nil, handler.NewError(http.StatusBadRequest, "invalid unit status, must be one of 'Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed'")
	}

	unitType, isValidUnitType := enum.ParseUnitType(request.Type)
	if !isValidUnitType {
		return nil, handler.NewError(http.StatusBadRequest, "invalid unit type, must be 'cabin' or 'capsule'")
	}

	unit := domain.Units{
		Name:   request.Name,
		Status: status,
		Type:   unitType,
	}

	createdUnit, errSave := u.unitRepository.Create(unit)
	if errSave != nil {
		return nil, handler.NewError(http.StatusInternalServerError, errSave.Error())
	}

	return &createdUnit, nil
}

func (u *UnitServiceImpl) FindByID(id string) (domain.Units, *handler.CustomError) {
	var unit domain.Units

	unit, err := u.unitRepository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return unit, handler.NewError(http.StatusNotFound, "unit with that id was not found")
		}
		return unit, handler.NewError(http.StatusInternalServerError, err.Error())
	}

	return unit, nil
}

func (u *UnitServiceImpl) GetDetailByID(id string) (response.UnitDetailResponse, *handler.CustomError) {
	var responseUnit response.UnitDetailResponse

	unit, err := u.FindByID(id)
	if err != nil {
		return responseUnit, handler.NewError(err.Code, err.Message)
	}

	return response.BuildUnitDetailResponseFromUnit(unit), nil
}

func (u *UnitServiceImpl) DeleteByID(id string) *handler.CustomError {
	unit, err := u.FindByID(id)
	if err != nil {
		return handler.NewError(err.Code, err.Message)
	}

	errDelete := u.unitRepository.Delete(unit)
	if errDelete != nil {
		return handler.NewError(http.StatusInternalServerError, errDelete.Error())
	}

	return nil
}

func (u *UnitServiceImpl) FindUnits(status, unitType, name string, page, size int) (*dto.PaginationResponse, *handler.CustomError) {
	units, totalUnit, err := u.unitRepository.FindAll(status, unitType, name, page, size)
	if err != nil {
		return nil, handler.NewError(http.StatusInternalServerError, err.Error())
	}

	return dto.NewPaginationResponse(page, size, int(totalUnit), units), nil
}

func (u *UnitServiceImpl) Update(id string, request request.UpdateUnitDto) (*domain.Units, *handler.CustomError) {
	unit, err := u.FindByID(id)
	if err != nil {
		return nil, handler.NewError(err.Code, err.Message)
	}

	unitType, isValidUnitType := enum.ParseUnitType(request.Type)
	if !isValidUnitType {
		return nil, handler.NewError(http.StatusBadRequest, "invalid unit type, must be 'cabin' or 'capsule'")
	}

	newStatus, isValidStatus := enum.ParseUnitStatus(request.Status)
	if !isValidStatus {
		return nil, handler.NewError(http.StatusBadRequest, "invalid unit status, must be one of 'Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed'")
	}

	if unit.Status == enum.Occupied && newStatus == enum.Available {
		return nil, handler.NewError(http.StatusBadRequest, "unit cannot go directly from occupied to available")
	}

	unit.Name = request.Name
	unit.Type = unitType
	unit.Status = newStatus
	unit.LastUpdated = time.Now()

	errUpdate := u.unitRepository.Update(unit)
	if errUpdate != nil {
		return nil, handler.NewError(http.StatusInternalServerError, errUpdate.Error())
	}

	return &unit, nil
}
