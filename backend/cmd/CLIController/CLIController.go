// cli_controller.go
package clicontroller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CLIController 命令行控制器
type CLIController struct {
	userService       service.UserService
	accountService    service.AccountService
	projectService    service.ProjectService
	taskService       service.TaskService
	inviteCodeService service.InviteCodeService
	inviteUserService service.InviteUserService

	reader          *bufio.Reader
	currentUserID   uint64
	currentUsername string
}

// NewCLIController 创建控制器实例，需传入所有服务
func NewCLIController(
	userService service.UserService,
	accountService service.AccountService,
	projectService service.ProjectService,
	taskService service.TaskService,
	inviteCodeService service.InviteCodeService,
	inviteUserService service.InviteUserService,
) *CLIController {
	return &CLIController{
		userService:       userService,
		accountService:    accountService,
		projectService:    projectService,
		taskService:       taskService,
		inviteCodeService: inviteCodeService,
		inviteUserService: inviteUserService,
		reader:            bufio.NewReader(os.Stdin),
	}
}

// Run 启动交互式命令行
func (cli *CLIController) Run() {
	fmt.Println("CLI started. Unit 'help' for commands.")
	for {
		fmt.Print("> ")
		input, err := cli.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		parts := strings.Fields(input)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "exit", "quit":
			fmt.Println("Bye!")
			return
		case "help":
			cli.printHelp()
		case "login":
			cli.handleLogin(args)
		case "register":
			cli.handleRegister(args)
		case "get-user":
			cli.handleGetUser(args)
		case "delete-user":
			cli.handleDeleteUser(args)
		case "list-accounts":
			cli.handleListAccounts(args)
		case "get-account":
			cli.handleGetAccount(args)
		case "delete-account":
			cli.handleDeleteAccount(args)
		case "list-projects":
			cli.handleListProjects(args)
		case "get-project":
			cli.handleGetProject(args)
		case "delete-project":
			cli.handleDeleteProject(args)
		case "list-tasks":
			cli.handleListTasks(args)
		case "get-task":
			cli.handleGetTask(args)
		case "delete-task":
			cli.handleDeleteTask(args)
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}

// printHelp 打印帮助信息
func (cli *CLIController) printHelp() {
	helpText := `
Commands:
  login <username> <password>                 - 登录并设置当前用户
  register                                     - 注册新用户（交互式）
  get-user [--id <id> | --username <username>] - 查看用户信息（默认当前用户）
  delete-user [--id <id>]                      - 删除用户（默认当前用户）

  list-accounts                                 - 列出当前用户的所有账户
  get-account <account-id>                      - 查看指定账户详情
  delete-account <account-id>                    - 删除指定账户

  list-projects                                 - 列出当前用户的所有项目
  get-project <project-id>                       - 查看指定项目详情
  delete-project <project-id>                     - 删除指定项目

  list-tasks [--project <project-id>]            - 列出任务（默认所有任务，可指定项目）
  get-task <task-id>                              - 查看指定任务详情
  delete-task <task-id>                            - 删除指定任务

  help                                           - 显示本帮助
  exit / quit                                    - 退出CLI
`
	fmt.Println(helpText)
}

// parseFlags 解析命令行参数中的标志（如 --key value），返回标志映射和剩余位置参数
func parseFlags(args []string) (map[string]string, []string) {
	flags := make(map[string]string)
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = "" // 布尔标志，无值
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) == 2 {
			key := arg[1:]
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = ""
			}
		} else {
			positional = append(positional, arg)
		}
	}
	return flags, positional
}

// requireLogin 检查是否已登录，若未登录则提示并返回false
func (cli *CLIController) requireLogin() bool {
	if cli.currentUserID == 0 {
		fmt.Println("You must login first. Use 'login' command.")
		return false
	}
	return true
}

// handleLogin 处理登录
func (cli *CLIController) handleLogin(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: login <username> <password>")
		return
	}
	username := args[0]
	password := args[1]
	user, err := cli.userService.Login(username, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	cli.currentUserID = user.ID
	cli.currentUsername = user.Username
	fmt.Printf("Login successful. Welcome %s (ID: %d)\n", user.Username, user.ID)
}

// handleRegister 交互式注册
func (cli *CLIController) handleRegister(args []string) {
	fmt.Println("Enter user information:")
	fmt.Print("Username: ")
	username, _ := cli.reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := cli.reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Nickname: ")
	nickname, _ := cli.reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	fmt.Print("Email: ")
	email, _ := cli.reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Phone: ")
	phone, _ := cli.reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	user := &models.User{
		Username: username,
		Password: password,
		Nickname: nickname,
		Email:    email,
		Phone:    phone,
		Role:     "user", // 默认角色
	}
	err := cli.userService.Register(user)
	if err != nil {
		fmt.Printf("Registration failed: %v\n", err)
		return
	}
	fmt.Printf("User %s registered successfully with ID %d\n", username, user.ID)
}

// handleGetUser 查看用户信息
func (cli *CLIController) handleGetUser(args []string) {
	if !cli.requireLogin() {
		return
	}
	flags, pos := parseFlags(args)
	var targetID uint64
	if idStr, ok := flags["id"]; ok {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid user ID")
			return
		}
		targetID = id
	} else if username, ok := flags["username"]; ok {
		// 通过用户名查询，需要调用 GetByUsername，但该服务需要当前用户ID，我们使用当前用户ID来获取（只能查自己）
		user, err := cli.userService.GetByUsername(username, cli.currentUserID)
		if err != nil {
			fmt.Printf("Failed to get user: %v\n", err)
			return
		}
		targetID = user.ID
	} else if len(pos) > 0 {
		// 尝试将第一个位置参数作为ID
		id, err := strconv.ParseUint(pos[0], 10, 64)
		if err != nil {
			fmt.Println("Invalid user ID")
			return
		}
		targetID = id
	} else {
		// 默认当前用户
		targetID = cli.currentUserID
	}

	user, err := cli.userService.GetByID(targetID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		return
	}
	// 打印用户信息（隐藏密码）
	fmt.Printf("ID: %d\nUsername: %s\nNickname: %s\nEmail: %s\nPhone: %s\nRole: %s\nCreatedAt: %s\n",
		user.ID, user.Username, user.Nickname, user.Email, user.Phone, user.Role, user.CreatedAt.Format(time.RFC3339))
}

// handleDeleteUser 删除用户
func (cli *CLIController) handleDeleteUser(args []string) {
	if !cli.requireLogin() {
		return
	}
	flags, pos := parseFlags(args)
	var targetID uint64
	if idStr, ok := flags["id"]; ok {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid user ID")
			return
		}
		targetID = id
	} else if len(pos) > 0 {
		id, err := strconv.ParseUint(pos[0], 10, 64)
		if err != nil {
			fmt.Println("Invalid user ID")
			return
		}
		targetID = id
	} else {
		targetID = cli.currentUserID
	}
	// 询问确认
	fmt.Printf("Are you sure you want to delete user %d? (yes/no): ", targetID)
	confirm, _ := cli.reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	err := cli.userService.HardDeleteByID(targetID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to delete user: %v\n", err)
		return
	}
	fmt.Printf("User %d deleted successfully.\n", targetID)
	if targetID == cli.currentUserID {
		// 如果删除了当前用户，登出
		cli.currentUserID = 0
		cli.currentUsername = ""
		fmt.Println("You have been logged out.")
	}
}

// handleListAccounts 列出当前用户的所有账户
func (cli *CLIController) handleListAccounts(args []string) {
	if !cli.requireLogin() {
		return
	}
	accounts, err := cli.accountService.ListByUserID(cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to list accounts: %v\n", err)
		return
	}
	if len(accounts) == 0 {
		fmt.Println("No accounts found.")
		return
	}
	fmt.Println("Accounts:")
	for _, acc := range accounts {
		fmt.Printf("  ID: %d, Unit: %s, Balance: %.2f\n", acc.ID, acc.Type, acc.Balance)
	}
}

// handleGetAccount 查看指定账户
func (cli *CLIController) handleGetAccount(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: get-account <account-id>")
		return
	}
	accountID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid account ID")
		return
	}
	account, err := cli.accountService.GetByAccountID(cli.currentUserID, accountID)
	if err != nil {
		fmt.Printf("Failed to get account: %v\n", err)
		return
	}
	fmt.Printf("Account ID: %d\nUnit: %s\nBalance: %.2f\n", account.ID, account.Type, account.Balance)
}

// handleDeleteAccount 删除指定账户
func (cli *CLIController) handleDeleteAccount(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: delete-account <account-id>")
		return
	}
	accountID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid account ID")
		return
	}
	// 先获取账户对象
	account, err := cli.accountService.GetByAccountID(cli.currentUserID, accountID)
	if err != nil {
		fmt.Printf("Failed to get account: %v\n", err)
		return
	}
	// 确认
	fmt.Printf("Are you sure you want to delete account %d (Unit: %s)? (yes/no): ", account.ID, account.Type)
	confirm, _ := cli.reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	err = cli.accountService.DeleteAccount(account)
	if err != nil {
		fmt.Printf("Failed to delete account: %v\n", err)
		return
	}
	fmt.Printf("Account %d deleted successfully.\n", accountID)
}

// handleListProjects 列出当前用户的所有项目
func (cli *CLIController) handleListProjects(args []string) {
	if !cli.requireLogin() {
		return
	}
	// 分页参数可扩展，这里简单列出前100条
	projects, err := cli.projectService.ListByUserID(cli.currentUserID, 0, 100)
	if err != nil {
		fmt.Printf("Failed to list projects: %v\n", err)
		return
	}
	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}
	fmt.Println("Projects:")
	for _, p := range projects {
		fmt.Printf("  ID: %d, Name: %s, Description: %s\n", p.ID, p.Name, p.Description)
	}
}

// handleGetProject 查看指定项目
func (cli *CLIController) handleGetProject(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: get-project <project-id>")
		return
	}
	projectID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid project ID")
		return
	}
	project, err := cli.projectService.GetByID(projectID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to get project: %v\n", err)
		return
	}
	fmt.Printf("Project ID: %d\nName: %s\nDescription: %s\n", project.ID, project.Name, project.Description)
}

// handleDeleteProject 删除指定项目
func (cli *CLIController) handleDeleteProject(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: delete-project <project-id>")
		return
	}
	projectID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid project ID")
		return
	}
	// 确认
	fmt.Printf("Are you sure you want to delete project %d? (yes/no): ", projectID)
	confirm, _ := cli.reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	err = cli.projectService.Delete(projectID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to delete project: %v\n", err)
		return
	}
	fmt.Printf("Project %d deleted successfully.\n", projectID)
}

// handleListTasks 列出任务
func (cli *CLIController) handleListTasks(args []string) {
	if !cli.requireLogin() {
		return
	}
	flags, _ := parseFlags(args)
	var projectID uint64
	if pidStr, ok := flags["project"]; ok {
		pid, err := strconv.ParseUint(pidStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid project ID")
			return
		}
		projectID = pid
	}
	var tasks []models.Task
	var total int64
	var err error
	if projectID != 0 {
		tasks, total, err = cli.taskService.ListByProjectID(projectID, 1, 100, cli.currentUserID)
	} else {
		tasks, total, err = cli.taskService.ListByUserID(cli.currentUserID, 0, 100)
	}
	if err != nil {
		fmt.Printf("Failed to list tasks: %v\n", err)
		return
	}
	fmt.Printf("Total tasks: %d\n", total)
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Println("Tasks:")
	for _, t := range tasks {
		deadlineStr := "none"
		if t.Deadline != nil {
			deadlineStr = t.Deadline.Format(time.RFC3339)
		}
		fmt.Printf("  ID: %d, Name: %s, Status: %d, Deadline: %s\n", t.ID, t.Name, t.Status, deadlineStr)
	}
}

// handleGetTask 查看指定任务
func (cli *CLIController) handleGetTask(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: get-task <task-id>")
		return
	}
	taskID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}
	task, err := cli.taskService.GetByID(taskID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to get task: %v\n", err)
		return
	}
	deadlineStr := "none"
	if task.Deadline != nil {
		deadlineStr = task.Deadline.Format(time.RFC3339)
	}
	fmt.Printf("Task ID: %d\nName: %s\nDescription: %s\nStatus: %d\nDeadline: %s\n",
		task.ID, task.Name, task.Description, task.Status, deadlineStr)
}

// handleDeleteTask 删除指定任务
func (cli *CLIController) handleDeleteTask(args []string) {
	if !cli.requireLogin() {
		return
	}
	if len(args) < 1 {
		fmt.Println("Usage: delete-task <task-id>")
		return
	}
	taskID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		fmt.Println("Invalid task ID")
		return
	}
	// 确认
	fmt.Printf("Are you sure you want to delete task %d? (yes/no): ", taskID)
	confirm, _ := cli.reader.ReadString('\n')
	confirm = strings.TrimSpace(confirm)
	if confirm != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}
	err = cli.taskService.Delete(taskID, cli.currentUserID)
	if err != nil {
		fmt.Printf("Failed to delete task: %v\n", err)
		return
	}
	fmt.Printf("Task %d deleted successfully.\n", taskID)
}
