package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"unit-management-be/pkg/model/dto"
	"unit-management-be/pkg/model/dto/request"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const baseURL = "http://localhost:5000/api"

type UnitTestSuite struct {
	suite.Suite
	unitIDs []string
}

func (s *UnitTestSuite) SetupSuite() {
	s.unitIDs = make([]string, 5)

	// create 5 initial units
	for i := 0; i < 5; i++ {
		createDto := request.CreateUnitDto{
			Name:   fmt.Sprintf("Unit Test %d", i+1),
			Type:   "capsule",
			Status: "Available",
		}
		body, err := json.Marshal(createDto)
		s.Require().NoError(err)

		resp, err := http.Post(baseURL+"/unit", "application/json", bytes.NewBuffer(body))
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().Equal(http.StatusCreated, resp.StatusCode)

		var res dto.Response
		err = json.NewDecoder(resp.Body).Decode(&res)
		s.Require().NoError(err)
		s.Require().True(res.Success)

		// extract ID from response
		dataBytes, err := json.Marshal(res.Data)
		s.Require().NoError(err)

		var unit struct {
			ID string `json:"id"`
		}
		err = json.Unmarshal(dataBytes, &unit)
		s.Require().NoError(err)

		s.unitIDs[i] = unit.ID
	}
}

func (s *UnitTestSuite) TestCreateUnit() {
	s.Run("Positive Case: Should create a new unit successfully", func() {
		createDto := request.CreateUnitDto{
			Name:   "New Test Unit",
			Type:   "capsule",
			Status: "Available",
		}
		body, err := json.Marshal(createDto)
		s.NoError(err)

		resp, err := http.Post(baseURL+"/unit", "application/json", bytes.NewBuffer(body))
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusCreated, resp.StatusCode)
	})

	s.Run("Negative Case: Should return 400 for missing name", func() {
		createDto := request.CreateUnitDto{
			Name:   "",
			Type:   "capsule",
			Status: "Available",
		}
		body, err := json.Marshal(createDto)
		s.NoError(err)

		resp, err := http.Post(baseURL+"/unit", "application/json", bytes.NewBuffer(body))
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *UnitTestSuite) TestGetUnits() {
	s.Run("Positive Case: Should get all units with pagination", func() {
		resp, err := http.Get(baseURL + "/unit?page=1&size=2")
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)

		var res struct {
			Data dto.PaginationResponse `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&res)
		s.NoError(err)
		s.Len(res.Data.Content.([]interface{}), 2)
	})

	s.Run("Negative Case: Should return 400 for invalid page parameter", func() {
		resp, err := http.Get(baseURL + "/unit?page=abc")
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *UnitTestSuite) TestGetDetailUnitByID() {
	s.Run("Positive Case: Should get unit detail by ID", func() {
		unitID := s.unitIDs[0]
		resp, err := http.Get(baseURL + "/unit/" + unitID)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)
	})

	s.Run("Negative Case: Should return 404 for non-existent unit ID", func() {
		nonExistentID := uuid.New().String()
		resp, err := http.Get(baseURL + "/unit/" + nonExistentID)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *UnitTestSuite) TestUpdateUnit() {
	s.Run("Positive Case: Should update a unit successfully", func() {
		unitID := s.unitIDs[1]
		updateDto := request.UpdateUnitDto{
			CreateUnitDto: request.CreateUnitDto{
				Name:   "Updated Unit Name",
				Type:   "capsule",
				Status: "Maintenance Needed",
			},
		}
		body, err := json.Marshal(updateDto)
		s.NoError(err)

		req, err := http.NewRequest(http.MethodPut, baseURL+"/unit/"+unitID, bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)
	})

	s.Run("Negative Case: Should return 400 for direct status change from Occupied to Available", func() {
		unitID := s.unitIDs[2]

		// first set status to Occupied
		occupiedDto := request.UpdateUnitDto{CreateUnitDto: request.CreateUnitDto{Name: "Unit B", Type: "cabin", Status: "Occupied"}}
		body, err := json.Marshal(occupiedDto)
		s.NoError(err)
		req, err := http.NewRequest(http.MethodPut, baseURL+"/unit/"+unitID, bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(req)
		s.NoError(err)

		// then try to set it to Available directly
		availableDto := request.UpdateUnitDto{CreateUnitDto: request.CreateUnitDto{Name: "Unit B", Type: "cabin", Status: "Available"}}
		body, err = json.Marshal(availableDto)
		s.NoError(err)
		req, err = http.NewRequest(http.MethodPut, baseURL+"/unit/"+unitID, bytes.NewBuffer(body))
		s.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})
}

func (s *UnitTestSuite) TestDeleteUnit() {
	s.Run("Positive Case: Should delete a unit successfully", func() {
		unitIDToDelete := s.unitIDs[3]
		req, err := http.NewRequest(http.MethodDelete, baseURL+"/unit/"+unitIDToDelete, nil)
		s.NoError(err)

		resp, err := http.DefaultClient.Do(req)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusOK, resp.StatusCode)

		// verify unit is no longer found
		resp, err = http.Get(baseURL + "/unit/" + unitIDToDelete)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})

	s.Run("Negative Case: Should return 404 for non-existent unit ID", func() {
		nonExistentID := uuid.New().String()
		req, err := http.NewRequest(http.MethodDelete, baseURL+"/unit/"+nonExistentID, nil)
		s.NoError(err)

		resp, err := http.DefaultClient.Do(req)
		s.NoError(err)
		defer resp.Body.Close()
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
