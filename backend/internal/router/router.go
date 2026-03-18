package router

import (
	"LifeNavigator/internal/controller"
	"LifeNavigator/pkg/jwt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化路由引擎，注入所有控制器依赖
func SetupRouter(
	accountCtrl *controller.AccountController,
	aiFeatureCtrl *controller.AIFeatureController,
	inviteCtrl *controller.InviteController,
	kanbanCtrl *controller.KanbanController,
	projectCtrl *controller.ProjectController,
	taskCtrl *controller.TaskController,
	userCtrl *controller.UserController,
) *gin.Engine {
	r := gin.Default()
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowOrigins []string
	if allowedOriginsEnv == "" {
		// 开发环境默认允许所有（可注释掉以强制生产环境必须配置）
		allowOrigins = []string{"*"}
	} else {
		// 分割逗号并去除空格
		for _, origin := range strings.Split(allowedOriginsEnv, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				allowOrigins = append(allowOrigins, trimmed)
			}
		}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins, // TODO 上线的时候应当指定该项为一个确定的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "device-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	v0 := r.Group("/api/v0")
	{
		// 公开路由（无需认证）
		public := v0.Group("/auth")
		{
			public.POST("/register", userCtrl.Register)
			public.POST("/login", userCtrl.Login)
			public.POST("/refresh", userCtrl.Refresh)
		}

		// 需要认证的路由
		authorized := v0.Group("")
		authorized.Use(jwt.JWTAuthMiddleware())
		{
			// 用户相关
			authorized.GET("/users/me", userCtrl.Profile)
			authorized.GET("/users/:id", userCtrl.GetUser)

			// 邀请码相关
			authorized.POST("/invite-codes", inviteCtrl.CreateInviteCode)
			authorized.GET("/invite-codes/:token", inviteCtrl.GetInviteCode)
			authorized.GET("/users/:id/invite-codes", inviteCtrl.ListUserInviteCodes)

			// 账户相关
			authorized.POST("/accounts", accountCtrl.CreateAccount)
			authorized.GET("/accounts", accountCtrl.ListAccounts)
			authorized.GET("/accounts/:id", accountCtrl.GetAccount)
			authorized.DELETE("/accounts/:id", accountCtrl.DeleteAccount)
			authorized.POST("/accounts/:id/balance", accountCtrl.AdjustBalance) // 调整余额
			authorized.GET("/accounts/:id/tasks", accountCtrl.ListLinkedTasks)  // 关联任务

			// 项目相关
			authorized.POST("/projects", projectCtrl.CreateProject)
			authorized.GET("/projects", projectCtrl.GetProjectsByUser) // 当前用户的项目列表
			authorized.GET("/projects/:id", projectCtrl.GetProject)
			authorized.PUT("/projects/:id", projectCtrl.UpdateProject)
			authorized.DELETE("/projects/:id", projectCtrl.DeleteProject)

			// 项目预算
			authorized.POST("/projects/:id/budgets", projectCtrl.AddBudget)
			authorized.PUT("/projects/budgets/:budgetId", projectCtrl.UpdateBudget)
			authorized.DELETE("/projects/budgets/:budgetId", projectCtrl.DeleteBudget)

			// 任务相关
			authorized.POST("/tasks", taskCtrl.CreateTask)
			authorized.GET("/tasks", taskCtrl.ListTasks) // 支持 project_id 查询
			authorized.GET("/tasks/:id", taskCtrl.GetTask)
			authorized.PUT("/tasks/:id", taskCtrl.UpdateTask)
			authorized.DELETE("/tasks/:id", taskCtrl.DeleteTask)
			authorized.POST("/tasks/:id/finish", taskCtrl.FinishTask)
			authorized.PUT("/tasks/:id/status", taskCtrl.UpdateTaskStatus) // 更新状态

			// 任务依赖
			authorized.GET("/tasks/:id/prerequisites", taskCtrl.GetPrerequisites)
			authorized.GET("/tasks/:id/postrequisites", taskCtrl.GetPostrequisites)
			authorized.POST("/tasks/:id/prerequisites", taskCtrl.SetPrerequisites)
			authorized.DELETE("/tasks/:id/prerequisites", taskCtrl.UnsetPrerequisites) // 需在body中传prerequisite_id

			// 任务付款
			authorized.POST("/tasks/:id/payments", taskCtrl.SetPayment)
			authorized.GET("/tasks/:id/payments", taskCtrl.GetPayments)
			authorized.PUT("/tasks/payments/:id", taskCtrl.UpdatePayment)
			authorized.DELETE("/tasks/payments/:id", taskCtrl.DeletePayment)

			// 看板相关
			authorized.POST("/kanbans", kanbanCtrl.CreateKanban)
			authorized.GET("/kanbans", kanbanCtrl.ListKanbans)
			authorized.GET("/kanbans/:id", kanbanCtrl.GetKanban)
			authorized.PUT("/kanbans/:id", kanbanCtrl.UpdateKanban)
			authorized.DELETE("/kanbans/:id", kanbanCtrl.DeleteKanban)
			authorized.GET("/kanbans/:id/tasks", kanbanCtrl.GetKanbanTasks)            // 看板任务列表（带状态过滤）
			authorized.GET("/kanbans/default/tasks", kanbanCtrl.GetDefaultKanbanTasks) // 默认看板任务
			authorized.POST("/kanbans/:id/default", kanbanCtrl.SetDefaultKanban)       // 设为默认看板

			// AI 功能（SSE 流式响应）
			authorized.POST("/ai/reduce-project", aiFeatureCtrl.ReduceProject)
			authorized.POST("/ai/summary", aiFeatureCtrl.Summary)
		}
	}

	return r
}
