const InviteCodesPage = {
    render() {
        const page = document.getElementById('page-invite-codes');
        if (!page) return;
        
        const user = Store.getUser();
        const codes = Store.state.inviteCodes || [];
        
        page.innerHTML = `
            <div class="page-header flex-between">
                <h1>邀请码管理</h1>
                <button class="btn btn-primary" onclick="InviteCodesPage.showCreateModal()">创建邀请码</button>
            </div>
            
            <div class="card">
                <div class="card-body">
                    ${codes.length === 0 ? this.renderEmptyState() : this.renderCodeList(codes)}
                </div>
            </div>
            
            <div class="modal-overlay" id="invite-modal">
                <div class="modal">
                    <div class="modal-header">
                        <h3 class="modal-title">创建邀请码</h3>
                        <button class="modal-close" onclick="InviteCodesPage.hideModal()">&times;</button>
                    </div>
                    <div class="modal-body">
                        <form id="invite-form">
                            <div class="form-group">
                                <label class="form-label" for="invite-amount">数量 *</label>
                                <input type="number" id="invite-amount" name="amount" class="form-input" 
                                    placeholder="请输入创建数量" min="1" max="10" value="1" required>
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="invite-role">角色 *</label>
                                <select id="invite-role" name="role" class="form-select" required>
                                    <option value="user">普通用户</option>
                                    <option value="admin">管理员</option>
                                </select>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-outline" onclick="InviteCodesPage.hideModal()">取消</button>
                        <button class="btn btn-primary" onclick="InviteCodesPage.createCodes()">创建</button>
                    </div>
                </div>
            </div>
        `;
        
        this.loadCodes();
    },
    
    renderEmptyState() {
        return `
            <div class="empty-state">
                <div class="empty-state-icon">🎫</div>
                <h3 class="empty-state-title">暂无邀请码</h3>
                <p class="empty-state-description">创建邀请码分享给朋友，邀请他们加入</p>
            </div>
        `;
    },
    
    renderCodeList(codes) {
        return `
            <table class="table">
                <thead>
                    <tr>
                        <th>邀请码</th>
                        <th>角色</th>
                        <th>使用情况</th>
                        <th>状态</th>
                        <th>过期时间</th>
                        <th>操作</th>
                    </tr>
                </thead>
                <tbody>
                    ${codes.map(c => `
                        <tr>
                            <td>
                                <code>${Utils.escapeHtml(c.token)}</code>
                            </td>
                            <td>${Utils.escapeHtml(c.invite_as_role || 'user')}</td>
                            <td>${c.count || 0} / ${c.amount === -1 ? '∞' : c.amount}</td>
                            <td>
                                <span class="badge ${Utils.getStatusClass(c.status, Config.INVITE_CODE_STATUS)}">
                                    ${Utils.getStatusLabel(c.status, Config.INVITE_CODE_STATUS)}
                                </span>
                            </td>
                            <td>${c.expires_at ? Utils.formatDateShort(c.expires_at) : '永不过期'}</td>
                            <td>
                                <button class="btn btn-sm btn-outline" onclick="InviteCodesPage.copyCode('${Utils.escapeHtml(c.token)}')">
                                    复制
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    },
    
    async loadCodes() {
        const user = Store.getUser();
        if (!user) return;
        
        try {
            const codes = await InviteCodeAPI.listByUser(user.id);
            Store.setInviteCodes(codes);
            
            const page = document.getElementById('page-invite-codes');
            if (page) {
                const cardBody = page.querySelector('.card-body');
                if (cardBody) {
                    cardBody.innerHTML = codes.length === 0 
                        ? this.renderEmptyState() 
                        : this.renderCodeList(codes);
                }
            }
        } catch (error) {
            console.error('Failed to load invite codes:', error);
        }
    },
    
    showCreateModal() {
        const modal = document.getElementById('invite-modal');
        document.getElementById('invite-form').reset();
        modal.classList.add('active');
    },
    
    hideModal() {
        const modal = document.getElementById('invite-modal');
        modal.classList.remove('active');
    },
    
    async createCodes() {
        const form = document.getElementById('invite-form');
        
        const amount = parseInt(form.amount.value);
        const role = form.role.value;
        
        if (!amount || amount < 1) {
            Utils.showNotification('请输入有效的数量', 'warning');
            return;
        }
        
        try {
            const codes = await InviteCodeAPI.create({ amount, role });
            Store.setInviteCodes([...(Store.state.inviteCodes || []), ...codes]);
            
            this.hideModal();
            this.loadCodes();
            Utils.showNotification(`成功创建 ${codes.length} 个邀请码`, 'success');
        } catch (error) {
            Utils.showNotification('创建失败: ' + error.message, 'danger');
        }
    },
    
    async copyCode(token) {
        const success = await Utils.copyToClipboard(token);
        if (success) {
            Utils.showNotification('邀请码已复制到剪贴板', 'success');
        } else {
            Utils.showNotification('复制失败，请手动复制', 'danger');
        }
    }
};
