package main

import (
	"log"
	"todo-api/config"
	"todo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	r := gin.Default()

	routes.RegisterRoutes(r, db)

	log.Fatal(r.Run(":8080"))
}
