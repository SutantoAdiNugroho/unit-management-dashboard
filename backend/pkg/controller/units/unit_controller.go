package units

import (
	"net/http"
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/dto"
	"unit-management-be/pkg/model/dto/request"
	unitService "unit-management-be/pkg/service/units"
	"unit-management-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UnitController struct {
	unitService unitService.UnitService
}

func NewUnitController(unitService unitService.UnitService) *UnitController {
	return &UnitController{unitService: unitService}
}

func SetupUnitRoutes(r *gin.RouterGroup, uc *UnitController) {
	unitGroup := r.Group("/unit")
	unitGroup.POST("", uc.CreateUnit)
}

func (uc *UnitController) CreateUnit(c *gin.Context) {
	var body request.CreateUnitDto
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(handler.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	if utils.IsEmptyString(body.Name) {
		c.Error(handler.NewError(http.StatusBadRequest, "unit name is required"))
		return
	}

	if utils.IsEmptyString(body.Status) {
		c.Error(handler.NewError(http.StatusBadRequest, "unit status is required"))
		return
	}

	if utils.IsEmptyString(body.Type) {
		c.Error(handler.NewError(http.StatusBadRequest, "unit type is required"))
		return
	}

	errUnit := uc.unitService.CreateUnit(body)
	if errUnit != nil {
		c.Error(handler.NewError(errUnit.Code, errUnit.Message))
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse(true, "OK", nil))
}
