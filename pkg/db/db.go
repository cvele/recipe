package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/cvele/recipe/pkg/config"
	"github.com/cvele/recipe/pkg/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var connStr string
	if cfg.DBType == "mysql" {
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	} else if cfg.DBType == "postgres" {
		connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)
	} else {
		log.Fatalf("Unsupported database type: %s", cfg.DBType)
	}

	db, err := gorm.Open(cfg.DBType, connStr)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Recipe{}, &models.Ingredient{})

	return db, nil
}
