const ProjectsPage = {
    render() {
        const page = document.getElementById('page-projects');
        if (!page) return;
        
        const projects = Store.state.projects || [];
        
        page.innerHTML = `
            <div class="page-header flex-between">
                <h1>项目管理</h1>
                <a href="#projects/create" class="btn btn-primary">新建项目</a>
            </div>
            
            ${projects.length === 0 ? this.renderEmptyState() : this.renderProjectList(projects)}
        `;
    },
    
    renderEmptyState() {
        return `
            <div class="card">
                <div class="empty-state">
                    <div class="empty-state-icon">📁</div>
                    <h3 class="empty-state-title">暂无项目</h3>
                    <p class="empty-state-description">创建你的第一个项目，开始管理任务和预算</p>
                    <a href="#projects/create" class="btn btn-primary">创建项目</a>
                </div>
            </div>
        `;
    },
    
    renderProjectList(projects) {
        return `
            <div class="grid grid-3 gap-3">
                ${projects.map(p => this.renderProjectCard(p)).join('')}
            </div>
        `;
    },
    
    renderProjectCard(project) {
        const budgets = project.budgets || [];
        const totalBudget = budgets.reduce((sum, b) => sum + (b.budget || 0), 0);
        const totalUsed = budgets.reduce((sum, b) => sum + (b.used || 0), 0);
        const progress = Utils.calculateProgress(totalUsed, totalBudget);
        
        return `
            <div class="card project-card">
                <div class="card-header">
                    <h3 class="card-title">
                        <a href="#projects/detail/${project.id}">${Utils.escapeHtml(project.name)}</a>
                    </h3>
                    <div class="dropdown">
                        <button class="btn btn-sm btn-outline" onclick="ProjectsPage.toggleDropdown(this)">⋮</button>
                        <div class="dropdown-menu">
                            <a class="dropdown-item" href="#projects/edit/${project.id}">编辑</a>
                            <a class="dropdown-item" href="#tasks/create?project_id=${project.id}">添加任务</a>
                            <div class="dropdown-divider"></div>
                            <span class="dropdown-item text-danger" onclick="ProjectsPage.deleteProject(${project.id})">删除</span>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <p class="text-muted mb-2">${Utils.escapeHtml(project.description || '暂无描述')}</p>
                    <div class="project-stats flex-between mb-2">
                        <span class="text-muted">任务数: ${project.max_task_id || 0}</span>
                        <span class="text-muted">${Utils.formatDateShort(project.created_at)}</span>
                    </div>
                    ${totalBudget > 0 ? `
                        <div class="budget-progress">
                            <div class="flex-between mb-1">
                                <span class="text-muted">预算使用</span>
                                <span>${Utils.formatCurrency(totalUsed)} / ${Utils.formatCurrency(totalBudget)}</span>
                            </div>
                            <div class="progress-bar">
                                <div class="progress-fill" style="width: ${progress}%"></div>
                            </div>
                        </div>
                    ` : ''}
                </div>
            </div>
        `;
    },
    
    toggleDropdown(btn) {
        const dropdown = btn.parentElement;
        const menu = dropdown.querySelector('.dropdown-menu');
        
        document.querySelectorAll('.dropdown-menu.active').forEach(m => {
            if (m !== menu) m.classList.remove('active');
        });
        
        menu.classList.toggle('active');
    },
    
    async deleteProject(id) {
        if (!Utils.confirm('确定要删除这个项目吗？相关的任务和预算也会被删除。')) {
            return;
        }
        
        try {
            await ProjectAPI.delete(id);
            Store.removeProject(id);
            this.render();
            Utils.showNotification('项目已删除', 'success');
        } catch (error) {
            Utils.showNotification('删除失败: ' + error.message, 'danger');
        }
    }
};

const ProjectFormPage = {
    isEdit: false,
    projectId: null,
    
    render(params = {}) {
        const page = document.getElementById('page-project-form');
        if (!page) return;
        
        this.isEdit = params.id ? true : false;
        this.projectId = params.id || null;
        
        page.innerHTML = `
            <div class="page-header">
                <h1>${this.isEdit ? '编辑项目' : '新建项目'}</h1>
            </div>
            
            <div class="card">
                <form id="project-form">
                    <div class="form-group">
                        <label class="form-label" for="name">项目名称 *</label>
                        <input type="text" id="name" name="name" class="form-input" placeholder="请输入项目名称" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="description">项目描述</label>
                        <textarea id="description" name="description" class="form-textarea" placeholder="请输入项目描述"></textarea>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="refreshInterval">刷新间隔（秒）</label>
                        <input type="number" id="refreshInterval" name="refreshInterval" class="form-input" placeholder="0表示不自动刷新" min="0" value="0">
                    </div>
                    <div id="form-error" class="form-error hidden"></div>
                    <div class="flex gap-2">
                        <button type="submit" class="btn btn-primary">${this.isEdit ? '保存' : '创建'}</button>
                        <a href="#projects" class="btn btn-outline">取消</a>
                    </div>
                </form>
            </div>
        `;
        
        this.bindEvents();
        
        if (this.isEdit) {
            this.loadProject();
        }
    },
    
    bindEvents() {
        const form = document.getElementById('project-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleSubmit(form);
            });
        }
    },
    
    async loadProject() {
        try {
            const project = await ProjectAPI.get(this.projectId);
            
            document.getElementById('name').value = project.name || '';
            document.getElementById('description').value = project.description || '';
            document.getElementById('refreshInterval').value = project.refresh_interval || 0;
            
        } catch (error) {
            Utils.showNotification('加载项目失败: ' + error.message, 'danger');
            Router.navigate('projects');
        }
    },
    
    async handleSubmit(form) {
        const errorEl = document.getElementById('form-error');
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const name = form.name.value.trim();
        const description = form.description.value.trim();
        const refreshInterval = parseInt(form.refreshInterval.value) || 0;
        
        if (!name) {
            errorEl.textContent = '请输入项目名称';
            errorEl.classList.remove('hidden');
            return;
        }
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '保存中...';
            errorEl.classList.add('hidden');
            
            const data = {
                name,
                description,
                refresh_interval: refreshInterval
            };
            
            if (this.isEdit) {
                const project = await ProjectAPI.update(this.projectId, data);
                Store.updateProject(this.projectId, project);
                Utils.showNotification('项目已更新', 'success');
            } else {
                const project = await ProjectAPI.create(data);
                Store.addProject(project);
                Utils.showNotification('项目创建成功', 'success');
            }
            
            Router.navigate('projects');
            
        } catch (error) {
            errorEl.textContent = error.message;
            errorEl.classList.remove('hidden');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = this.isEdit ? '保存' : '创建';
        }
    }
};

const ProjectDetailPage = {
    projectId: null,
    
    async render(params = {}) {
        const page = document.getElementById('page-project-detail');
        if (!page) return;
        
        this.projectId = params.id;
        
        if (!this.projectId) {
            Router.navigate('projects');
            return;
        }
        
        page.innerHTML = `
            <div class="loading-overlay">
                <div class="loading-spinner"></div>
            </div>
        `;
        
        try {
            const project = await ProjectAPI.get(this.projectId);
            page.innerHTML = this.renderContent(project);
            this.bindEvents();
        } catch (error) {
            page.innerHTML = `
                <div class="alert alert-danger">加载失败: ${error.message}</div>
                <a href="#projects" class="btn btn-outline">返回列表</a>
            `;
        }
    },
    
    renderContent(project) {
        const budgets = project.budgets || [];
        const totalBudget = budgets.reduce((sum, b) => sum + (b.budget || 0), 0);
        const totalUsed = budgets.reduce((sum, b) => sum + (b.used || 0), 0);
        
        return `
            <div class="page-header flex-between">
                <div>
                    <h1>${Utils.escapeHtml(project.name)}</h1>
                    <p class="text-muted">${Utils.escapeHtml(project.description || '暂无描述')}</p>
                </div>
                <div class="flex gap-2">
                    <a href="#projects/edit/${project.id}" class="btn btn-outline">编辑</a>
                    <a href="#tasks/create?project_id=${project.id}" class="btn btn-primary">添加任务</a>
                </div>
            </div>
            
            <div class="grid grid-2 gap-3">
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">项目信息</h3>
                    </div>
                    <div class="card-body">
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">创建时间</span>
                            <span>${Utils.formatDate(project.created_at)}</span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">更新时间</span>
                            <span>${Utils.formatDate(project.updated_at)}</span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">刷新间隔</span>
                            <span>${project.refresh_interval || 0} 秒</span>
                        </div>
                        <div class="info-row flex-between">
                            <span class="text-muted">任务数</span>
                            <span>${project.max_task_id || 0}</span>
                        </div>
                    </div>
                </div>
                
                <div class="card">
                    <div class="card-header flex-between">
                        <h3 class="card-title">预算管理</h3>
                        <button class="btn btn-sm btn-outline" onclick="ProjectDetailPage.showBudgetModal()">添加预算</button>
                    </div>
                    <div class="card-body">
                        ${budgets.length === 0 ? `
                            <p class="text-muted text-center">暂无预算</p>
                        ` : `
                            <div class="mb-2">
                                <div class="flex-between mb-1">
                                    <span>总预算</span>
                                    <strong>${Utils.formatCurrency(totalBudget)}</strong>
                                </div>
                                <div class="flex-between mb-1">
                                    <span>已使用</span>
                                    <strong>${Utils.formatCurrency(totalUsed)}</strong>
                                </div>
                                <div class="progress-bar mt-2">
                                    <div class="progress-fill" style="width: ${Utils.calculateProgress(totalUsed, totalBudget)}%"></div>
                                </div>
                            </div>
                            <table class="table mt-3">
                                <thead>
                                    <tr>
                                        <th>类型</th>
                                        <th>账户</th>
                                        <th>预算</th>
                                        <th>已用</th>
                                        <th>操作</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${budgets.map(b => `
                                        <tr>
                                            <td><span class="badge badge-primary">${b.type || 'money'}</span></td>
                                            <td>${b.account_id || '未关联'}</td>
                                            <td>${Utils.formatCurrency(b.budget)}</td>
                                            <td>${Utils.formatCurrency(b.used)}</td>
                                            <td>
                                                <button class="btn btn-sm btn-outline" onclick="ProjectDetailPage.editBudget(${b.id})">编辑</button>
                                                <button class="btn btn-sm btn-danger" onclick="ProjectDetailPage.deleteBudget(${b.id})">删除</button>
                                            </td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        `}
                    </div>
                </div>
            </div>
            
            <div class="modal-overlay" id="budget-modal">
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title" id="budget-modal-title">添加预算</h3>
                        <button class="modal-close" onclick="ProjectDetailPage.hideBudgetModal()">&times;</button>
                    </div>
                    <div class="modal-body">
                        <form id="budget-form">
                            <input type="hidden" id="budget-id" name="id">
                            <div class="form-group">
                                <label class="form-label" for="budget-account">关联账户</label>
                                <select id="budget-account" name="account_id" class="form-select">
                                    <option value="0">不关联账户</option>
                                    ${(Store.state.accounts || []).map(a => `
                                        <option value="${a.id}">${Utils.escapeHtml(a.type)} (${Utils.formatCurrency(a.balance)})</option>
                                    `).join('')}
                                </select>
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="budget-type">预算类型</label>
                                <select id="budget-type" name="type" class="form-select" required>
                                    <option value="money">金钱</option>
                                    <option value="time">时间</option>
                                    <option value="token">Token</option>
                                    <option value="energy">精力</option>
                                </select>
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="budget-amount">预算金额</label>
                                <input type="number" id="budget-amount" name="budget" class="form-input" placeholder="请输入预算金额" min="0" step="0.01" required>
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="budget-used">已用金额</label>
                                <input type="number" id="budget-used" name="used" class="form-input" placeholder="0" min="0" step="0.01" value="0">
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-outline" onclick="ProjectDetailPage.hideBudgetModal()">取消</button>
                        <button class="btn btn-primary" onclick="ProjectDetailPage.saveBudget()">保存</button>
                    </div>
                </div>
            </div>
        `;
    },
    
    bindEvents() {
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal-overlay')) {
                this.hideBudgetModal();
            }
        });
    },
    
    showBudgetModal(budgetId = null) {
        const modal = document.getElementById('budget-modal');
        const title = document.getElementById('budget-modal-title');
        const form = document.getElementById('budget-form');
        
        form.reset();
        document.getElementById('budget-id').value = '';
        
        if (budgetId) {
            title.textContent = '编辑预算';
            const budgets = Store.state.projects.find(p => p.id === this.projectId)?.budgets || [];
            const budget = budgets.find(b => b.id === budgetId);
            if (budget) {
                document.getElementById('budget-id').value = budget.id;
                document.getElementById('budget-account').value = budget.account_id || 0;
                document.getElementById('budget-type').value = budget.type || 'money';
                document.getElementById('budget-amount').value = budget.budget;
                document.getElementById('budget-used').value = budget.used;
            }
        } else {
            title.textContent = '添加预算';
            document.getElementById('budget-type').value = 'money';
        }
        
        modal.classList.add('active');
    },
    
    hideBudgetModal() {
        const modal = document.getElementById('budget-modal');
        modal.classList.remove('active');
    },
    
    async saveBudget() {
        const form = document.getElementById('budget-form');
        const budgetId = document.getElementById('budget-id').value;
        
        const data = {
            account_id: parseInt(form.account_id.value) || 0,
            type: form.type.value || 'money',
            budget: parseFloat(form.budget.value) || 0,
            used: parseFloat(form.used.value) || 0
        };
        
        try {
            if (budgetId) {
                await ProjectAPI.updateBudget(budgetId, data);
                Utils.showNotification('预算已更新', 'success');
            } else {
                await ProjectAPI.addBudget(this.projectId, data);
                Utils.showNotification('预算已添加', 'success');
            }
            
            this.hideBudgetModal();
            this.render({ id: this.projectId });
        } catch (error) {
            Utils.showNotification('保存失败: ' + error.message, 'danger');
        }
    },
    
    editBudget(budgetId) {
        this.showBudgetModal(budgetId);
    },
    
    async deleteBudget(budgetId) {
        if (!Utils.confirm('确定要删除这个预算吗？')) return;
        
        try {
            await ProjectAPI.deleteBudget(budgetId);
            Utils.showNotification('预算已删除', 'success');
            this.render({ id: this.projectId });
        } catch (error) {
            Utils.showNotification('删除失败: ' + error.message, 'danger');
        }
    }
};
