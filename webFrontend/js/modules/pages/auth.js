const LoginPage = {
    render() {
        const container = document.getElementById('login-form-container');
        if (!container) return;
        
        container.innerHTML = `
            <form id="login-form" class="auth-form">
                <div class="form-group">
                    <label class="form-label" for="username">用户名</label>
                    <input type="text" id="username" name="username" class="form-input" placeholder="请输入用户名" required>
                </div>
                <div class="form-group">
                    <label class="form-label" for="password">密码</label>
                    <input type="password" id="password" name="password" class="form-input" placeholder="请输入密码" required>
                </div>
                <div id="login-error" class="form-error hidden"></div>
                <button type="submit" class="btn btn-primary btn-block btn-lg">登录</button>
                <div class="auth-footer">
                    <p>还没有账号？<a href="#register">立即注册</a></p>
                </div>
            </form>
        `;
        
        this.bindEvents();
    },
    
    bindEvents() {
        const form = document.getElementById('login-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleLogin(form);
            });
        }
    },
    
    async handleLogin(form) {
        const errorEl = document.getElementById('login-error');
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const username = form.username.value.trim();
        const password = form.password.value;
        
        if (!username || !password) {
            errorEl.textContent = '请填写用户名和密码';
            errorEl.classList.remove('hidden');
            return;
        }
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '登录中...';
            errorEl.classList.add('hidden');
            
            const response = await AuthAPI.login({ username, password });
            
            Store.setToken(response.access_token);
            Store.setRefreshToken(response.refresh_token);
            Store.setUser(response.user);
            
            App.renderLayout();
            Router.navigate('dashboard');
            Utils.showNotification('登录成功', 'success');
            
        } catch (error) {
            errorEl.textContent = error.message;
            errorEl.classList.remove('hidden');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = '登录';
        }
    }
};

const RegisterPage = {
    render() {
        const container = document.getElementById('register-form-container');
        if (!container) return;
        
        container.innerHTML = `
            <form id="register-form" class="auth-form">
                <div class="form-group">
                    <label class="form-label" for="username">用户名 *</label>
                    <input type="text" id="username" name="username" class="form-input" placeholder="请输入用户名" required>
                </div>
                <div class="form-group">
                    <label class="form-label" for="password">密码 *</label>
                    <input type="password" id="password" name="password" class="form-input" placeholder="请输入密码" required>
                </div>
                <div class="form-group">
                    <label class="form-label" for="confirmPassword">确认密码 *</label>
                    <input type="password" id="confirmPassword" name="confirmPassword" class="form-input" placeholder="请再次输入密码" required>
                </div>
                <div class="form-group">
                    <label class="form-label" for="nickname">昵称</label>
                    <input type="text" id="nickname" name="nickname" class="form-input" placeholder="请输入昵称（可选）">
                </div>
                <div class="form-group">
                    <label class="form-label" for="email">邮箱</label>
                    <input type="email" id="email" name="email" class="form-input" placeholder="请输入邮箱（可选）">
                </div>
                <div class="form-group">
                    <label class="form-label" for="phone">手机号</label>
                    <input type="tel" id="phone" name="phone" class="form-input" placeholder="请输入手机号（可选）">
                </div>
                <div class="form-group">
                    <label class="form-label" for="inviteCode">邀请码 *</label>
                    <input type="text" id="inviteCode" name="inviteCode" class="form-input" placeholder="请输入邀请码" required>
                </div>
                <div id="register-error" class="form-error hidden"></div>
                <button type="submit" class="btn btn-primary btn-block btn-lg">注册</button>
                <div class="auth-footer">
                    <p>已有账号？<a href="#login">立即登录</a></p>
                </div>
            </form>
        `;
        
        this.bindEvents();
    },
    
    bindEvents() {
        const form = document.getElementById('register-form');
        if (form) {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                await this.handleRegister(form);
            });
        }
    },
    
    async handleRegister(form) {
        const errorEl = document.getElementById('register-error');
        const submitBtn = form.querySelector('button[type="submit"]');
        
        const username = form.username.value.trim();
        const password = form.password.value;
        const confirmPassword = form.confirmPassword.value;
        const nickname = form.nickname.value.trim();
        const email = form.email.value.trim();
        const phone = form.phone.value.trim();
        const inviteCode = form.inviteCode.value.trim();
        
        if (!username || !password || !confirmPassword || !inviteCode) {
            errorEl.textContent = '请填写所有必填项';
            errorEl.classList.remove('hidden');
            return;
        }
        
        if (password !== confirmPassword) {
            errorEl.textContent = '两次输入的密码不一致';
            errorEl.classList.remove('hidden');
            return;
        }
        
        if (password.length < 6) {
            errorEl.textContent = '密码长度至少6位';
            errorEl.classList.remove('hidden');
            return;
        }
        
        if (email && !Utils.validateEmail(email)) {
            errorEl.textContent = '邮箱格式不正确';
            errorEl.classList.remove('hidden');
            return;
        }
        
        if (phone && !Utils.validatePhone(phone)) {
            errorEl.textContent = '手机号格式不正确';
            errorEl.classList.remove('hidden');
            return;
        }
        
        try {
            submitBtn.disabled = true;
            submitBtn.textContent = '注册中...';
            errorEl.classList.add('hidden');
            
            const response = await AuthAPI.register({
                username,
                password,
                nickname: nickname || username,
                email: email || undefined,
                phone: phone || undefined,
                invite_code: inviteCode
            });
            
            Store.setToken(response.access_token);
            Store.setRefreshToken(response.refresh_token);
            Store.setUser(response.user);
            
            App.renderLayout();
            Router.navigate('dashboard');
            Utils.showNotification('注册成功', 'success');
            
        } catch (error) {
            errorEl.textContent = error.message;
            errorEl.classList.remove('hidden');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = '注册';
        }
    }
};
