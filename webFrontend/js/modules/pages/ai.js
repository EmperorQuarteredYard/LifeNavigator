const AIPage = {
    currentTab: 'reduce',
    isStreaming: false,
    
    render() {
        const page = document.getElementById('page-ai');
        if (!page) return;
        
        page.innerHTML = `
            <div class="page-header">
                <h1>AI助手</h1>
            </div>
            
            <div class="card mb-3">
                <div class="tabs">
                    <span class="tab ${this.currentTab === 'reduce' ? 'active' : ''}" data-tab="reduce">AI创建项目</span>
                    <span class="tab ${this.currentTab === 'summary' ? 'active' : ''}" data-tab="summary">AI总结成就</span>
                </div>
            </div>
            
            <div id="ai-content">
                ${this.currentTab === 'reduce' ? this.renderReduceForm() : this.renderSummaryForm()}
            </div>
        `;
        
        this.bindEvents();
    },
    
    renderReduceForm() {
        return `
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">AI辅助创建项目</h3>
                </div>
                <div class="card-body">
                    <p class="text-muted mb-3">描述你想要创建的项目，AI会自动为你生成项目、任务和预算。</p>
                    <form id="ai-reduce-form">
                        <div class="form-group">
                            <label class="form-label" for="project-description">项目描述 *</label>
                            <textarea id="project-description" name="description" class="form-textarea" 
                                placeholder="例如：我想创建一个学习Go语言的计划，包括基础语法、并发编程、Web开发等内容，预计需要3个月时间，预算500元用于购买书籍和课程。"
                                rows="5" required></textarea>
                        </div>
                        <div class="flex gap-2">
                            <button type="submit" class="btn btn-primary" id="reduce-submit-btn">
                                <span class="btn-text">开始生成</span>
                                <span class="loading hidden"></span>
                            </button>
                            <button type="button" class="btn btn-outline hidden" id="reduce-stop-btn">停止</button>
                        </div>
                    </form>
                    
                    <div id="reduce-output" class="mt-3 hidden">
                        <h4 class="mb-2">生成结果</h4>
                        <div class="sse-output" id="reduce-sse-output"></div>
                    </div>
                </div>
            </div>
        `;
    },
    
    renderSummaryForm() {
        const now = new Date();
        const oneMonthAgo = new Date(now.getFullYear(), now.getMonth() - 1, now.getDate());
        
        return `
            <div class="card">
                <div class="card-header">
                    <h3 class="card-title">AI总结成就</h3>
                </div>
                <div class="card-body">
                    <p class="text-muted mb-3">选择时间范围，AI会总结你在这段时间内完成的任务和成就。</p>
                    <form id="ai-summary-form">
                        <div class="grid grid-2 gap-2">
                            <div class="form-group">
                                <label class="form-label" for="start-time">开始时间 *</label>
                                <input type="datetime-local" id="start-time" name="start_time" class="form-input" required
                                    value="${oneMonthAgo.toISOString().slice(0, 16)}">
                            </div>
                            <div class="form-group">
                                <label class="form-label" for="end-time">结束时间 *</label>
                                <input type="datetime-local" id="end-time" name="end_time" class="form-input" required
                                    value="${now.toISOString().slice(0, 16)}">
                            </div>
                        </div>
                        <div class="flex gap-2">
                            <button type="submit" class="btn btn-primary" id="summary-submit-btn">
                                <span class="btn-text">开始总结</span>
                                <span class="loading hidden"></span>
                            </button>
                            <button type="button" class="btn btn-outline hidden" id="summary-stop-btn">停止</button>
                        </div>
                    </form>
                    
                    <div id="summary-output" class="mt-3 hidden">
                        <h4 class="mb-2">总结结果</h4>
                        <div class="sse-output" id="summary-sse-output"></div>
                    </div>
                </div>
            </div>
        `;
    },
    
    bindEvents() {
        const tabs = document.querySelectorAll('.tab');
        tabs.forEach(tab => {
            tab.addEventListener('click', () => {
                tabs.forEach(t => t.classList.remove('active'));
                tab.classList.add('active');
                this.currentTab = tab.dataset.tab;
                
                const content = document.getElementById('ai-content');
                content.innerHTML = this.currentTab === 'reduce' 
                    ? this.renderReduceForm() 
                    : this.renderSummaryForm();
                
                this.bindFormEvents();
            });
        });
        
        this.bindFormEvents();
    },
    
    bindFormEvents() {
        const reduceForm = document.getElementById('ai-reduce-form');
        if (reduceForm) {
            reduceForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handleReduceProject();
            });
        }
        
        const summaryForm = document.getElementById('ai-summary-form');
        if (summaryForm) {
            summaryForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.handleSummary();
            });
        }
    },
    
    async handleReduceProject() {
        const description = document.getElementById('project-description').value.trim();
        const submitBtn = document.getElementById('reduce-submit-btn');
        const stopBtn = document.getElementById('reduce-stop-btn');
        const outputDiv = document.getElementById('reduce-output');
        const sseOutput = document.getElementById('reduce-sse-output');
        
        if (!description) {
            Utils.showNotification('请输入项目描述', 'warning');
            return;
        }
        
        this.isStreaming = true;
        submitBtn.querySelector('.btn-text').textContent = '生成中...';
        submitBtn.querySelector('.loading').classList.remove('hidden');
        submitBtn.disabled = true;
        stopBtn.classList.remove('hidden');
        outputDiv.classList.remove('hidden');
        sseOutput.innerHTML = '';
        
        try {
            await AIAPI.reduceProject(description, {
                onEvent: (event) => {
                    this.handleReduceEvent(event, sseOutput);
                },
                onStreamEnd: () => {
                    this.appendOutput(sseOutput, '\n[流结束]', 'event-stream_complete');
                },
                onComplete: () => {
                    this.finishStreaming('reduce');
                }
            });
        } catch (error) {
            this.appendOutput(sseOutput, `\n错误: ${error.message}`, 'event-stream_error');
            this.finishStreaming('reduce');
        }
    },
    
    async handleSummary() {
        const startTime = document.getElementById('start-time').value;
        const endTime = document.getElementById('end-time').value;
        const submitBtn = document.getElementById('summary-submit-btn');
        const stopBtn = document.getElementById('summary-stop-btn');
        const outputDiv = document.getElementById('summary-output');
        const sseOutput = document.getElementById('summary-sse-output');
        
        if (!startTime || !endTime) {
            Utils.showNotification('请选择时间范围', 'warning');
            return;
        }
        
        this.isStreaming = true;
        submitBtn.querySelector('.btn-text').textContent = '总结中...';
        submitBtn.querySelector('.loading').classList.remove('hidden');
        submitBtn.disabled = true;
        stopBtn.classList.remove('hidden');
        outputDiv.classList.remove('hidden');
        sseOutput.innerHTML = '';
        
        try {
            await AIAPI.summary(
                new Date(startTime).toISOString(),
                new Date(endTime).toISOString(),
                {
                    onEvent: (event) => {
                        this.handleSummaryEvent(event, sseOutput);
                    },
                    onStreamEnd: () => {
                        this.appendOutput(sseOutput, '\n[流结束]', 'event-stream_complete');
                    },
                    onComplete: () => {
                        this.finishStreaming('summary');
                    }
                }
            );
        } catch (error) {
            this.appendOutput(sseOutput, `\n错误: ${error.message}`, 'event-stream_error');
            this.finishStreaming('summary');
        }
    },
    
    handleReduceEvent(event, output) {
        switch (event.type) {
            case 'project_created':
                this.appendOutput(output, 
                    `\n[项目创建] ID: ${event.content.id}, 名称: ${event.content.name}`,
                    'event-project_created'
                );
                break;
            case 'task_created':
                this.appendOutput(output,
                    `\n[任务创建] ID: ${event.content.id}, 名称: ${event.content.name}`,
                    'event-task_created'
                );
                break;
            case 'budget_created':
                this.appendOutput(output,
                    `\n[预算创建] ID: ${event.content.id}, 类型: ${event.content.type}, 金额: ${Utils.formatCurrency(event.content.budget)}`,
                    'event-budget_created'
                );
                break;
            case 'stream_complete':
                this.appendOutput(output,
                    `\n[完成] ${event.content.message}`,
                    'event-stream_complete'
                );
                break;
            case 'stream_error':
                this.appendOutput(output,
                    `\n[错误] ${event.content.error}`,
                    'event-stream_error'
                );
                break;
        }
    },
    
    handleSummaryEvent(event, output) {
        switch (event.type) {
            case 'summary_content':
                this.appendOutput(output, event.content.content, 'event-summary_content');
                break;
            case 'stream_complete':
                this.appendOutput(output,
                    `\n\n[完成] ${event.content.message}`,
                    'event-stream_complete'
                );
                break;
            case 'stream_error':
                this.appendOutput(output,
                    `\n[错误] ${event.content.error}`,
                    'event-stream_error'
                );
                break;
        }
    },
    
    appendOutput(element, text, className = '') {
        const span = document.createElement('span');
        span.className = className;
        span.textContent = text;
        element.appendChild(span);
        element.scrollTop = element.scrollHeight;
    },
    
    finishStreaming(type) {
        this.isStreaming = false;
        
        if (type === 'reduce') {
            const submitBtn = document.getElementById('reduce-submit-btn');
            const stopBtn = document.getElementById('reduce-stop-btn');
            submitBtn.querySelector('.btn-text').textContent = '开始生成';
            submitBtn.querySelector('.loading').classList.add('hidden');
            submitBtn.disabled = false;
            stopBtn.classList.add('hidden');
            
            App.loadInitialData();
        } else {
            const submitBtn = document.getElementById('summary-submit-btn');
            const stopBtn = document.getElementById('summary-stop-btn');
            submitBtn.querySelector('.btn-text').textContent = '开始总结';
            submitBtn.querySelector('.loading').classList.add('hidden');
            submitBtn.disabled = false;
            stopBtn.classList.add('hidden');
        }
    }
};
