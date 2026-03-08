const API = {
    async request(endpoint, options = {}) {
        const url = Config.API_BASE_URL + endpoint;
        const token = Store.getToken();
        
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };
        
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        try {
            const response = await fetch(url, {
                ...options,
                headers
            });
            
            if (response.status === 401) {
                const refreshed = await this.refreshToken();
                if (refreshed) {
                    headers['Authorization'] = `Bearer ${Store.getToken()}`;
                    const retryResponse = await fetch(url, {
                        ...options,
                        headers
                    });
                    return this.handleResponse(retryResponse);
                } else {
                    Store.clearAuth();
                    Router.navigate('login');
                    throw new Error('认证已过期，请重新登录');
                }
            }
            
            return this.handleResponse(response);
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    },
    
    async handleResponse(response) {
        const contentType = response.headers.get('content-type');
        
        if (contentType && contentType.includes('application/json')) {
            const data = await response.json();
            
            console.log('[API Response]', {
                url: response.url,
                status: response.status,
                data: data
            });
            
            if (!response.ok) {
                throw new Error(data.message || data.error || `请求失败 (${response.status})`);
            }
            
            if (data.code !== undefined && data.code !== 0) {
                throw new Error(data.message || `请求失败 (code: ${data.code})`);
            }
            
            return data.data !== undefined ? data.data : data;
        }
        
        if (!response.ok) {
            throw new Error(`请求失败 (${response.status})`);
        }
        
        return response;
    },
    
    async refreshToken() {
        const refreshToken = Store.getRefreshToken();
        if (!refreshToken) return false;
        
        try {
            const response = await fetch(Config.API_BASE_URL + '/v0/refresh', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ refresh_token: refreshToken })
            });
            
            if (response.ok) {
                const data = await response.json();
                console.log('[Token Refresh Response]', data);
                Store.setToken(data.access_token);
                Store.setRefreshToken(data.refresh_token);
                return true;
            }
        } catch (error) {
            console.error('Token refresh failed:', error);
        }
        
        return false;
    },
    
    get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    },
    
    post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    },
    
    put(endpoint, data) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    },
    
    delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    },
    
    async sseRequest(endpoint, data, callbacks) {
        const url = Config.API_BASE_URL + endpoint;
        const token = Store.getToken();
        
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'SSE请求失败');
        }
        
        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';
        
        while (true) {
            const { done, value } = await reader.read();
            
            if (done) {
                if (callbacks.onComplete) {
                    callbacks.onComplete();
                }
                break;
            }
            
            buffer += decoder.decode(value, { stream: true });
            const lines = buffer.split('\n');
            buffer = lines.pop() || '';
            
            for (const line of lines) {
                if (line.startsWith('data: ')) {
                    const data = line.slice(6);
                    
                    if (data === '[STREAM_END]') {
                        if (callbacks.onStreamEnd) {
                            callbacks.onStreamEnd();
                        }
                        return;
                    }
                    
                    try {
                        const event = JSON.parse(data);
                        console.log('[SSE Event]', event);
                        if (callbacks.onEvent) {
                            callbacks.onEvent(event);
                        }
                    } catch (e) {
                        console.error('Failed to parse SSE event:', data);
                    }
                }
            }
        }
    }
};

const AuthAPI = {
    register(data) {
        return API.post('/v0/register', data);
    },
    
    login(data) {
        return API.post('/v0/login', data);
    },
    
    refresh() {
        return API.post('/v0/refresh', {
            refresh_token: Store.getRefreshToken()
        });
    },
    
    getProfile() {
        return API.get('/v0/user/profile');
    },
    
    getUser(id) {
        return API.get(`/v0/users/${id}`);
    }
};

const ProjectAPI = {
    create(data) {
        return API.post('/v0/projects', data);
    },
    
    get(id) {
        return API.get(`/v0/projects/${id}`);
    },
    
    list(params = {}) {
        const query = Utils.buildQueryString(params);
        return API.get(`/v0/projects${query ? '?' + query : ''}`);
    },
    
    update(id, data) {
        return API.put(`/v0/projects/${id}`, data);
    },
    
    delete(id) {
        return API.delete(`/v0/projects/${id}`);
    },
    
    addBudget(projectId, data) {
        return API.post(`/v0/projects/${projectId}/budgets`, data);
    },
    
    updateBudget(budgetId, data) {
        return API.put(`/v0/projects/budgets/${budgetId}`, data);
    },
    
    deleteBudget(budgetId) {
        return API.delete(`/v0/projects/budgets/${budgetId}`);
    }
};

const TaskAPI = {
    create(data) {
        return API.post('/v0/tasks', data);
    },
    
    get(id) {
        return API.get(`/v0/tasks/${id}`);
    },
    
    list(params = {}) {
        const query = Utils.buildQueryString(params);
        return API.get(`/v0/tasks${query ? '?' + query : ''}`);
    },
    
    listByProject(projectId, params = {}) {
        const query = Utils.buildQueryString(params);
        return API.get(`/v0/projects/${projectId}/tasks${query ? '?' + query : ''}`);
    },
    
    update(id, data) {
        return API.put(`/v0/tasks/${id}`, data);
    },
    
    delete(id) {
        return API.delete(`/v0/tasks/${id}`);
    },
    
    setPayment(taskId, data) {
        return API.post(`/v0/tasks/${taskId}/budgets`, data);
    },
    
    updatePayment(paymentId, data) {
        return API.put(`/v0/tasks/budgets/${paymentId}`, data);
    },
    
    deletePayment(paymentId) {
        return API.delete(`/v0/tasks/budgets/${paymentId}`);
    },
    
    getPayments(taskId) {
        return API.get(`/v0/tasks/${taskId}/budgets`);
    }
};

const AccountAPI = {
    create(data) {
        return API.post('/v0/accounts', data);
    },
    
    list() {
        return API.get('/v0/accounts');
    },
    
    get(id) {
        return API.get(`/v0/accounts/${id}`);
    },
    
    delete(id) {
        return API.delete(`/v0/accounts/${id}`);
    },
    
    adjustBalance(id, amount) {
        return API.post(`/v0/accounts/${id}/balance`, { amount });
    },
    
    getLinkedTasks(id) {
        return API.get(`/v0/accounts/${id}/tasks`);
    },
    
    getLinkedPayments(id) {
        return API.get(`/v0/accounts/${id}/payments`);
    }
};

const InviteCodeAPI = {
    create(data) {
        return API.post('/v0/invite-codes', data);
    },
    
    get(token) {
        return API.get(`/v0/invite-codes/${token}`);
    },
    
    listByUser(userId) {
        return API.get(`/v0/users/${userId}/invite-codes`);
    }
};

const AIAPI = {
    reduceProject(description, callbacks) {
        return API.sseRequest('/v0/ai/reduce-project', {
            project_description: description
        }, callbacks);
    },
    
    summary(startTime, endTime, callbacks) {
        return API.sseRequest('/v0/ai/summary', {
            start_time: startTime,
            end_time: endTime
        }, callbacks);
    }
};
