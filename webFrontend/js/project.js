// js/project.js
import { api } from './api.js';
import { requireAuth } from './auth.js';

requireAuth();

const urlParams = new URLSearchParams(window.location.search);
const projectId = urlParams.get('id');

// 加载项目详情
async function loadProject() {
const project = await api.get(`/projects/${projectId}`);
document.getElementById('projectName').textContent = project.name;
document.getElementById('projectDesc').textContent = project.description || '';

// 显示预算列表
const budgetList = document.getElementById('budgetList');
budgetList.innerHTML = (project.budgets || []).map(b => `
<tr>
  <td>${b.type}</td>
  <td>${b.budget}</td>
  <td>${b.used}</td>
  <td>${b.budget - b.used}</td>
</tr>
`).join('');

// 加载任务列表
const tasks = await api.get(`/projects/${projectId}/tasks?page=0&pageSize=50`);
const taskList = document.getElementById('taskList');
taskList.innerHTML = tasks.tasks.map(t => `
<tr>
  <td><a href="task.html?id=${t.id}">${t.name}</a></td>
  <td>${t.status}</td>
  <td>${t.deadline ? new Date(t.deadline).toLocaleString() : '-'}</td>
</tr>
`).join('');
}

// 创建任务
document.getElementById('createTaskForm')?.addEventListener('submit', async (e) => {
e.preventDefault();
const name = document.getElementById('taskName').value;
const description = document.getElementById('taskDesc').value;
const deadline = document.getElementById('taskDeadline').value;
await api.post('/tasks', {
project_id: parseInt(projectId),
name,
description,
deadline: deadline ? new Date(deadline).toISOString() : null
});
loadProject(); // 刷新
});

// 创建预算
document.getElementById('createBudgetForm')?.addEventListener('submit', async (e) => {
e.preventDefault();
const type = document.getElementById('budgetType').value;
const budget = parseFloat(document.getElementById('budgetAmount').value);
await api.post(`/projects/${projectId}/budgets`, { type, budget, used: 0 });
loadProject();
});

loadProject();