// js/task.js
import { api } from './api.js';
import { requireAuth } from './auth.js';

requireAuth();

const urlParams = new URLSearchParams(window.location.search);
const taskId = urlParams.get('id');

async function loadTask() {
const task = await api.get(`/tasks/${taskId}`);
document.getElementById('taskName').textContent = task.name;
document.getElementById('taskDesc').textContent = task.description || '-';
document.getElementById('taskStatus').textContent = task.status;
document.getElementById('taskDeadline').textContent = task.deadline ? new Date(task.deadline).toLocaleString() : '-';

// 加载任务预算
const budgets = await api.get(`/tasks/${taskId}/budgets`);
const budgetList = document.getElementById('taskBudgetList');
budgetList.innerHTML = budgets.map(b => `
<tr>
    <td>${b.type}</td>
    <td>${b.budget}</td>
    <td>${b.used}</td>
    <td>${b.budget - b.used}</td>
    <td><button onclick="deleteBudget(${b.id})">删除</button></td>
</tr>
`).join('');
}

// 更新任务状态
window.updateStatus = async function(status) {
await api.put(`/tasks/${taskId}`, { status });
loadTask();
};

// 添加任务预算
document.getElementById('addTaskBudgetForm')?.addEventListener('submit', async (e) => {
e.preventDefault();
const type = document.getElementById('budgetType').value;
const budget = parseFloat(document.getElementById('budgetAmount').value);
await api.post(`/tasks/${taskId}/budgets`, { type, budget, used: 0 });
loadTask();
});

// 删除任务预算（全局函数）
window.deleteBudget = async function(budgetId) {
if (confirm('确定删除该预算项？')) {
await api.delete(`/tasks/budgets/${budgetId}`);
loadTask();
}
};

loadTask();