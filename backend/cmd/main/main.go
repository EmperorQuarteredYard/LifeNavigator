package main

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/internal/database"
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/internal/router"
	"LifeNavigator/internal/service"
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
	inviteCodeRepo := repository.NewInviteCodeRepository(db)
	kanbanRepo := repository.NewKanbanRepository(db)

	transactor := repository.NewTransactor(db)

	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(transactor, projectRepo, projectBudgetRepo, taskPaymentRepo, taskRepo)
	taskService := service.NewTaskService(transactor, taskRepo, taskPaymentRepo, projectRepo, accountRepo)
	accountService := service.NewAccountService(accountRepo, taskRepo, taskPaymentRepo, userRepo, transactor)
	aiFeatureService := service.NewAIFeatureService(transactor)
	inviteCodeService := service.NewInviteCodeService(inviteCodeRepo)
	inviteUserService := service.NewInviteUserService(userService, inviteCodeService)
	kanbanSerivice := service.NewKanbanService(kanbanRepo, projectRepo, taskRepo, taskPaymentRepo)

	projectCtl := controller.NewProjectController(projectService)
	taskCtl := controller.NewTaskController(taskService)
	accountCtl := controller.NewAccountController(accountService)
	aIFeatureController := controller.NewAIFeatureController(aiFeatureService, accountService)
	inviteController := controller.NewInviteController(inviteCodeService, inviteUserService, userService)
	kanbanController := controller.NewKanbanController(kanbanSerivice)
	userController := controller.NewUserController(userService, inviteUserService)

	initAdministrator(db, userService)

	r := router.SetupRouter(accountCtl, aIFeatureController, inviteController, kanbanController, projectCtl, taskCtl, userController)
	log.Println("Listening on port :5083")
	if err := r.Run(":5083"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}

func initAdministrator(db *gorm.DB, userService service.UserService) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", roles.Administrator).Count(&count)
	if count == 0 {
		userService.Register(&models.User{
			Username: "Administrator",
			Password: "Administrator",
			Nickname: "Administrator",
		})
		log.Println("Administrator account created")
	}
}
