const TasksPage = {
    render() {
        const page = document.getElementById('page-tasks');
        if (!page) return;
        
        const tasks = Store.state.tasks || [];
        const projects = Store.state.projects || [];
        
        page.innerHTML = `
            <div class="page-header flex-between">
                <h1>任务管理</h1>
                <a href="#tasks/create" class="btn btn-primary">新建任务</a>
            </div>
            
            <div class="card mb-3">
                <div class="tabs">
                    <span class="tab active" data-status="all">全部</span>
                    <span class="tab" data-status="0">待办</span>
                    <span class="tab" data-status="1">进行中</span>
                    <span class="tab" data-status="2">已完成</span>
                </div>
            </div>
            
            ${tasks.length === 0 ? this.renderEmptyState() : this.renderTaskList(tasks, projects)}
        `;
        
        this.bindEvents();
    },
    
    renderEmptyState() {
        return `
            <div class="card">
                <div class="empty-state">
                    <div class="empty-state-icon">✅</div>
                    <h3 class="empty-state-title">暂无任务</h3>
                    <p class="empty-state-description">创建你的第一个任务，开始规划你的工作</p>
                    <a href="#tasks/create" class="btn btn-primary">创建任务</a>
                </div>
            </div>
        `;
    },
    
    renderTaskList(tasks, projects) {
        const getProjectName = (projectId) => {
            const project = projects.find(p => p.id === projectId);
            return project ? project.name : '未知项目';
        };
        
        return `
            <div class="card">
                <table class="table">
                    <thead>
                        <tr>
                            <th>任务名称</th>
                            <th>所属项目</th>
                            <th>状态</th>
                            <th>截止时间</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${tasks.map(t => `
                            <tr class="task-row" data-status="${t.status}">
                                <td>
                                    <a href="#tasks/detail/${t.id}">${Utils.escapeHtml(t.name)}</a>
                                    ${t.description ? `<br><small class="text-muted">${Utils.escapeHtml(t.description.substring(0, 50))}${t.description.length > 50 ? '...' : ''}</small>` : ''}
                                </td>
                                <td>
                                    <a href="#projects/detail/${t.project_id}">${Utils.escapeHtml(getProjectName(t.project_id))}</a>
                                </td>
                                <td>
                                    <span class="badge ${Utils.getStatusClass(t.status, Config.TASK_STATUS)}">
                                        ${Utils.getStatusLabel(t.status, Config.TASK_STATUS)}
                                    </span>
                                </td>
                                <td class="${t.deadline && Utils.isOverdue(t.deadline) && t.status !== 2 ? 'text-danger' : ''}">
                                    ${t.deadline ? Utils.formatDateShort(t.deadline) : '-'}
                                    ${t.deadline && Utils.isOverdue(t.deadline) && t.status !== 2 ? '<br><small>已逾期</small>' : ''}
                                </td>
                                <td>
                                    <div class="flex gap-1">
                                        <button class="btn btn-sm btn-outline" onclick="TasksPage.changeStatus(${t.id}, ${t.status})">
                                            ${t.status === 2 ? '重开' : '完成'}
                                        </button>
                                        <a href="#tasks/edit/${t.id}" class="btn btn-sm btn-outline">编辑</a>
                                        <button class="btn btn-sm btn-danger" onclick="TasksPage.deleteTask(${t.id})">删除</button>
                                    </div>
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            </div>
        `;
    },
    
    bindEvents() {
        const tabs = document.querySelectorAll('.tab');
        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                tabs.forEach(t => t.classList.remove('active'));
                tab.classList.add('active');
                
                const status = tab.dataset.status;
                const rows = document.querySelectorAll('.task-row');
                
                rows.forEach(row => {
                    if (status === 'all' || row.dataset.status === status) {
                        row.style.display = '';
                    } else {
                        row.style.display = 'none';
                    }
                });
            });
        });
    },
    
    async changeStatus(id, currentStatus) {
        const newStatus = currentStatus === 2 ? 0 : 2;
        
        try {
            const task = await TaskAPI.update(id, { status: newStatus });
            Store.updateTask(id, task);
            this.render();
            Utils.showNotification(newStatus === 2 ? '任务已完成' : '任务已重开', 'success');
        } catch (error) {
            Utils.showNotification('操作失败: ' + error.message, 'danger');
        }
    },
    
    async deleteTask(id) {
        if (!Utils.confirm('确定要删除这个任务吗？')) return;
        
        try {
            await TaskAPI.delete(id);
            Store.removeTask(id);
            this.render();
            Utils.showNotification('任务已删除', 'success');
        } catch (error) {
            Utils.showNotification('删除失败: ' + error.message, 'danger');
        }
    }
};

const TaskFormPage = {
    isEdit: false,
    taskId: null,
    defaultProjectId: null,
    
    render(params = {}) {
        const page = document.getElementById('page-task-form');
        if (!page) return;
        
        this.isEdit = params.id ? true : false;
        this.taskId = params.id || null;
        this.defaultProjectId = params.project_id || null;
        
        const projects = Store.state.projects || [];
        
        if (projects.length === 0 && !this.isEdit) {
            page.innerHTML = `
                <div class="alert alert-warning">请先创建项目再添加任务</div>
                <a href="#projects/create" class="btn btn-primary">创建项目</a>
            `;
            return;
        }
        
        page.innerHTML = `
            <div class="page-header">
                <h1>${this.isEdit ? '编辑任务' : '新建任务'}</h1>
            </div>
            
            <div class="card">
                <form id="task-form">
                    <div class="form-group">
                        <label class="form-label" for="project_id">所属项目 *</label>
                        <select id="project_id" name="project_id" class="form-select" required>
                            <option value="">请选择项目</option>
                            ${projects.map(p => `
                                <option value="${p.id}">${Utils.escapeHtml(p.name)}</option>
                            `).join('')}
                        </select>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="name">任务名称 *</label>
                        <input type="text" id="name" name="name" class="form-input" placeholder="请输入任务名称" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="description">任务描述</label>
                        <textarea id="description" name="description" class="form-textarea" placeholder="请输入任务描述"></textarea>
                    </div>
                    <div class="grid grid-2 gap-2">
                        <div class="form-group">
                            <label class="form-label" for="status">状态</label>
                            <select id="status" name="status" class="form-select">
                                <option value="0">待办</option>
                                <option value="1">进行中</option>
                                <option value="2">已完成</option>
                                <option value="3">已取消</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label class="form-label" for="category">分类</label>
                            <input type="text" id="category" name="category" class="form-input" placeholder="如：学习、工作、生活">
                        </div>
                    </div>
                    <div class="grid grid-2 gap-2">
                        <div class="form-group">
                            <label class="form-label" for="deadline">截止时间</label>
                            <input type="datetime-local" id="deadline" name="deadline" class="form-input">
                        </div>
                        <div class="form-group">
                            <label class="form-label" for="type">任务类型</label>
                            <input type="number" id="type" name="type" class="form-input" placeholder="0" min="0" value="0">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">
                            <input type="checkbox" id="auto_calculated" name="auto_calculated"> 自动计算
                        </label>
                    </div>
                    <div id="form-error" class="form-error hidden"></div>
                    <div class="flex gap-2">
                        <button type="submit" class="btn btn-primary">${this.isEdit ? '保存' : '创建'}</button>
                        <a href="#tasks" class="btn btn-outline">取消</a>
                    </div>
                </form>
            </div>
        `;
        
        this.bindEvents();
        
        if (this.isEdit) {
            this.loadTask();
        } else if (this.defaultProjectId) {
            document.getElementById('project_id').value = this.defaultProjectId;
        }
    },
    
    bindEvents() {
        const form = document.getElementById('task-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleSubmit(form);
            });
        }
    },
    
    async loadTask() {
        try {
            const task = await TaskAPI.get(this.taskId);
            
            document.getElementById('project_id').value = task.project_id || '';
            document.getElementById('name').value = task.name || '';
            document.getElementById('description').value = task.description || '';
            document.getElementById('status').value = task.status || 0;
            document.getElementById('category').value = task.category || '';
            document.getElementById('type').value = task.type || 0;
            document.getElementById('auto_calculated').checked = task.auto_calculated || false;
            
            if (task.deadline) {
                const deadline = new Date(task.deadline);
                document.getElementById('deadline').value = deadline.toISOString().slice(0, 16);
            }
            
        } catch (error) {
            Utils.showNotification('加载任务失败: ' + error.message, 'danger');
            Router.navigate('tasks');
        }
    },
    
    async handleSubmit(form) {
        const errorEl = document.getElementById('form-error');
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const projectId = parseInt(form.project_id.value);
        const name = form.name.value.trim();
        const description = form.description.value.trim();
        const status = parseInt(form.status.value);
        const category = form.category.value.trim();
        const type = parseInt(form.type.value) || 0;
        const autoCalculated = form.auto_calculated.checked;
        const deadlineValue = form.deadline.value;
        
        if (!projectId || !name) {
            errorEl.textContent = '请填写必填项';
            errorEl.classList.remove('hidden');
            return;
        }
        
        const data = {
            project_id: projectId,
            name,
            description,
            status,
            category,
            type,
            auto_calculated: autoCalculated,
            deadline: deadlineValue ? new Date(deadlineValue).toISOString() : null
        };
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '保存中...';
            errorEl.classList.add('hidden');
            
            if (this.isEdit) {
                const task = await TaskAPI.update(this.taskId, data);
                Store.updateTask(this.taskId, task);
                Utils.showNotification('任务已更新', 'success');
            } else {
                const task = await TaskAPI.create(data);
                Store.addTask(task);
                Utils.showNotification('任务创建成功', 'success');
            }
            
            Router.navigate('tasks');
            
        } catch (error) {
            errorEl.textContent = error.message;
            errorEl.classList.remove('hidden');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = this.isEdit ? '保存' : '创建';
        }
    }
};

const TaskDetailPage = {
    taskId: null,
    
    async render(params = {}) {
        const page = document.getElementById('page-task-detail');
        if (!page) return;
        
        this.taskId = params.id;
        
        if (!this.taskId) {
            Router.navigate('tasks');
            return;
        }
        
        page.innerHTML = `
            <div class="loading-overlay">
                <div class="loading-spinner"></div>
            </div>
        `;
        
        try {
            const task = await TaskAPI.get(this.taskId);
            const payments = task.payments || [];
            
            page.innerHTML = this.renderContent(task, payments);
            this.bindEvents();
        } catch (error) {
            page.innerHTML = `
                <div class="alert alert-danger">加载失败: ${error.message}</div>
                <a href="#tasks" class="btn btn-outline">返回列表</a>
            `;
        }
    },
    
    renderContent(task, payments) {
        const projects = Store.state.projects || [];
        const project = projects.find(p => p.id === task.project_id);
        
        return `
            <div class="page-header flex-between">
                <div>
                    <h1>${Utils.escapeHtml(task.name)}</h1>
                    <p class="text-muted">${Utils.escapeHtml(task.description || '暂无描述')}</p>
                </div>
                <div class="flex gap-2">
                    <a href="#tasks/edit/${task.id}" class="btn btn-outline">编辑</a>
                    <button class="btn ${task.status === 2 ? 'btn-secondary' : 'btn-success'}" onclick="TaskDetailPage.toggleComplete()">
                        ${task.status === 2 ? '重开任务' : '完成任务'}
                    </button>
                </div>
            </div>
            
            <div class="grid grid-2 gap-3">
                <div class="card">
                    <div class="card-header">
                        <h3 class="card-title">任务信息</h3>
                    </div>
                    <div class="card-body">
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">所属项目</span>
                            <a href="#projects/detail/${task.project_id}">${Utils.escapeHtml(project?.name || '未知项目')}</a>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">状态</span>
                            <span class="badge ${Utils.getStatusClass(task.status, Config.TASK_STATUS)}">
                                ${Utils.getStatusLabel(task.status, Config.TASK_STATUS)}
                            </span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">分类</span>
                            <span>${Utils.escapeHtml(task.category || '未分类')}</span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">截止时间</span>
                            <span class="${task.deadline && Utils.isOverdue(task.deadline) && task.status !== 2 ? 'text-danger' : ''}">
                                ${task.deadline ? Utils.formatDate(task.deadline) : '-'}
                            </span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">创建时间</span>
                            <span>${Utils.formatDate(task.created_at)}</span>
                        </div>
                        <div class="info-row flex-between mb-2">
                            <span class="text-muted">更新时间</span>
                            <span>${Utils.formatDate(task.updated_at)}</span>
                        </div>
                        ${task.completed_at ? `
                            <div class="info-row flex-between">
                                <span class="text-muted">完成时间</span>
                                <span class="text-success">${Utils.formatDate(task.completed_at)}</span>
                            </div>
                        ` : ''}
                    </div>
                </div>
                
                <div class="card">
                    <div class="card-header flex-between">
                        <h3 class="card-title">付款记录</h3>
                        <button class="btn btn-sm btn-outline" onclick="TaskDetailPage.showPaymentModal()">添加付款</button>
                    </div>
                    <div class="card-body">
                        ${payments.length === 0 ? `
                            <p class="text-muted text-center">暂无付款记录</p>
                        ` : `
                            <table class="table">
                                <thead>
                                    <tr>
                                        <th>预算ID</th>
                                        <th>金额</th>
                                        <th>操作</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${payments.map(p => `
                                        <tr>
                                            <td>${p.budget_id}</td>
                                            <td>${Utils.formatCurrency(p.amount)}</td>
                                            <td>
                                                <button class="btn btn-sm btn-danger" onclick="TaskDetailPage.deletePayment(${p.id})">删除</button>
                                            </td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        `}
                    </div>
                </div>
            </div>
            
            <div class="modal-overlay" id="payment-modal">
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">添加付款</h3>
                        <button class="modal-close" onclick="TaskDetailPage.hidePaymentModal()">&times;</button>
                    </div>
                    <div class="modal-body">
                        <form id="payment-form">
                            <div class="form-group">
                                <label class="form-label" for="payment-budget">预算ID *</label>
                                <input type="number" id="payment-budget" name="budget_id" class="form-input" placeholder="请输入预算ID" required>
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="payment-amount">金额 *</label>
                                <input type="number" id="payment-amount" name="amount" class="form-input" placeholder="请输入金额" min="0" step="0.01" required>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-outline" onclick="TaskDetailPage.hidePaymentModal()">取消</button>
                        <button class="btn btn-primary" onclick="TaskDetailPage.savePayment()">保存</button>
                    </div>
                </div>
            </div>
        `;
    },
    
    bindEvents() {
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal-overlay')) {
                this.hidePaymentModal();
            }
        });
    },
    
    async toggleComplete() {
        const newStatus = Store.state.tasks.find(t => t.id === this.taskId)?.status === 2 ? 0 : 2;
        
        try {
            const task = await TaskAPI.update(this.taskId, { status: newStatus });
            Store.updateTask(this.taskId, task);
            this.render({ id: this.taskId });
            Utils.showNotification(newStatus === 2 ? '任务已完成' : '任务已重开', 'success');
        } catch (error) {
            Utils.showNotification('操作失败: ' + error.message, 'danger');
        }
    },
    
    showPaymentModal() {
        const modal = document.getElementById('payment-modal');
        document.getElementById('payment-form').reset();
        modal.classList.add('active');
    },
    
    hidePaymentModal() {
        const modal = document.getElementById('payment-modal');
        modal.classList.remove('active');
    },
    
    async savePayment() {
        const form = document.getElementById('payment-form');
        
        const data = {
            budget_id: parseInt(form.budget_id.value),
            amount: parseFloat(form.amount.value)
        };
        
        try {
            await TaskAPI.setPayment(this.taskId, data);
            Utils.showNotification('付款已添加', 'success');
            this.hidePaymentModal();
            this.render({ id: this.taskId });
        } catch (error) {
            Utils.showNotification('保存失败: ' + error.message, 'danger');
        }
    },
    
    async deletePayment(paymentId) {
        if (!Utils.confirm('确定要删除这个付款记录吗？')) return;
        
        try {
            await TaskAPI.deletePayment(paymentId);
            Utils.showNotification('付款已删除', 'success');
            this.render({ id: this.taskId });
        } catch (error) {
            Utils.showNotification('删除失败: ' + error.message, 'danger');
        }
    }
};
