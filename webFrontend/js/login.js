<!-- login.html 节选 -->
<form id="loginForm">
    <input type="text" id="username" placeholder="用户名" required>
    <input type="password" id="password" placeholder="密码" required>
    <button type="submit">登录</button>
</form>
<script type="module">
    import { login } from './js/auth.js';
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        try {
            await login(username, password);
            window.location.href = '/dashboard.html';
        } catch (err) {
            alert('登录失败：' + err.message);
        }
    });
</script>