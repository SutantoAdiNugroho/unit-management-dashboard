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
