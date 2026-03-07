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
		auth.GET("/v0/projects", projectCtl.GetProjectsByUser) // 获取当前用户项目列表
		auth.PUT("/v0/projects/:id", projectCtl.UpdateProject)
		auth.DELETE("/v0/projects/:id", projectCtl.DeleteProject)

		// 项目预算
		auth.POST("/v0/projects/:id/budgets", projectCtl.AddBudget)
		auth.PUT("/v0/projects/budgets/:budgetId", projectCtl.UpdateBudget)
		auth.DELETE("/v0/projects/budgets/:budgetId", projectCtl.DeleteBudget)
		auth.GET("/v0/projects/:id/budgets/summary", projectCtl.GetBudgetSummary)

		// 任务
		auth.GET("/v0/tasks/:id", taskCtl.GetTask)
		auth.GET("/v0/projects/:id/tasks", taskCtl.ListTasks) // 根据 project_id 查询（从查询参数获取）
		auth.GET("/v0/tasks", taskCtl.ListTasks)              // 当前用户所有任务
		auth.PUT("/v0/tasks/:id", taskCtl.UpdateTask)
		auth.PUT("/v0/projects/:id/tasks/:id", taskCtl.UpdateTask) // 注意路径参数重复，需后端处理
		auth.POST("/v0/tasks", taskCtl.CreateTask)
		auth.DELETE("/v0/tasks/:id", taskCtl.DeleteTask)

		// 任务预算
		auth.POST("/v0/tasks/:id/budgets", taskCtl.SetPayment)
		auth.PUT("/v0/tasks/budgets/:budgetId", taskCtl.UpdatePayment)
		auth.DELETE("/v0/tasks/budgets/:budgetId", taskCtl.DeletePayment)
		auth.GET("/v0/tasks/:id/budgets", taskCtl.GetPayments)
	}

	return r
}
