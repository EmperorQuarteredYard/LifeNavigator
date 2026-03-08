const Store = {
    getToken() {
        return localStorage.getItem(Config.TOKEN_KEY);
    },
    
    setToken(token) {
        localStorage.setItem(Config.TOKEN_KEY, token);
    },
    
    getRefreshToken() {
        return localStorage.getItem(Config.REFRESH_TOKEN_KEY);
    },
    
    setRefreshToken(token) {
        localStorage.setItem(Config.REFRESH_TOKEN_KEY, token);
    },
    
    getUser() {
        const user = localStorage.getItem(Config.USER_KEY);
        return user ? JSON.parse(user) : null;
    },
    
    setUser(user) {
        localStorage.setItem(Config.USER_KEY, JSON.stringify(user));
    },
    
    clearAuth() {
        localStorage.removeItem(Config.TOKEN_KEY);
        localStorage.removeItem(Config.REFRESH_TOKEN_KEY);
        localStorage.removeItem(Config.USER_KEY);
    },
    
    isAuthenticated() {
        return !!this.getToken();
    },
    
    state: {
        projects: [],
        tasks: [],
        accounts: [],
        inviteCodes: [],
        currentProject: null,
        currentTask: null,
        currentAccount: null,
        loading: false,
        error: null
    },
    
    listeners: [],
    
    subscribe(listener) {
        this.listeners.push(listener);
        return () => {
            this.listeners = this.listeners.filter(l => l !== listener);
        };
    },
    
    notify() {
        this.listeners.forEach(listener => listener(this.state));
    },
    
    setState(newState) {
        this.state = { ...this.state, ...newState };
        this.notify();
    },
    
    setProjects(projects) {
        this.state.projects = projects;
        this.notify();
    },
    
    addProject(project) {
        this.state.projects.unshift(project);
        this.notify();
    },
    
    updateProject(id, updates) {
        const index = this.state.projects.findIndex(p => p.id === id);
        if (index !== -1) {
            this.state.projects[index] = { ...this.state.projects[index], ...updates };
            this.notify();
        }
    },
    
    removeProject(id) {
        this.state.projects = this.state.projects.filter(p => p.id !== id);
        this.notify();
    },
    
    setTasks(tasks) {
        this.state.tasks = tasks;
        this.notify();
    },
    
    addTask(task) {
        this.state.tasks.unshift(task);
        this.notify();
    },
    
    updateTask(id, updates) {
        const index = this.state.tasks.findIndex(t => t.id === id);
        if (index !== -1) {
            this.state.tasks[index] = { ...this.state.tasks[index], ...updates };
            this.notify();
        }
    },
    
    removeTask(id) {
        this.state.tasks = this.state.tasks.filter(t => t.id !== id);
        this.notify();
    },
    
    setAccounts(accounts) {
        this.state.accounts = accounts;
        this.notify();
    },
    
    addAccount(account) {
        this.state.accounts.unshift(account);
        this.notify();
    },
    
    updateAccount(id, updates) {
        const index = this.state.accounts.findIndex(a => a.id === id);
        if (index !== -1) {
            this.state.accounts[index] = { ...this.state.accounts[index], ...updates };
            this.notify();
        }
    },
    
    removeAccount(id) {
        this.state.accounts = this.state.accounts.filter(a => a.id !== id);
        this.notify();
    },
    
    setInviteCodes(codes) {
        this.state.inviteCodes = codes;
        this.notify();
    },
    
    addInviteCode(code) {
        this.state.inviteCodes.unshift(code);
        this.notify();
    },
    
    setLoading(loading) {
        this.state.loading = loading;
        this.notify();
    },
    
    setError(error) {
        this.state.error = error;
        this.notify();
    },
    
    clearError() {
        this.state.error = null;
        this.notify();
    }
};
