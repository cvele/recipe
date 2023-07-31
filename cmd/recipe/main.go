package main

import (
	"github.com/cvele/recipe/pkg/config"
	"github.com/cvele/recipe/pkg/controllers"
	"github.com/cvele/recipe/pkg/db"
	"github.com/cvele/recipe/pkg/repositories"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Error parsing log level: %v", err)
	}
	log.SetLevel(logLevel)

	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	repo := repositories.NewGormRecipeRepository(db)
	controller := controllers.NewRecipeController(repo)

	router := gin.Default()
	controller.RegisterRoutes(router.Group("/api"))

	log.Infof("Starting server on port %s", cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}
