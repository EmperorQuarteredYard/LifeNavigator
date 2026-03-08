const Config = {
    API_BASE_URL: 'http://localhost:8080/api',
    TOKEN_KEY: 'life_navigator_token',
    REFRESH_TOKEN_KEY: 'life_navigator_refresh_token',
    USER_KEY: 'life_navigator_user',
    
    ENDPOINTS: {
        AUTH: {
            REGISTER: '/v0/register',
            LOGIN: '/v0/login',
            REFRESH: '/v0/refresh',
            PROFILE: '/v0/user/profile',
            USER: '/v0/users'
        },
        PROJECTS: {
            BASE: '/v0/projects',
            BUDGETS: '/v0/projects/budgets'
        },
        TASKS: {
            BASE: '/v0/tasks',
            BUDGETS: '/v0/tasks/budgets'
        },
        ACCOUNTS: {
            BASE: '/v0/accounts',
            BALANCE: '/balance',
            TASKS: '/tasks',
            PAYMENTS: '/payments'
        },
        INVITE_CODES: {
            BASE: '/v0/invite-codes',
            USER_CODES: '/v0/users'
        },
        AI: {
            REDUCE_PROJECT: '/v0/ai/reduce-project',
            SUMMARY: '/v0/ai/summary'
        }
    },
    
    TASK_STATUS: {
        0: { label: '待办', class: 'badge-secondary' },
        1: { label: '进行中', class: 'badge-primary' },
        2: { label: '已完成', class: 'badge-success' },
        3: { label: '已取消', class: 'badge-danger' }
    },
    
    INVITE_CODE_STATUS: {
        0: { label: '无效', class: 'badge-danger' },
        1: { label: '有效', class: 'badge-success' }
    }
};
