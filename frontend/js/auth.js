// 认证状态管理
const auth = {
    isAuthenticated: false,
    token: null,
    user: null,

    init() {
        const token = localStorage.getItem('token');
        if (token) {
            this.token = token;
            this.isAuthenticated = true;
            // 兼容后续API统一读取loginResult
            if (!localStorage.getItem('loginResult')) {
                localStorage.setItem('loginResult', JSON.stringify({ token: this.token }));
            }
            this.updateUI();
        }
    },

    async login(username, password) {
        try {
            const response = await authAPI.login(username, password);
            // 兼容后端返回结构
            if (response.code === 1000 && response.data && response.data.token) {
                this.token = response.data.token;
                this.isAuthenticated = true;
                localStorage.setItem('token', this.token);
                localStorage.setItem('loginResult', JSON.stringify({
                    token: this.token,
                    user_id: response.data.user_id,
                    user_name: response.data.user_name
                }));
                this.updateUI && this.updateUI();
                return true;
            } else {
                return false;
            }
        } catch (error) {
            console.error('Login failed:', error);
            return false;
        }
    },



    async register(username, password, rePassword) {
        try {
            console.log('开始注册，用户名:', username);
            
            // 调用注册API，传递用户名、密码和重复密码
            const response = await authAPI.register(username, password, rePassword);
            console.log('注册响应:', response);
            
            if (response.code === 1000) {
                console.log('注册成功，准备自动登录');
                return await this.login(username, password);
            } else {
                // 注册失败，返回错误信息
                console.error('注册失败，错误码:', response.code, '消息:', response.msg);
                alert('注册失败: ' + (response.msg || '未知错误'));
                return false;
            }
        } catch (error) {
            console.error('注册过程中出错:', error);
            return false;
        }
    },

    logout() {
        this.token = null;
        this.isAuthenticated = false;
        this.user = null;
        localStorage.removeItem('token');
        this.updateUI();
    },

    updateUI() {
        const loginBtn = document.getElementById('loginBtn');
        const registerBtn = document.getElementById('registerBtn');
        const logoutBtn = document.getElementById('logoutBtn');
        const loginForm = document.getElementById('loginForm');
        const registerForm = document.getElementById('registerForm');

        if (this.isAuthenticated) {
            if (loginBtn) loginBtn.style.display = 'none';
            if (registerBtn) registerBtn.style.display = 'none';
            if (logoutBtn) logoutBtn.style.display = 'inline';
            if (loginForm) loginForm.style.display = 'none';
            if (registerForm) registerForm.style.display = 'none';
        } else {
            if (loginBtn) loginBtn.style.display = 'inline';
            if (registerBtn) registerBtn.style.display = 'inline';
            if (logoutBtn) logoutBtn.style.display = 'none';
        }
    }
};

// 表单处理
document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login');
    const registerForm = document.getElementById('register');
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const logoutBtn = document.getElementById('logoutBtn');

    // 登录表单提交
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('loginUsername').value;
            const password = document.getElementById('loginPassword').value;
            
            if (await auth.login(username, password)) {
                document.getElementById('loginForm').style.display = 'none';
                window.location.reload();
            } else {
                alert('登录失败，请检查用户名和密码');
            }
        });
    }

    // 注册表单提交
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const username = document.getElementById('registerUsername').value;
            const password = document.getElementById('registerPassword').value;
            
            if (await auth.register(username, password)) {
                document.getElementById('registerForm').style.display = 'none';
                window.location.reload();
            } else {
                alert('注册失败，请稍后重试');
            }
        });
    }

    // 登录按钮点击
    if (loginBtn) {
        loginBtn.addEventListener('click', (e) => {
            // 不阻止默认行为，允许链接跳转到login.html
            // 但需要判断是否在主页，如果在主页才跳转，否则直接显示嵌入式表单
            if (window.location.pathname.includes('index.html') || window.location.pathname === '/' || window.location.pathname.endsWith('/')) {
                // 允许跳转，不做任何处理
            } else {
                // 在其他页面保持原有逻辑，显示嵌入式表单
                e.preventDefault();
                const loginForm = document.getElementById('loginForm');
                const registerForm = document.getElementById('registerForm');
                if (loginForm) loginForm.style.display = 'block';
                if (registerForm) registerForm.style.display = 'none';
            }
        });
    }

    // 注册按钮点击
    if (registerBtn) {
        registerBtn.addEventListener('click', (e) => {
            // 不阻止默认行为，允许链接跳转到register.html
            // 但需要判断是否在主页，如果在主页才跳转，否则直接显示嵌入式表单
            if (window.location.pathname.includes('index.html') || window.location.pathname === '/' || window.location.pathname.endsWith('/')) {
                // 允许跳转，不做任何处理
            } else {
                // 在其他页面保持原有逻辑，显示嵌入式表单
                e.preventDefault();
                const loginForm = document.getElementById('loginForm');
                const registerForm = document.getElementById('registerForm');
                if (registerForm) registerForm.style.display = 'block';
                if (loginForm) loginForm.style.display = 'none';
            }
        });
    }

    // 退出按钮点击
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            auth.logout();
            window.location.reload();
        });
    }

    // 初始化认证状态
    auth.init();
}); 