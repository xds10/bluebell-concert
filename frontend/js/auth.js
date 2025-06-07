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
            this.token = response.data.token;
            this.isAuthenticated = true;
            localStorage.setItem('token', this.token);
            // 兼容后续API统一读取loginResult
            localStorage.setItem('loginResult', JSON.stringify({ token: this.token }));
            this.updateUI();
            return true;
        } catch (error) {
            console.error('Login failed:', error);
            return false;
        }
    },

    async register(username, password) {
        try {
            await authAPI.register(username, password);
            return await this.login(username, password);
        } catch (error) {
            console.error('Registration failed:', error);
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
            loginBtn.style.display = 'none';
            registerBtn.style.display = 'none';
            logoutBtn.style.display = 'inline';
            loginForm.style.display = 'none';
            registerForm.style.display = 'none';
        } else {
            loginBtn.style.display = 'inline';
            registerBtn.style.display = 'inline';
            logoutBtn.style.display = 'none';
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

    // 注册表单提交
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

    // 登录按钮点击
    loginBtn.addEventListener('click', (e) => {
        e.preventDefault();
        document.getElementById('loginForm').style.display = 'block';
        document.getElementById('registerForm').style.display = 'none';
    });

    // 注册按钮点击
    registerBtn.addEventListener('click', (e) => {
        e.preventDefault();
        document.getElementById('registerForm').style.display = 'block';
        document.getElementById('loginForm').style.display = 'none';
    });

    // 退出按钮点击
    logoutBtn.addEventListener('click', (e) => {
        e.preventDefault();
        auth.logout();
        window.location.reload();
    });

    // 初始化认证状态
    auth.init();
}); 