package routes

import (
	"fmt"
	"log"
	"os"
	"unit-management-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()

	port := os.Getenv("PORT")
	if utils.IsEmptyString(port) {
		port = "5000"
	}

	log.Printf("application running on port : %s", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("failed to run backend")
	}
}
