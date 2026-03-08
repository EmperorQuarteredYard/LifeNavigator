package router

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(
	userCtl *controller.UserController,
	inviteCtl *controller.InviteController,
	projectCtl *controller.ProjectController,
	taskCtl *controller.TaskController,
	accountCtl *controller.AccountController,
	aiFeatureCtl *controller.AIFeatureController,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Device-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

		// 任务
		auth.GET("/v0/tasks/:id", taskCtl.GetTask)
		auth.GET("/v0/projects/:id/tasks", taskCtl.ListTasks) // 根据 project_id 查询（从查询参数获取）
		auth.GET("/v0/tasks", taskCtl.ListTasks)              // 当前用户所有任务
		auth.PUT("/v0/tasks/:id", taskCtl.UpdateTask)
		auth.POST("/v0/tasks", taskCtl.CreateTask)
		auth.DELETE("/v0/tasks/:id", taskCtl.DeleteTask)

		// 任务预算
		auth.POST("/v0/tasks/:id/budgets", taskCtl.SetPayment)
		auth.PUT("/v0/tasks/budgets/:budgetId", taskCtl.UpdatePayment)
		auth.DELETE("/v0/tasks/budgets/:budgetId", taskCtl.DeletePayment)
		auth.GET("/v0/tasks/:id/budgets", taskCtl.GetPayments)

		// 账户
		auth.POST("/v0/accounts", accountCtl.CreateAccount)             // 创建账户
		auth.GET("/v0/accounts", accountCtl.ListAccounts)               // 获取当前用户所有账户
		auth.GET("/v0/accounts/:id", accountCtl.GetAccount)             // 获取单个账户详情
		auth.DELETE("/v0/accounts/:id", accountCtl.DeleteAccount)       // 删除账户
		auth.POST("/v0/accounts/:id/balance", accountCtl.AdjustBalance) // 调整余额

		// 账户关联查询
		auth.GET("/v0/accounts/:id/tasks", accountCtl.ListLinkedTasks)       // 获取账户关联的任务和付款
		auth.GET("/v0/accounts/:id/payments", accountCtl.ListLinkedPayments) // 获取账户关联的付款

		// AI 功能
		auth.POST("/v0/ai/reduce-project", aiFeatureCtl.ReduceProject) // AI 辅助创建项目
		auth.POST("/v0/ai/summary", aiFeatureCtl.Summary)              // AI 总结用户成就
	}

	return r
}
