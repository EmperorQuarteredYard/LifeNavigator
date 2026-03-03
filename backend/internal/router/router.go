package router

import (
	"LifeNavigator/backend/internal/controller"
	"LifeNavigator/backend/pkg/jwt"

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
		public.POST("/register", userCtl.Register)
		public.POST("/login", userCtl.Login)
		public.POST("/refresh", userCtl.Refresh)
		public.GET("/users/:id", userCtl.GetUser) // 公开获取用户基本信息
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(jwt.JWTAuthMiddleware()) // 将用户信息存入 context
	{
		// 用户
		auth.GET("/user/profile", userCtl.Profile)

		// 邀请码
		auth.POST("/invite-codes", inviteCtl.CreateInviteCode)
		auth.GET("/invite-codes/:token", inviteCtl.GetInviteCode)
		auth.GET("/users/:userID/invite-codes", inviteCtl.ListUserInviteCodes)

		// 项目
		auth.POST("/projects", projectCtl.CreateProject)
		auth.GET("/projects/:id", projectCtl.GetProject)
		auth.GET("/projects", projectCtl.ListProjects)
		auth.PUT("/projects/:id", projectCtl.UpdateProject)
		auth.DELETE("/projects/:id", projectCtl.DeleteProject)

		// 项目预算
		auth.POST("/projects/:id/budgets", projectCtl.AddProjectBudget)
		auth.PUT("/projects/budgets/:budgetId", projectCtl.UpdateProjectBudget)
		auth.DELETE("/projects/budgets/:budgetId", projectCtl.DeleteProjectBudget)
		auth.GET("/projects/:id/budgets/summary", projectCtl.GetProjectBudgetSummary)
		auth.GET("/projects/:id/tasks/budgets/summary", projectCtl.GetTaskBudgetSummary)

		// 任务
		auth.POST("/tasks", taskCtl.CreateTask)
		auth.GET("/tasks/:id", taskCtl.GetTask)
		auth.GET("/projects/:projectId/tasks", taskCtl.ListProjectTasks)
		auth.GET("/tasks", taskCtl.ListUserTasks)
		auth.PUT("/tasks/:id", taskCtl.UpdateTask)
		auth.DELETE("/tasks/:id", taskCtl.DeleteTask)

		// 任务查询
		auth.GET("/projects/:projectId/tasks/status/:status", taskCtl.GetTasksByStatus)
		auth.GET("/projects/:projectId/tasks/deadline/before", taskCtl.GetTasksByDeadlineBefore)
		auth.GET("/projects/:projectId/tasks/deadline/after", taskCtl.GetTasksByDeadlineAfter)
		auth.GET("/projects/:projectId/tasks/period", taskCtl.GetTasksByTimePeriod)

		// 任务预算
		auth.POST("/tasks/:id/budgets", taskCtl.AddTaskBudget)
		auth.PUT("/tasks/budgets/:budgetId", taskCtl.UpdateTaskBudget)
		auth.DELETE("/tasks/budgets/:budgetId", taskCtl.DeleteTaskBudget)
		auth.GET("/tasks/:id/budgets", taskCtl.GetTaskBudgets)
	}

	return r
}
