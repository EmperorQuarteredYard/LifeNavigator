package database

import (
	"LifeNavigator/internal/models"
	"encoding/json"
	"fmt"
	"log"
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
	file, err := os.Open("config/database.json")
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
		//config, err := getConfig()
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}
		//if config == nil {
		//	fmt.Println("config is nil")
		//	return nil
		//}
		username := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		database := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		var err error

		dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("fail to connect to database\ndsn:"+dsn+"\nreason:", err)
			return nil
		}

		log.Println("successfully connect to database")

		err = db.AutoMigrate(
			&models.User{},
			&models.InviteCode{},
			&models.Project{},
			&models.ProjectBudget{},
			&models.Task{},
			&models.TaskPayment{},
			&models.TaskDependency{},
			&models.Account{},
		)

		if err != nil {
			fmt.Println("Fail to migrate model\nreason:", err)
		} else {
			fmt.Println("successfully migrate model")
		}
		initialized = true
	}
	return db
}
