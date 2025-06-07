// 页面导航管理
const navigation = {
    init() {
        const navLinks = document.querySelectorAll('.nav-links a[data-page]');
        navLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const page = e.target.dataset.page;
                this.navigateTo(page);
            });
        });
    },

    navigateTo(page) {
        // 更新导航栏激活状态
        document.querySelectorAll('.nav-links a').forEach(link => {
            link.classList.remove('active');
        });
        document.querySelector(`.nav-links a[data-page="${page}"]`).classList.add('active');

        // 显示对应页面内容
        switch (page) {
            case 'home':
                document.getElementById('concertList').style.display = 'block';
                document.getElementById('concertDetail').style.display = 'none';
                document.getElementById('loginForm').style.display = 'none';
                document.getElementById('registerForm').style.display = 'none';
                concertManager.loadConcerts();
                break;
            case 'concerts':
                document.getElementById('concertList').style.display = 'block';
                document.getElementById('concertDetail').style.display = 'none';
                document.getElementById('loginForm').style.display = 'none';
                document.getElementById('registerForm').style.display = 'none';
                concertManager.loadConcerts();
                break;
            case 'login':
                document.getElementById('loginForm').style.display = 'block';
                document.getElementById('registerForm').style.display = 'none';
                break;
            case 'register':
                document.getElementById('registerForm').style.display = 'block';
                document.getElementById('loginForm').style.display = 'none';
                break;
        }
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    navigation.init();
    auth.init();
    concertManager.loadConcerts();
}); 