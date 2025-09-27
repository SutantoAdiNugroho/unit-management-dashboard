package units

import (
	"net/http"
	"testing"

	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/model/domain"
	"unit-management-be/pkg/model/domain/enum"
	"unit-management-be/pkg/model/dto/request"
	"unit-management-be/pkg/model/dto/response"
	unitrepository "unit-management-be/pkg/repository/units"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUnitRepository of unit repository
type MockUnitRepository struct {
	mock.Mock
}

func (m *MockUnitRepository) Create(unit domain.Units) (domain.Units, error) {
	args := m.Called(unit)
	return args.Get(0).(domain.Units), args.Error(1)
}

func (m *MockUnitRepository) GetByID(id string) (domain.Units, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Units), args.Error(1)
}

func (m *MockUnitRepository) Delete(unit domain.Units) error {
	args := m.Called(unit)
	return args.Error(0)
}

func (m *MockUnitRepository) Update(unit domain.Units) error {
	args := m.Called(unit)
	return args.Error(0)
}

func (m *MockUnitRepository) FindAll(status, unitType, name string, page, size int) ([]response.UnitDetailResponse, int64, error) {
	args := m.Called(status, unitType, name, page, size)
	return args.Get(0).([]response.UnitDetailResponse), args.Get(1).(int64), args.Error(2)
}

var _ unitrepository.UnitRepository = &MockUnitRepository{}

// initialization service and unit repository
func setupTest(t *testing.T) (*MockUnitRepository, UnitService) {
	mockRepo := new(MockUnitRepository)
	unitService := NewUnitService(mockRepo)
	return mockRepo, unitService
}

func TestCreateUnit(t *testing.T) {
	t.Run("Positive Case: Create unit successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		req := request.CreateUnitDto{
			Name:   "Unit Test",
			Status: "Available",
			Type:   "capsule",
		}
		expectedUnit := domain.Units{
			ID:     uuid.New(),
			Name:   "Unit Test",
			Status: enum.Available,
			Type:   enum.Capsule,
		}

		mockRepo.On("Create", mock.Anything).Return(expectedUnit, nil).Once()

		result, err := unitService.CreateUnit(req)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.Equal(t, expectedUnit.Name, result.Name)
		assert.Equal(t, expectedUnit.Status, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Invalid status", func(t *testing.T) {
		_, unitService := setupTest(t)
		req := request.CreateUnitDto{
			Name:   "Unit Test",
			Status: "InvalidStatus",
			Type:   "cabin",
		}
		result, err := unitService.CreateUnit(req)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid unit status, must be one of 'Available', 'Occupied', 'Cleaning In Progress', 'Maintenance Needed'", err.Message)
	})

	t.Run("Negative Case: Invalid unit type", func(t *testing.T) {
		_, unitService := setupTest(t)
		req := request.CreateUnitDto{
			Name:   "Unit Test",
			Status: "Available",
			Type:   "invalid_type",
		}
		result, err := unitService.CreateUnit(req)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "invalid unit type, must be 'cabin' or 'capsule'", err.Message)
	})

	t.Run("Negative Case: Repository returns error", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		req := request.CreateUnitDto{
			Name:   "Unit Test",
			Status: "Available",
			Type:   "capsule",
		}
		expectedUnit := domain.Units{}
		expectedErr := handler.NewError(http.StatusInternalServerError, "database error")
		mockRepo.On("Create", mock.Anything).Return(expectedUnit, expectedErr).Once()
		result, err := unitService.CreateUnit(req)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindByID(t *testing.T) {
	t.Run("Positive Case: Find unit by ID successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		expectedUnit := domain.Units{ID: uuid.MustParse(id)}

		mockRepo.On("GetByID", id).Return(expectedUnit, nil).Once()

		result, err := unitService.FindByID(id)

		assert.Nil(t, err)
		assert.Equal(t, expectedUnit.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Unit not found", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		expectedUnit := domain.Units{}

		mockRepo.On("GetByID", id).Return(expectedUnit, gorm.ErrRecordNotFound).Once()

		_, err := unitService.FindByID(id)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Repository returns other error", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		expectedUnit := domain.Units{}
		expectedErr := gorm.ErrInvalidDB

		mockRepo.On("GetByID", id).Return(expectedUnit, expectedErr).Once()

		_, err := unitService.FindByID(id)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetDetailByID(t *testing.T) {
	t.Run("Positive Case: Get unit detail by ID successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		unit := domain.Units{
			ID:     uuid.MustParse(id),
			Name:   "Unit 1",
			Status: enum.Available,
			Type:   enum.Cabin,
		}
		expectedResponse := response.BuildUnitDetailResponseFromUnit(unit)

		mockRepo.On("GetByID", id).Return(unit, nil).Once()

		result, err := unitService.GetDetailByID(id)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Name, result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Unit not found", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		unit := domain.Units{}

		mockRepo.On("GetByID", id).Return(unit, gorm.ErrRecordNotFound).Once()

		result, err := unitService.GetDetailByID(id)

		assert.Equal(t, response.UnitDetailResponse{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteByID(t *testing.T) {
	t.Run("Positive Case: Delete unit by ID successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		unit := domain.Units{ID: uuid.MustParse(id)}

		mockRepo.On("GetByID", id).Return(unit, nil).Once()
		mockRepo.On("Delete", unit).Return(nil).Once()

		err := unitService.DeleteByID(id)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Unit to delete not found", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		unit := domain.Units{}

		mockRepo.On("GetByID", id).Return(unit, gorm.ErrRecordNotFound).Once()

		err := unitService.DeleteByID(id)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Repository returns error on delete", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		unit := domain.Units{ID: uuid.MustParse(id)}
		expectedErr := gorm.ErrInvalidDB

		mockRepo.On("GetByID", id).Return(unit, nil).Once()
		mockRepo.On("Delete", unit).Return(expectedErr).Once()

		err := unitService.DeleteByID(id)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindUnits(t *testing.T) {
	t.Run("Positive Case: Find units successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		status := "Available"
		unitType := "cabin"
		name := "Unit 1"
		page := 1
		size := 10
		unitsData := []response.UnitDetailResponse{{ID: uuid.New()}, {ID: uuid.New()}}
		totalUnit := int64(2)

		mockRepo.On("FindAll", status, unitType, name, page, size).Return(unitsData, totalUnit, nil).Once()

		result, err := unitService.FindUnits(status, unitType, name, page, size)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, page, result.Pagination.Page)
		assert.Equal(t, size, result.Pagination.Size)
		assert.Equal(t, int(totalUnit), result.Pagination.Total)
		assert.Len(t, result.Content.([]response.UnitDetailResponse), 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Repository returns error", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		status := "Available"
		unitType := "cabin"
		name := ""
		page := 1
		size := 10
		expectedErr := gorm.ErrInvalidDB

		mockRepo.On("FindAll", status, unitType, name, page, size).Return([]response.UnitDetailResponse{}, int64(0), expectedErr).Once()

		result, err := unitService.FindUnits(status, unitType, name, page, size)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Positive Case: Update unit successfully", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		oldUnit := domain.Units{
			ID:     uuid.MustParse(id),
			Name:   "Old Name",
			Status: enum.Available,
			Type:   enum.Capsule,
		}
		updateReq := request.UpdateUnitDto{
			request.CreateUnitDto{
				Name:   "New Name",
				Status: "Cleaning In Progress",
				Type:   "cabin",
			},
		}

		mockRepo.On("GetByID", id).Return(oldUnit, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("domain.Units")).Return(nil).Run(func(args mock.Arguments) {
			argUnit := args.Get(0).(domain.Units)
			assert.Equal(t, updateReq.Name, argUnit.Name)
			assert.Equal(t, enum.CleaningInProgress, argUnit.Status)
			assert.Equal(t, enum.Cabin, argUnit.Type)
			assert.NotZero(t, argUnit.LastUpdated)
		}).Once()

		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateReq.Name, result.Name)
		assert.Equal(t, enum.CleaningInProgress, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Unit to update not found", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		updateReq := request.UpdateUnitDto{}
		unit := domain.Units{}

		mockRepo.On("GetByID", id).Return(unit, gorm.ErrRecordNotFound).Once()

		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Invalid new status", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		oldUnit := domain.Units{ID: uuid.MustParse(id)}
		updateReq := request.UpdateUnitDto{
			request.CreateUnitDto{Status: "InvalidStatus"},
		}

		mockRepo.On("GetByID", id).Return(oldUnit, nil).Once()

		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Invalid unit type", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		oldUnit := domain.Units{ID: uuid.MustParse(id)}
		updateReq := request.UpdateUnitDto{
			request.CreateUnitDto{Status: "Available",
				Type: "invalid_type"},
		}
		mockRepo.On("GetByID", id).Return(oldUnit, nil).Once()
		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Cannot change from 'Occupied' to 'Available'", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		oldUnit := domain.Units{
			ID:     uuid.MustParse(id),
			Status: enum.Occupied,
			Type:   enum.Cabin,
		}
		updateReq := request.UpdateUnitDto{
			request.CreateUnitDto{Status: "Available", Type: "cabin"},
		}
		mockRepo.On("GetByID", id).Return(oldUnit, nil).Once()
		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "unit cannot go directly from occupied to available", err.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Negative Case: Repository returns error on update", func(t *testing.T) {
		mockRepo, unitService := setupTest(t)
		id := uuid.New().String()
		oldUnit := domain.Units{ID: uuid.MustParse(id), Status: enum.Available, Type: enum.Cabin}
		updateReq := request.UpdateUnitDto{request.CreateUnitDto{Status: "Cleaning In Progress", Type: "cabin"}}
		expectedErr := gorm.ErrInvalidDB

		mockRepo.On("GetByID", id).Return(oldUnit, nil).Once()
		mockRepo.On("Update", mock.Anything).Return(expectedErr).Once()

		result, err := unitService.Update(id, updateReq)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockRepo.AssertExpectations(t)
	})
}
