package units

import (
	"net/http"
	"strconv"
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
	unitGroup.GET("/:unitId", uc.GetDetailUnitByID)
	unitGroup.DELETE("/:unitId", uc.DeleteUnit)
	unitGroup.GET("", uc.GetUnits)
	unitGroup.PUT("/:unitId", uc.UpdateUnit)
}

// @Summary Create Unit
// @Description Create new unit with name, status and type
// @Tags Units
// @Accept json
// @Produce json
// @Param unit body request.CreateUnitDto true "Unit creation request"
// @Success 201 {object} dto.Response "Unit created successfully"
// @Failure 400 {object} dto.Response "Bad request: Missing required fields"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /unit [post]
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

	unit, errUnit := uc.unitService.CreateUnit(body)
	if errUnit != nil {
		c.Error(handler.NewError(errUnit.Code, errUnit.Message))
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse(true, "OK", unit))
}

// @Summary Get Unit Detail by ID
// @Description Retrieve details of specific unit using its ID
// @Tags Units
// @Produce json
// @Param unitId path string true "Unit ID"
// @Success 200 {object} dto.Response "Successfully retrieved unit detail"
// @Failure 400 {object} dto.Response "Bad request"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /unit/{unitId} [get]
func (uc *UnitController) GetDetailUnitByID(c *gin.Context) {
	unitId := c.Param("unitId")

	unit, err := uc.unitService.GetDetailByID(unitId)
	if err != nil {
		c.Error(handler.NewError(err.Code, err.Message))
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse(true, "OK", unit))
}

// @Summary Delete Unit by ID
// @Description Delete unit using its ID
// @Tags Units
// @Produce json
// @Param unitId path string true "Unit ID"
// @Success 200 {object} dto.Response "Unit successfully deleted"
// @Failure 404 {object} dto.Response "Unit not found"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /unit/{unitId} [delete]
func (uc *UnitController) DeleteUnit(c *gin.Context) {
	unitId := c.Param("unitId")

	err := uc.unitService.DeleteByID(unitId)
	if err != nil {
		c.Error(handler.NewError(err.Code, err.Message))
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse(true, "OK", nil))
}

// @Summary Get List of Units
// @Description Retrieve list of units with optional filtering and pagination
// @Tags Units
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param size query int false "Number of items per page (default 10)"
// @Param name query string false "Filter by unit name"
// @Param status query string false "Filter by unit status (Available, Occupied)"
// @Param type query string false "Filter by unit type (capsule, cabin)"
// @Success 200 {object} dto.Response{data=dto.PaginationResponse} "Successfully retrieved list of units"
// @Failure 400 {object} dto.Response "Bad request (invalid page/size parameter)"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /unit [get]
func (uc *UnitController) GetUnits(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	nameStr := c.DefaultQuery("name", "")
	statusStr := c.DefaultQuery("status", "")
	typeStr := c.DefaultQuery("type", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.Error(handler.NewError(http.StatusBadRequest, "invalid page parameter, must be number"))
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		c.Error(handler.NewError(http.StatusBadRequest, "invalid size parameter, must be number"))
		return
	}

	units, errUnits := uc.unitService.FindUnits(statusStr, typeStr, nameStr, page, size)
	if errUnits != nil {
		c.Error(handler.NewError(errUnits.Code, errUnits.Message))
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse(true, "OK", units))
}

// @Summary Update Unit
// @Description Update existing units name, status or type.
// @Tags Units
// @Accept json
// @Produce json
// @Param unitId path string true "Unit ID"
// @Param unit body request.UpdateUnitDto true "Unit update request"
// @Success 200 {object} dto.Response "Unit successfully updated"
// @Failure 400 {object} dto.Response "Bad request (missing required fields or invalid ID)"
// @Failure 404 {object} dto.Response "Unit not found"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /unit/{unitId} [put]
func (uc *UnitController) UpdateUnit(c *gin.Context) {
	unitId := c.Param("unitId")

	var body request.UpdateUnitDto
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

	unit, errUnit := uc.unitService.Update(unitId, body)
	if errUnit != nil {
		c.Error(handler.NewError(errUnit.Code, errUnit.Message))
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse(true, "OK", unit))
}
