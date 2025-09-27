package enum

type UnitType string
type UnitStatus string

const (
	Capsule UnitType = "capsule"
	Cabin   UnitType = "cabin"

	Available          UnitStatus = "Available"
	Occupied           UnitStatus = "Occupied"
	CleaningInProgress UnitStatus = "Cleaning In Progress"
	MaintenanceNeeded  UnitStatus = "Maintenance Needed"
)

func ParseUnitType(value string) (UnitType, bool) {
	switch value {
	case string(Capsule):
		return Capsule, true
	case string(Cabin):
		return Cabin, true
	default:
		return "", false
	}
}

func ParseUnitStatus(value string) (UnitStatus, bool) {
	switch value {
	case string(Available):
		return Available, true
	case string(Occupied):
		return Occupied, true
	case string(CleaningInProgress):
		return CleaningInProgress, true
	case string(MaintenanceNeeded):
		return MaintenanceNeeded, true
	default:
		return "", false
	}
}
