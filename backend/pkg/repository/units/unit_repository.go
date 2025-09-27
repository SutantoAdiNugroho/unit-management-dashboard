package units

import "unit-management-be/pkg/model/domain"

type UnitRepository interface {
	Create(unit domain.Units) error
}
