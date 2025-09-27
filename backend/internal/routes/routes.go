package routes

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unit-management-be/internal/db"
	"unit-management-be/pkg/handler"
	"unit-management-be/pkg/utils"

	unitcontroller "unit-management-be/pkg/controller/units"
	unitrepository "unit-management-be/pkg/repository/units"
	unitservice "unit-management-be/pkg/service/units"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	r.Use(handler.ErrorHandler())

	allowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	allowMethods := os.Getenv("CORS_ALLOW_METHOD")

	// cors configuration
	config := cors.Config{
		AllowOrigins:     strings.Split(allowOrigins, ","),
		AllowMethods:     strings.Split(allowMethods, ","),
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	db := db.GetDB()
	unitRepository := unitrepository.NewUnitRepository(db)
	unitService := unitservice.NewUnitService(unitRepository)
	unitController := unitcontroller.NewUnitController(unitService)

	api := r.Group("/api")
	unitcontroller.SetupUnitRoutes(api, unitController)

	port := os.Getenv("PORT")
	if utils.IsEmptyString(port) {
		port = "5000"
	}

	log.Printf("application running on port : %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("failed to run backend")
	}
}
