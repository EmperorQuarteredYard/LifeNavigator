const DashboardPage = {
    render() {
        const page = document.getElementById('page-dashboard');
        if (!page) return;
        
        const user = Store.getUser();
        const projects = Store.state.projects || [];
        const tasks = Store.state.tasks || [];
        const accounts = Store.state.accounts || [];
        
        const activeTasks = tasks.filter(t => t.status === 0 || t.status === 1);
        const completedTasks = tasks.filter(t => t.status === 2);
        const overdueTasks = tasks.filter(t => t.deadline && Utils.isOverdue(t.deadline) && t.status !== 2);
        
        const totalBalance = accounts.reduce((sum, a) => sum + (a.balance || 0), 0);
        
        page.innerHTML = `
            <div class="page-header">
                <h1>欢迎回来，${Utils.escapeHtml(user?.nickname || user?.username)}</h1>
                <button class="btn btn-outline" onclick="App.refreshData()">刷新数据</button>
            </div>
            
            <div class="grid grid-4 gap-3 mt-3">
                <div class="stat-card">
                    <div class="stat-value">${projects.length}</div>
                    <div class="stat-label">项目总数</div>
                </div>
                <div class="stat-card success">
                    <div class="stat-value">${activeTasks.length}</div>
                    <div class="stat-label">进行中任务</div>
                </div>
                <div class="stat-card warning">
                    <div class="stat-value">${completedTasks.length}</div>
                    <div class="stat-label">已完成任务</div>
                </div>
                <div class="stat-card danger">
                    <div class="stat-value">${overdueTasks.length}</div>
                    <div class="stat-label">逾期任务</div>
                </div>
            </div>
            
            <div class="grid grid-2 gap-3 mt-4">
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">最近项目</h3>
                        <a href="#projects" class="btn btn-sm btn-outline">查看全部</a>
                    </div>
                    <div class="card-body">
                        ${this.renderRecentProjects(projects.slice(0, 5))}
                    </div>
                </div>
                
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">待办任务</h3>
                        <a href="#tasks" class="btn btn-sm btn-outline">查看全部</a>
                    </div>
                    <div class="card-body">
                        ${this.renderPendingTasks(activeTasks.slice(0, 5))}
                    </div>
                </div>
            </div>
            
            <div class="grid grid-2 gap-3 mt-3">
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">账户概览</h3>
                        <a href="#accounts" class="btn btn-sm btn-outline">管理账户</a>
                    </div>
                    <div class="card-body">
                        ${this.renderAccountSummary(accounts, totalBalance)}
                    </div>
                </div>
                
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">快捷操作</h3>
                    </div>
                    <div class="card-body">
                        <div class="grid grid-2 gap-2">
                            <a href="#projects/create" class="btn btn-primary btn-block">新建项目</a>
                            <a href="#tasks/create" class="btn btn-success btn-block">新建任务</a>
                            <a href="#accounts/create" class="btn btn-secondary btn-block">新建账户</a>
                            <a href="#ai" class="btn btn-outline btn-block">AI助手</a>
                        </div>
                    </div>
                </div>
            </div>
        `;
    },
    
    renderRecentProjects(projects) {
        if (projects.length === 0) {
            return `
                <div class="empty-state">
                    <p class="text-muted">暂无项目</p>
                    <a href="#projects/create" class="btn btn-primary btn-sm">创建项目</a>
                </div>
            `;
        }
        
        return `
            <table class="table">
                <thead>
                    <tr>
                        <th>名称</th>
                        <th>任务数</th>
                        <th>创建时间</th>
                    </tr>
                </thead>
                <tbody>
                    ${projects.map(p => `
                        <tr>
                            <td><a href="#projects/detail/${p.id}">${Utils.escapeHtml(p.name)}</a></td>
                            <td>${p.max_task_id || 0}</td>
                            <td>${Utils.formatDateShort(p.created_at)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    },
    
    renderPendingTasks(tasks) {
        if (tasks.length === 0) {
            return `
                <div class="empty-state">
                    <p class="text-muted">暂无待办任务</p>
                </div>
            `;
        }
        
        return `
            <div class="task-list">
                ${tasks.map(t => `
                    <div class="task-item flex-between ${t.deadline && Utils.isOverdue(t.deadline) ? 'overdue' : ''}">
                        <div>
                            <div class="task-name">${Utils.escapeHtml(t.name)}</div>
                            <div class="task-meta text-muted">
                                ${t.deadline ? `截止: ${Utils.formatDateShort(t.deadline)}` : ''}
                            </div>
                        </div>
                        <span class="badge ${Utils.getStatusClass(t.status, Config.TASK_STATUS)}">
                            ${Utils.getStatusLabel(t.status, Config.TASK_STATUS)}
                        </span>
                    </div>
                `).join('')}
            </div>
        `;
    },
    
    renderAccountSummary(accounts, totalBalance) {
        if (accounts.length === 0) {
            return `
                <div class="empty-state">
                    <p class="text-muted">暂无账户</p>
                    <a href="#accounts/create" class="btn btn-primary btn-sm">创建账户</a>
                </div>
            `;
        }
        
        return `
            <div class="mb-2">
                <div class="text-muted mb-1">总资产</div>
                <div class="stat-value">${Utils.formatCurrency(totalBalance)}</div>
            </div>
            <div class="account-list">
                ${accounts.slice(0, 4).map(a => `
                    <div class="account-item flex-between">
                        <span>${Utils.escapeHtml(a.type)}</span>
                        <span class="${a.balance >= 0 ? 'text-success' : 'text-danger'}">
                            ${Utils.formatCurrency(a.balance)}
                        </span>
                    </div>
                `).join('')}
            </div>
        `;
    }
};
