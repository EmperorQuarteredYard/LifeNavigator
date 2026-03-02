package database

import (
	"LifeNavigator/backend/internal/models"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	initialized bool = false
	db          *gorm.DB
)

type databaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     string `json:"port"`
	Host     string `json:"host"`
}

func getConfig() (*databaseConfig, error) {
	file, err := os.Open("backend/config/database.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var newConfig databaseConfig
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&newConfig)
	if err != nil {
		return nil, err
	}
	return &newConfig, nil
}

func GetDatabase() *gorm.DB {
	if !initialized {
		config, err := getConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if config == nil {
			fmt.Println("config is nil")
			return nil
		}
		dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println("fail to connect to database\nreason:", err)
			return nil
		}

		fmt.Println("successfully connect to database")

		err = db.AutoMigrate(&models.User{}, &models.Arrangement{}, &models.Project{})

		if err != nil {
			fmt.Println("Fail to migrate model\nreason:", err)
		} else {
			fmt.Println("successfully migrate model")
		}
		initialized = true
	}
	return db
}
