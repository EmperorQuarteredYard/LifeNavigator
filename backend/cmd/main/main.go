package main

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/internal/database"
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/internal/router"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/roles"
	"log"

	"gorm.io/gorm"
)

func main() {
	ServeByRealDatabase()
}

func ServeByRealDatabase() {
	db := database.GetDatabase()
	if db == nil {
		log.Println("failed to connect to database")
		log.Fatal("failed to connect to database")
	}

	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	projectBudgetRepo := repository.NewProjectBudgetRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	taskPaymentRepo := repository.NewTaskPaymentRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	transactor := repository.NewTransactor(db)

	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(transactor, projectRepo, projectBudgetRepo, taskPaymentRepo, taskRepo)
	taskService := service.NewTaskService(transactor, taskRepo, taskPaymentRepo, projectRepo, projectService)
	accountService := service.NewAccountService(accountRepo, taskRepo, taskPaymentRepo, transactor)

	authCtl := controller.NewAuthController(userService)
	projectCtl := controller.NewProjectController(projectService)
	taskCtl := controller.NewTaskController(taskService)
	accountCtl := controller.NewAccountController(accountService)

	initAdministrator(db, userService)

	r := router.InitRouter(authCtl, projectCtl, taskCtl, accountCtl)
	log.Println("Listening on port :5083")
	if err := r.Run(":5083"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}

func initAdministrator(db *gorm.DB, userService service.UserService) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", roles.Administrator).Count(&count)
	if count == 0 {
		userService.Register(&dto.RegisterRequest{
			Username: "Administrator",
			Password: "Administrator",
			Nickname: "Administrator",
		})
		log.Println("Administrator account created")
	}
}
