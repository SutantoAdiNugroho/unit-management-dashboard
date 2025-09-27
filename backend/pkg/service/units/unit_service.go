package units

import (
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/dto/request"
)

type UnitService interface {
	CreateUnit(request request.CreateUnitDto) *handler.CustomError
}
