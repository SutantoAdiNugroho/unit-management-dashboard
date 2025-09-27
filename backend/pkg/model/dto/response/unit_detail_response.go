package response

import (
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/domain/enum"

	"github.com/google/uuid"
)

type UnitDetailResponse struct {
	ID     uuid.UUID       `json:"id"`
	Name   string          `json:"name"`
	Type   enum.UnitType   `json:"type"`
	Status enum.UnitStatus `json:"status"`
}

func BuildUnitDetailResponseFromUnit(unit domain.Units) UnitDetailResponse {
	return UnitDetailResponse{
		ID:     unit.ID,
		Name:   unit.Name,
		Type:   unit.Type,
		Status: unit.Status,
	}
}
