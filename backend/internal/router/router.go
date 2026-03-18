package router

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	authCtl *controller.AuthController,
	projectCtl *controller.ProjectController,
	taskCtl *controller.TaskController,
	accountCtl *controller.AccountController,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtl.Register)
			auth.POST("/login", authCtl.Login)
			auth.POST("/refresh", authCtl.RefreshToken)
		}

		users := api.Group("/users")
		{
			users.GET("/:id", authCtl.GetUserByID)
		}

		protected := api.Group("")
		protected.Use(jwt.JWTAuthMiddleware())
		{
			protected.GET("/auth/profile", authCtl.GetProfile)

			projects := protected.Group("/projects")
			{
				projects.POST("", projectCtl.CreateProject)
				projects.GET("", projectCtl.ListProjects)
				projects.GET("/:id", projectCtl.GetProject)
				projects.PUT("/:id", projectCtl.UpdateProject)
				projects.DELETE("/:id", projectCtl.DeleteProject)
				projects.POST("/:id/budgets", projectCtl.AddBudget)
				projects.PUT("/budgets/:budgetId", projectCtl.UpdateBudget)
				projects.DELETE("/budgets/:budgetId", projectCtl.DeleteBudget)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.POST("", taskCtl.CreateTask)
				tasks.GET("", taskCtl.ListTasks)
				tasks.GET("/:id", taskCtl.GetTask)
				tasks.PUT("/:id", taskCtl.UpdateTask)
				tasks.PATCH("/:id/status", taskCtl.UpdateTaskStatus)
				tasks.DELETE("/:id", taskCtl.DeleteTask)
				tasks.POST("/:id/finish", taskCtl.FinishTask)
				tasks.GET("/:id/prerequisites", taskCtl.GetPrerequisites)
				tasks.POST("/:id/prerequisites", taskCtl.SetPrerequisites)
				tasks.DELETE("/:id/prerequisites", taskCtl.UnsetPrerequisites)
				tasks.GET("/:id/postrequisites", taskCtl.GetPostrequisites)
				tasks.POST("/:id/payments", taskCtl.SetPayment)
				tasks.PUT("/payments/:id", taskCtl.UpdatePayment)
				tasks.DELETE("/payments/:id", taskCtl.DeletePayment)
				tasks.GET("/:id/payments", taskCtl.GetPayments)
			}

			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountCtl.CreateAccount)
				accounts.GET("", accountCtl.ListAccounts)
				accounts.GET("/:id", accountCtl.GetAccount)
				accounts.POST("/:id/adjust", accountCtl.AdjustBalance)
				accounts.DELETE("/:id", accountCtl.DeleteAccount)
				accounts.GET("/:id/tasks", accountCtl.GetAccountTasks)
			}
		}
	}

	return r
}

func InitRouterWithKanban(
	authCtl *controller.AuthController,
	projectCtl *controller.ProjectController,
	taskCtl *controller.TaskController,
	accountCtl *controller.AccountController,
	kanbanCtl *controller.KanbanController,
) *gin.Engine {
	r := InitRouter(authCtl, projectCtl, taskCtl, accountCtl)

	api := r.Group("/api")
	protected := api.Group("")
	protected.Use(jwt.JWTAuthMiddleware())
	{
		kanbans := protected.Group("/kanbans")
		{
			kanbans.POST("", kanbanCtl.CreateKanban)
			kanbans.GET("", kanbanCtl.ListKanbans)
			kanbans.GET("/default/tasks", kanbanCtl.GetDefaultKanbanTasks)
			kanbans.GET("/:id", kanbanCtl.GetKanban)
			kanbans.PUT("/:id", kanbanCtl.UpdateKanban)
			kanbans.DELETE("/:id", kanbanCtl.DeleteKanban)
			kanbans.GET("/:id/tasks", kanbanCtl.GetKanbanTasks)
			kanbans.POST("/:id/default", kanbanCtl.SetDefaultKanban)
		}
	}

	return r
}

func InitRouterWithAll(
	authCtl *controller.AuthController,
	projectCtl *controller.ProjectController,
	taskCtl *controller.TaskController,
	accountCtl *controller.AccountController,
	kanbanCtl *controller.KanbanController,
	inviteCtl *controller.InviteController,
	userCtl *controller.UserController,
	aiCtl *controller.AIFeatureController,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authCtl.Register)
			auth.POST("/login", authCtl.Login)
			auth.POST("/refresh", authCtl.RefreshToken)
		}

		users := api.Group("/users")
		{
			users.GET("/:id", authCtl.GetUserByID)
		}

		inviteCodes := api.Group("/invite-codes")
		inviteCodes.Use(jwt.JWTAuthMiddleware())
		{
			inviteCodes.POST("", inviteCtl.CreateInviteCode)
			inviteCodes.GET("/:token", inviteCtl.GetInviteCode)
		}

		usersProtected := api.Group("/users")
		usersProtected.Use(jwt.JWTAuthMiddleware())
		{
			usersProtected.GET("/:id/invite-codes", inviteCtl.GetUserInviteCodes)
		}

		protected := api.Group("")
		protected.Use(jwt.JWTAuthMiddleware())
		{
			protected.GET("/auth/profile", authCtl.GetProfile)

			projects := protected.Group("/projects")
			{
				projects.POST("", projectCtl.CreateProject)
				projects.GET("", projectCtl.ListProjects)
				projects.GET("/:id", projectCtl.GetProject)
				projects.PUT("/:id", projectCtl.UpdateProject)
				projects.DELETE("/:id", projectCtl.DeleteProject)
				projects.POST("/:id/budgets", projectCtl.AddBudget)
				projects.PUT("/budgets/:budgetId", projectCtl.UpdateBudget)
				projects.DELETE("/budgets/:budgetId", projectCtl.DeleteBudget)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.POST("", taskCtl.CreateTask)
				tasks.GET("", taskCtl.ListTasks)
				tasks.GET("/:id", taskCtl.GetTask)
				tasks.PUT("/:id", taskCtl.UpdateTask)
				tasks.PATCH("/:id/status", taskCtl.UpdateTaskStatus)
				tasks.DELETE("/:id", taskCtl.DeleteTask)
				tasks.POST("/:id/finish", taskCtl.FinishTask)
				tasks.GET("/:id/prerequisites", taskCtl.GetPrerequisites)
				tasks.POST("/:id/prerequisites", taskCtl.SetPrerequisites)
				tasks.DELETE("/:id/prerequisites", taskCtl.UnsetPrerequisites)
				tasks.GET("/:id/postrequisites", taskCtl.GetPostrequisites)
				tasks.POST("/:id/payments", taskCtl.SetPayment)
				tasks.PUT("/payments/:id", taskCtl.UpdatePayment)
				tasks.DELETE("/payments/:id", taskCtl.DeletePayment)
				tasks.GET("/:id/payments", taskCtl.GetPayments)
			}

			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountCtl.CreateAccount)
				accounts.GET("", accountCtl.ListAccounts)
				accounts.GET("/:id", accountCtl.GetAccount)
				accounts.POST("/:id/adjust", accountCtl.AdjustBalance)
				accounts.DELETE("/:id", accountCtl.DeleteAccount)
				accounts.GET("/:id/tasks", accountCtl.GetAccountTasks)
			}

			kanbans := protected.Group("/kanbans")
			{
				kanbans.POST("", kanbanCtl.CreateKanban)
				kanbans.GET("", kanbanCtl.ListKanbans)
				kanbans.GET("/default/tasks", kanbanCtl.GetDefaultKanbanTasks)
				kanbans.GET("/:id", kanbanCtl.GetKanban)
				kanbans.PUT("/:id", kanbanCtl.UpdateKanban)
				kanbans.DELETE("/:id", kanbanCtl.DeleteKanban)
				kanbans.GET("/:id/tasks", kanbanCtl.GetKanbanTasks)
				kanbans.POST("/:id/default", kanbanCtl.SetDefaultKanban)
			}

			ai := protected.Group("/ai")
			{
				ai.POST("/reduce-project", aiCtl.ReduceProject)
				ai.POST("/summary", aiCtl.Summary)
			}
		}
	}

	return r
}
