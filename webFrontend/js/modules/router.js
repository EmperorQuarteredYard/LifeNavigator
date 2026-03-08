const Router = {
    routes: {},
    currentRoute: null,
    beforeNavigate: null,
    
    init(routes) {
        this.routes = routes;
        
        window.addEventListener('popstate', (event) => {
            if (event.state && event.state.route) {
                this.loadRoute(event.state.route, event.state.params, false);
            } else {
                this.handleUrlChange();
            }
        });
        
        this.handleUrlChange();
    },
    
    handleUrlChange() {
        const hash = window.location.hash.slice(1) || 'dashboard';
        const [route, ...paramParts] = hash.split('/');
        const params = {};
        
        paramParts.forEach((part, index) => {
            if (index % 2 === 0 && paramParts[index + 1]) {
                params[part] = paramParts[index + 1];
            }
        });
        
        this.loadRoute(route, params, false);
    },
    
    navigate(route, params = {}) {
        if (this.beforeNavigate) {
            const result = this.beforeNavigate(route, params);
            if (result === false) return;
        }
        
        const hash = this.buildHash(route, params);
        window.history.pushState({ route, params }, '', `#${hash}`);
        this.loadRoute(route, params, true);
    },
    
    buildHash(route, params) {
        const paramPairs = Object.entries(params)
            .filter(([_, value]) => value !== undefined && value !== null)
            .map(([key, value]) => `${key}/${value}`);
        
        return paramPairs.length > 0 
            ? `${route}/${paramPairs.join('/')}` 
            : route;
    },
    
    loadRoute(routeName, params = {}, pushState = true) {
        const route = this.routes[routeName];
        
        if (!route) {
            this.loadRoute('not-found', {}, false);
            return;
        }
        
        if (route.authRequired && !Store.isAuthenticated()) {
            this.navigate('login');
            return;
        }
        
        this.currentRoute = { name: routeName, params };
        
        document.querySelectorAll('.page').forEach(page => {
            page.classList.remove('active');
        });
        
        const pageElement = document.getElementById(`page-${routeName}`);
        if (pageElement) {
            pageElement.classList.add('active');
        }
        
        this.updateNavigation(routeName);
        
        if (route.onLoad) {
            route.onLoad(params);
        }
    },
    
    updateNavigation(currentRoute) {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
            const href = link.getAttribute('href');
            if (href === `#${currentRoute}` || href.startsWith(`#${currentRoute}/`)) {
                link.classList.add('active');
            }
        });
    },
    
    getParams() {
        return this.currentRoute ? this.currentRoute.params : {};
    },
    
    getRoute() {
        return this.currentRoute ? this.currentRoute.name : null;
    },
    
    setBeforeNavigate(callback) {
        this.beforeNavigate = callback;
    }
};
