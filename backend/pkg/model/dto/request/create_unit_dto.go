package request

type CreateUnitDto struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
