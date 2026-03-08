const AccountsPage = {
    render() {
        const page = document.getElementById('page-accounts');
        if (!page) return;
        
        const accounts = Store.state.accounts || [];
        const totalBalance = accounts.reduce((sum, a) => sum + (a.balance || 0), 0);
        
        page.innerHTML = `
            <div class="page-header flex-between">
                <h1>账户管理</h1>
                <a href="#accounts/create" class="btn btn-primary">新建账户</a>
            </div>
            
            <div class="card mb-3">
                <div class="card-body">
                    <div class="flex-between">
                        <span class="text-muted">总资产</span>
                        <span class="stat-value ${totalBalance >= 0 ? 'text-success' : 'text-danger'}">
                            ${Utils.formatCurrency(totalBalance)}
                        </span>
                    </div>
                </div>
            </div>
            
            ${accounts.length === 0 ? this.renderEmptyState() : this.renderAccountList(accounts)}
        `;
    },
    
    renderEmptyState() {
        return `
            <div class="card">
                <div class="empty-state">
                    <div class="empty-state-icon">💰</div>
                    <h3 class="empty-state-title">暂无账户</h3>
                    <p class="empty-state-description">创建你的第一个账户，开始管理财务</p>
                    <a href="#accounts/create" class="btn btn-primary">创建账户</a>
                </div>
            </div>
        `;
    },
    
    renderAccountList(accounts) {
        return `
            <div class="grid grid-3 gap-3">
                ${accounts.map(a => this.renderAccountCard(a)).join('')}
            </div>
        `;
    },
    
    renderAccountCard(account) {
        return `
            <div class="card account-card">
                <div class="card-header">
                    <h3 class="card-title">${Utils.escapeHtml(account.type)}</h3>
                    <div class="dropdown">
                        <button class="btn btn-sm btn-outline" onclick="AccountsPage.toggleDropdown(this)">⋮</button>
                        <div class="dropdown-menu">
                            <a class="dropdown-item" href="#accounts/edit/${account.id}">编辑</a>
                            <a class="dropdown-item" href="#" onclick="AccountsPage.adjustBalance(${account.id})">调整余额</a>
                            <div class="dropdown-divider"></div>
                            <span class="dropdown-item text-danger" onclick="AccountsPage.deleteAccount(${account.id})">删除</span>
                        </div>
                    </div>
                </div>
                <div class="card-body">
                    <div class="account-balance">
                        <div class="text-muted mb-1">当前余额</div>
                        <div class="balance-value ${account.balance >= 0 ? 'text-success' : 'text-danger'}">
                            ${Utils.formatCurrency(account.balance)}
                        </div>
                    </div>
                    <div class="account-actions mt-2">
                        <a href="#" class="btn btn-sm btn-outline" onclick="AccountsPage.viewLinkedTasks(${account.id})">关联任务</a>
                        <a href="#" class="btn btn-sm btn-outline" onclick="AccountsPage.viewLinkedPayments(${account.id})">付款记录</a>
                    </div>
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
    
    async deleteAccount(id) {
        if (!Utils.confirm('确定要删除这个账户吗？')) return;
        
        try {
            await AccountAPI.delete(id);
            Store.removeAccount(id);
            this.render();
            Utils.showNotification('账户已删除', 'success');
        } catch (error) {
            Utils.showNotification('删除失败: ' + error.message, 'danger');
        }
    },
    
    adjustBalance(id) {
        const amount = prompt('请输入调整金额（正数增加，负数减少）：');
        if (amount === null) return;
        
        const numAmount = parseFloat(amount);
        if (isNaN(numAmount)) {
            Utils.showNotification('请输入有效的金额', 'warning');
            return;
        }
        
        this.doAdjustBalance(id, numAmount);
    },
    
    async doAdjustBalance(id, amount) {
        try {
            const account = await AccountAPI.adjustBalance(id, amount);
            Store.updateAccount(id, account);
            this.render();
            Utils.showNotification('余额已调整', 'success');
        } catch (error) {
            Utils.showNotification('调整失败: ' + error.message, 'danger');
        }
    },
    
    async viewLinkedTasks(id) {
        try {
            const tasks = await AccountAPI.getLinkedTasks(id);
            
            const modal = document.createElement('div');
            modal.className = 'modal-overlay active';
            modal.innerHTML = `
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">关联任务</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button>
                    </div>
                    <div class="modal-body">
                        ${tasks.length === 0 ? '<p class="text-muted text-center">暂无关联任务</p>' : `
                            <table class="table">
                                <thead>
                                    <tr>
                                        <th>任务名称</th>
                                        <th>状态</th>
                                        <th>截止时间</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${tasks.map(t => `
                                        <tr>
                                            <td><a href="#tasks/detail/${t.id}">${Utils.escapeHtml(t.name)}</a></td>
                                            <td>
                                                <span class="badge ${Utils.getStatusClass(t.status, Config.TASK_STATUS)}">
                                                    ${Utils.getStatusLabel(t.status, Config.TASK_STATUS)}
                                                </span>
                                            </td>
                                            <td>${t.deadline ? Utils.formatDateShort(t.deadline) : '-'}</td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        `}
                    </div>
                </div>
            `;
            
            document.body.appendChild(modal);
            modal.addEventListener('click', (e) => {
                if (e.target === modal) modal.remove();
            });
            
        } catch (error) {
            Utils.showNotification('加载失败: ' + error.message, 'danger');
        }
    },
    
    async viewLinkedPayments(id) {
        try {
            const payments = await AccountAPI.getLinkedPayments(id);
            
            const modal = document.createElement('div');
            modal.className = 'modal-overlay active';
            modal.innerHTML = `
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">付款记录</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">&times;</button>
                    </div>
                    <div class="modal-body">
                        ${payments.length === 0 ? '<p class="text-muted text-center">暂无付款记录</p>' : `
                            <table class="table">
                                <thead>
                                    <tr>
                                        <th>任务ID</th>
                                        <th>预算ID</th>
                                        <th>金额</th>
                                        <th>时间</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${payments.map(p => `
                                        <tr>
                                            <td>${p.task_id}</td>
                                            <td>${p.budget_id}</td>
                                            <td>${Utils.formatCurrency(p.amount)}</td>
                                            <td>${Utils.formatDate(p.created_at)}</td>
                                        </tr>
                                    `).join('')}
                                </tbody>
                            </table>
                        `}
                    </div>
                </div>
            `;
            
            document.body.appendChild(modal);
            modal.addEventListener('click', (e) => {
                if (e.target === modal) modal.remove();
            });
            
        } catch (error) {
            Utils.showNotification('加载失败: ' + error.message, 'danger');
        }
    }
};

const AccountFormPage = {
    isEdit: false,
    accountId: null,
    
    render(params = {}) {
        const page = document.getElementById('page-account-form');
        if (!page) return;
        
        this.isEdit = params.id ? true : false;
        this.accountId = params.id || null;
        
        page.innerHTML = `
            <div class="page-header">
                <h1>${this.isEdit ? '编辑账户' : '新建账户'}</h1>
            </div>
            
            <div class="card">
                <form id="account-form">
                    <div class="form-group">
                        <label class="form-label" for="type">账户类型 *</label>
                        <input type="text" id="type" name="type" class="form-input" placeholder="如：现金、银行卡、支付宝、微信" required>
                    </div>
                    <div class="form-group">
                        <label class="form-label" for="balance">初始余额</label>
                        <input type="number" id="balance" name="balance" class="form-input" placeholder="请输入初始余额" min="0" step="0.01" value="0">
                    </div>
                    <div id="form-error" class="form-error hidden"></div>
                    <div class="flex gap-2">
                        <button type="submit" class="btn btn-primary">${this.isEdit ? '保存' : '创建'}</button>
                        <a href="#accounts" class="btn btn-outline">取消</a>
                    </div>
                </form>
            </div>
        `;
        
        this.bindEvents();
        
        if (this.isEdit) {
            this.loadAccount();
        }
    },
    
    bindEvents() {
        const form = document.getElementById('account-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleSubmit(form);
            });
        }
    },
    
    async loadAccount() {
        try {
            const account = await AccountAPI.get(this.accountId);
            
            document.getElementById('type').value = account.type || '';
            document.getElementById('balance').value = account.balance || 0;
            
        } catch (error) {
            Utils.showNotification('加载账户失败: ' + error.message, 'danger');
            Router.navigate('accounts');
        }
    },
    
    async handleSubmit(form) {
        const errorEl = document.getElementById('form-error');
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const type = form.type.value.trim();
        const balance = parseFloat(form.balance.value) || 0;
        
        if (!type) {
            errorEl.textContent = '请输入账户类型';
            errorEl.classList.remove('hidden');
            return;
        }
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '保存中...';
            errorEl.classList.add('hidden');
            
            const data = { type, balance };
            
            if (this.isEdit) {
                const account = await AccountAPI.update(this.accountId, data);
                Store.updateAccount(this.accountId, account);
                Utils.showNotification('账户已更新', 'success');
            } else {
                const account = await AccountAPI.create(data);
                Store.addAccount(account);
                Utils.showNotification('账户创建成功', 'success');
            }
            
            Router.navigate('accounts');
            
        } catch (error) {
            errorEl.textContent = error.message;
            errorEl.classList.remove('hidden');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = this.isEdit ? '保存' : '创建';
        }
    }
};
