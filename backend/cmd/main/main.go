package main

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/internal/database"
	"LifeNavigator/internal/repository"
	"LifeNavigator/internal/router"
	"LifeNavigator/internal/service"
	"log"
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
	inviteCodeRepo := repository.NewInviteCodeRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	projectBudgetRepo := repository.NewProjectBudgetRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	taskBudgetRepo := repository.NewTaskBudgetRepository(db)

	transactor := repository.NewTransactor(db)

	userService := service.NewUserService(userRepo)
	inviteCodeService := service.NewInviteCodeService(inviteCodeRepo)
	inviteUserService := service.NewInviteUserService(userService, inviteCodeService)
	projectService := service.NewProjectService(transactor, projectRepo, projectBudgetRepo, taskBudgetRepo, taskRepo)
	taskService := service.NewTaskService(transactor, taskRepo, taskBudgetRepo, projectRepo)

	userCtl := controller.NewUserController(userService, inviteUserService)
	inviteCtl := controller.NewInviteController(inviteCodeService, inviteUserService, userService)
	projectCtl := controller.NewProjectController(projectService)
	taskCtl := controller.NewTaskController(taskService)

	r := router.InitRouter(userCtl, inviteCtl, projectCtl, taskCtl)
	log.Println("Listening on port :5083")
	if err := r.Run(":5083"); err != nil {
		log.Println("failed to start server: ", err)
		log.Fatal("failed to start server: ", err)
	}
}
