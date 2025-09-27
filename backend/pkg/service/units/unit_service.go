package units

import (
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/dto"
	"unit-management-be/pkg/model/dto/request"
	"unit-management-be/pkg/model/dto/response"
)

type UnitService interface {
	CreateUnit(request request.CreateUnitDto) (*domain.Units, *handler.CustomError)
	GetDetailByID(id string) (response.UnitDetailResponse, *handler.CustomError)
	DeleteByID(id string) *handler.CustomError
	FindByID(id string) (domain.Units, *handler.CustomError)
	FindUnits(status, unitType, name string, page, size int) (*dto.PaginationResponse, *handler.CustomError)
	Update(id string, request request.UpdateUnitDto) (*domain.Units, *handler.CustomError)
}
