package main

import (
	"LifeNavigator/internal/controller"
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
	//db := database.GetDatabase()
	db := &gorm.DB{}
	if db == nil {
		log.Println("failed to connect to database")
		log.Fatal("failed to connect to database")
	}

	userRepo := repository.NewUserRepository(db)
	inviteCodeRepo := repository.NewInviteCodeRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	projectBudgetRepo := repository.NewProjectBudgetRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	taskBudgetRepo := repository.NewTaskPaymentRepository(db)
	accountRepo := repository.NewAccountRepository(db)

	transactor := repository.NewTransactor(db)

	userService := service.NewUserService(userRepo)
	inviteCodeService := service.NewInviteCodeService(inviteCodeRepo)
	inviteUserService := service.NewInviteUserService(userService, inviteCodeService)
	projectService := service.NewProjectService(transactor, projectRepo, projectBudgetRepo, taskBudgetRepo, taskRepo)
	taskService := service.NewTaskService(transactor, taskRepo, taskBudgetRepo, projectRepo)
	accountService := service.NewAccountService(accountRepo, taskRepo, taskBudgetRepo, transactor)
	aiFeatureService := service.NewAIFeatureService(transactor)

	userCtl := controller.NewUserController(userService, inviteUserService)
	inviteCtl := controller.NewInviteController(inviteCodeService, inviteUserService, userService)
	projectCtl := controller.NewProjectController(projectService)
	taskCtl := controller.NewTaskController(taskService)
	accountCtl := controller.NewAccountController(accountService)
	aiFeatureCtl := controller.NewAIFeatureController(aiFeatureService, accountService)
	userService.Register(&models.User{
		Username: "Administrator",
		Password: "Administrator",
		Role:     roles.Administrator,
		Nickname: "Administrator",
	})

	r := router.InitRouter(userCtl, inviteCtl, projectCtl, taskCtl, accountCtl, aiFeatureCtl) //,
	log.Println("Listening on port :5083")
	//go func() {
	if err := r.Run(":5083"); err != nil {
		log.Fatal("failed to start server: ", err)
	}

	//}()
	//cliCTL := clicontroller.NewCLIController(userService, accountService, projectService, taskService, inviteCodeService, inviteUserService)
	//cliCTL.Run()
}
