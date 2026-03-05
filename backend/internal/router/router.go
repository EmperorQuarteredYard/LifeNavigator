package router

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	userCtl *controller.UserController,
	inviteCtl *controller.InviteController,
	projectCtl *controller.ProjectController,
	taskCtl *controller.TaskController,
) *gin.Engine {
	r := gin.Default()

	// 公开路由（无需认证）
	public := r.Group("/api")
	{
		public.POST("/v0/register", userCtl.Register)
		public.POST("/v0/login", userCtl.Login)
		public.POST("/v0/refresh", userCtl.Refresh)
		public.GET("/v0/users/:id", userCtl.GetUser) // 公开获取用户基本信息
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(jwt.JWTAuthMiddleware()) // 将用户信息存入 context
	{
		// 用户
		auth.GET("/v0/user/profile", userCtl.Profile)

		// 邀请码
		auth.POST("/v0/invite-codes", inviteCtl.CreateInviteCode)
		auth.GET("/v0/invite-codes/:token", inviteCtl.GetInviteCode)
		auth.GET("/v0/users/:id/invite-codes", inviteCtl.ListUserInviteCodes)

		// 项目
		auth.POST("/v0/projects", projectCtl.CreateProject)
		auth.GET("/v0/projects/:id", projectCtl.GetProject)
		auth.GET("/v0/projects", projectCtl.ListProjects)
		auth.PUT("/v0/projects/:id", projectCtl.UpdateProject)
		auth.DELETE("/v0/projects/:id", projectCtl.DeleteProject)

		// 项目预算
		auth.POST("/v0/projects/:id/budgets", projectCtl.AddProjectBudget)
		auth.PUT("/v0/projects/budgets/:budgetId", projectCtl.UpdateProjectBudget)
		auth.DELETE("/v0/projects/budgets/:budgetId", projectCtl.DeleteProjectBudget)
		auth.GET("/v0/projects/:id/budgets/summary", projectCtl.GetProjectBudgetSummary)
		auth.GET("/v0/projects/:id/tasks/budgets/summary", projectCtl.GetTaskBudgetSummary)

		// 任务
		auth.POST("/v0/tasks", taskCtl.CreateTask)
		auth.GET("/v0/tasks/:id", taskCtl.GetTask)
		auth.GET("/v0/projects/:id/tasks", taskCtl.ListProjectTasks)
		auth.GET("/v0/tasks", taskCtl.ListUserTasks)
		auth.PUT("/v0/tasks/:id", taskCtl.UpdateTask)
		auth.DELETE("/v0/tasks/:id", taskCtl.DeleteTask)
		auth.PUT("/v0/projects/:id/tasks/:id", taskCtl.UpdateTask)

		// 任务查询
		auth.GET("/v0/projects/:id/tasks/status/:status", taskCtl.GetTasksByStatus)
		auth.GET("/v0/projects/:id/tasks/period", taskCtl.GetTasksByTimePeriod)

		// 任务预算
		auth.POST("/v0/tasks/:id/budgets", taskCtl.AddTaskBudget)
		auth.PUT("/v0/tasks/budgets/:budgetId", taskCtl.UpdateTaskBudget)
		auth.DELETE("/v0/tasks/budgets/:budgetId", taskCtl.DeleteTaskBudget)
		auth.GET("/v0/tasks/:id/budgets", taskCtl.GetTaskBudgets)
	}

	return r
}
