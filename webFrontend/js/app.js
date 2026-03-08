const App = {
    async init() {
        console.log('LifeNavigator Web Frontend Initializing...');
        
        Router.init({
            'login': { page: 'login', authRequired: false, onLoad: () => LoginPage.render() },
            'register': { page: 'register', authRequired: false, onLoad: () => RegisterPage.render() },
            'dashboard': { page: 'dashboard', authRequired: true, onLoad: () => DashboardPage.render() },
            'projects': { page: 'projects', authRequired: true, onLoad: () => ProjectsPage.render() },
            'projects/create': { page: 'project-form', authRequired: true, onLoad: (params) => ProjectFormPage.render(params) },
            'projects/edit': { page: 'project-form', authRequired: true, onLoad: (params) => ProjectFormPage.render(params) },
            'projects/detail': { page: 'project-detail', authRequired: true, onLoad: (params) => ProjectDetailPage.render(params) },
            'tasks': { page: 'tasks', authRequired: true, onLoad: () => TasksPage.render() },
            'tasks/create': { page: 'task-form', authRequired: true, onLoad: (params) => TaskFormPage.render(params) },
            'tasks/edit': { page: 'task-form', authRequired: true, onLoad: (params) => TaskFormPage.render(params) },
            'tasks/detail': { page: 'task-detail', authRequired: true, onLoad: (params) => TaskDetailPage.render(params) },
            'accounts': { page: 'accounts', authRequired: true, onLoad: () => AccountsPage.render() },
            'accounts/create': { page: 'account-form', authRequired: true, onLoad: (params) => AccountFormPage.render(params) },
            'ai': { page: 'ai', authRequired: true, onLoad: () => AIPage.render() },
            'invite-codes': { page: 'invite-codes', authRequired: true, onLoad: () => InviteCodesPage.render() },
            'not-found': { page: 'not-found', authRequired: false }
        });
        
        Router.setBeforeNavigate((route, params) => {
            const routeConfig = Router.routes[route];
            if (routeConfig && routeConfig.authRequired && !Store.isAuthenticated()) {
                Router.navigate('login');
                return false;
            }
            return true;
        });
        
        this.renderLayout();
        
        if (Store.isAuthenticated()) {
            await this.loadInitialData();
            const hash = window.location.hash.slice(1) || 'dashboard';
            Router.navigate(hash);
        } else {
            const hash = window.location.hash.slice(1) || 'login';
            Router.navigate(hash);
        }
        
        console.log('LifeNavigator Web Frontend Initialized');
    },
    
    renderLayout() {
        const app = document.getElementById('app');
        
        if (Store.isAuthenticated()) {
            app.innerHTML = this.getMainLayout();
            this.setupNavigation();
        } else {
            app.innerHTML = this.getAuthLayout();
        }
    },
    
    getAuthLayout() {
        return `
            <div id="page-login" class="page auth-page"></div>
            <div id="page-register" class="page auth-page"></div>
        `;
    },
    
    getMainLayout() {
        const user = Store.getUser();
        const initials = Utils.getInitials(user?.nickname || user?.username);
        
        return `
            <div class="main-layout">
                <aside class="sidebar">
                    <div class="sidebar-header">
                        <h1 class="logo">LifeNavigator</h1>
                    </div>
                    <nav class="sidebar-nav">
                        <a href="#dashboard" class="nav-link">
                            <span class="nav-icon">📊</span>
                            <span>仪表盘</span>
                        </a>
                        <a href="#projects" class="nav-link">
                            <span class="nav-icon">📁</span>
                            <span>项目</span>
                        </a>
                        <a href="#tasks" class="nav-link">
                            <span class="nav-icon">✅</span>
                            <span>任务</span>
                        </a>
                        <a href="#accounts" class="nav-link">
                            <span class="nav-icon">💰</span>
                            <span>账户</span>
                        </a>
                        <a href="#ai" class="nav-link">
                            <span class="nav-icon">🤖</span>
                            <span>AI助手</span>
                        </a>
                        <a href="#invite-codes" class="nav-link">
                            <span class="nav-icon">🎫</span>
                            <span>邀请码</span>
                        </a>
                    </nav>
                    <div class="sidebar-footer">
                        <div class="user-info">
                            <div class="avatar">${initials}</div>
                            <div class="user-details">
                                <div class="user-name">${Utils.escapeHtml(user?.nickname || user?.username)}</div>
                                <div class="user-role">${Utils.escapeHtml(user?.role || 'user')}</div>
                            </div>
                        </div>
                        <button class="btn btn-outline btn-sm" onclick="App.logout()">退出</button>
                    </div>
                </aside>
                <main class="main-content">
                    <div id="page-dashboard" class="page"></div>
                    <div id="page-projects" class="page"></div>
                    <div id="page-project-form" class="page"></div>
                    <div id="page-project-detail" class="page"></div>
                    <div id="page-tasks" class="page"></div>
                    <div id="page-task-form" class="page"></div>
                    <div id="page-task-detail" class="page"></div>
                    <div id="page-accounts" class="page"></div>
                    <div id="page-account-form" class="page"></div>
                    <div id="page-ai" class="page"></div>
                    <div id="page-invite-codes" class="page"></div>
                    <div id="page-not-found" class="page">
                        <div class="empty-state">
                            <h2>404</h2>
                            <p>页面不存在</p>
                            <a href="#dashboard" class="btn btn-primary">返回首页</a>
                        </div>
                    </div>
                </main>
            </div>
        `;
    },
    
    setupNavigation() {
        const navLinks = document.querySelectorAll('.nav-link');
        navLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                navLinks.forEach(l => l.classList.remove('active'));
                link.classList.add('active');
            });
        });
    },
    
    async loadInitialData() {
        try {
            Store.setLoading(true);
            
            const [projects, tasks, accounts] = await Promise.all([
                ProjectAPI.list(),
                TaskAPI.list(),
                AccountAPI.list()
            ]);
            
            Store.setProjects(projects.items || []);
            Store.setTasks(tasks.list || []);
            Store.setAccounts(accounts || []);
            
        } catch (error) {
            console.error('Failed to load initial data:', error);
            Utils.showNotification('加载数据失败: ' + error.message, 'danger');
        } finally {
            Store.setLoading(false);
        }
    },
    
    logout() {
        if (Utils.confirm('确定要退出登录吗？')) {
            Store.clearAuth();
            this.renderLayout();
            Router.navigate('login');
            Utils.showNotification('已退出登录', 'info');
        }
    },
    
    async refreshData() {
        await this.loadInitialData();
        Utils.showNotification('数据已刷新', 'success');
    }
};

document.addEventListener('DOMContentLoaded', () => {
    App.init();
});
