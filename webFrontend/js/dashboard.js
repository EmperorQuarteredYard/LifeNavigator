// js/dashboard.js
import { api } from './api.js';
import { requireAuth, logout } from './auth.js';

requireAuth();

// 加载项目列表
async function loadProjects() {
const projects = await api.get('/projects?offset=0&limit=50');
const container = document.getElementById('projectList');
container.innerHTML = projects.map(p => `
<div class="project-card">
  <h3><a href="project.html?id=${p.id}">${p.name}</a></h3>
  <p>${p.description || '无描述'}</p>
</div>
`).join('');
}

// 创建项目
document.getElementById('createProjectForm').addEventListener('submit', async (e) => {
e.preventDefault();
const name = document.getElementById('projectName').value;
const description = document.getElementById('projectDesc').value;
await api.post('/projects', { name, description });
loadProjects(); // 刷新列表
});

loadProjects();

// 登出按钮
document.getElementById('logoutBtn').addEventListener('click', logout);