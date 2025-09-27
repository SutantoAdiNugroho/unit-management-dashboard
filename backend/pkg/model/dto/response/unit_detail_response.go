package response

import (
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/domain/enum"

	"github.com/google/uuid"
)

type UnitDetailResponse struct {
	ID     uuid.UUID
	Name   string
	Type   enum.UnitType
	Status enum.UnitStatus
}

func BuildUnitDetailResponseFromUnit(unit domain.Units) UnitDetailResponse {
	return UnitDetailResponse{
		ID:     unit.ID,
		Name:   unit.Name,
		Type:   unit.Type,
		Status: unit.Status,
	}
}
