// 演唱会列表管理
const concertManager = {
    // 存储所有演唱会数据
    allConcerts: [],
    
    async loadConcerts() {
        try {
            const response = await concertAPI.getList();
            console.log('API Response:', response);
            if (response.code === 1000 && response.data) {
                // 存储所有演唱会数据
                this.allConcerts = response.data;
                this.renderConcertList(this.allConcerts);
                
                // 初始化搜索功能
                this.initSearchFilter();
            } else {
                throw new Error(response.msg || '加载失败');
            }
        } catch (error) {
            console.error('加载演唱会列表失败:', error);
            alert('加载演唱会列表失败，请稍后重试');
        }
    },

    renderConcertList(concerts) {
        const formatDate = (dateStr) => {
            // 更严格的日期格式化
            try {
                const date = new Date(dateStr);
                return isNaN(date.getTime()) ? '无效时间' : date.toLocaleString('zh-CN');
            } catch (e) {
                return '无效时间';
            }
        };
        
        const grid = document.getElementById('concertsGrid');
        if (!grid) {
            console.error('Concerts grid element not found');
            return;
        }
        grid.innerHTML = '';
        
        if (concerts.length === 0) {
            grid.innerHTML = '<div class="no-results">没有找到匹配的演唱会</div>';
            return;
        }

        concerts.forEach(concert => {
            // 判断售票状态
            const statusInfo = this.calculateStatus(concert);
            
            const card = document.createElement('div');
            card.className = 'concert-card';
            card.innerHTML = `
                <div class="concert-card-content">
                    <h3>${concert.title}</h3>
                    <p>开始时间：${formatDate(concert.date)}</p>
                    <p>结束时间：${formatDate(concert.enddate)}</p>
                    <p>开售时间：${formatDate(concert.saletime)}</p>
                    <p>状态：${statusInfo.text}</p>
                    <button class="detail-btn" data-concert-id="${concert.id}">查看详情</button>
                </div>
            `;
            grid.appendChild(card);
        });

        // 使用事件委托处理按钮点击
        this.bindDetailButtons();
    },

    // 绑定详情按钮事件
    bindDetailButtons() {
        const grid = document.getElementById('concertsGrid');
        if (!grid) return;

        // 移除旧的事件监听器
        grid.removeEventListener('click', this.handleDetailClick);
        // 添加新的事件监听器
        grid.addEventListener('click', this.handleDetailClick);
    },

    // 处理详情按钮点击
    handleDetailClick(event) {
        if (event.target.classList.contains('detail-btn')) {
            const concertId = event.target.dataset.concertId;
            window.location.href = `concert_detail.html?concert_id=${concertId}`;
        }
    },
    
    // 计算演唱会实际状态
    calculateStatus(concert) {
        const now = new Date();
        const saleTime = new Date(concert.saletime);
        const startDate = new Date(concert.date);
        const endDate = new Date(concert.enddate);
        
        let statusText = "未知状态";
        let statusCode = -1;
        
        if (!isNaN(saleTime.getTime()) && !isNaN(startDate.getTime()) && !isNaN(endDate.getTime())) {
            if (now < saleTime) {
                statusText = "未开始售票";
                statusCode = 0;
            } else if (now >= saleTime && now < startDate) {
                statusText = "售票中";
                statusCode = 1;
            } else if (now >= startDate && now < endDate) {
                statusText = "演出中";
                statusCode = 3;
            } else {
                statusText = "已结束";
                statusCode = 2;
            }
        } else {
            // 如果日期无效，回退到使用后端状态
            statusText = this.getStatusText(concert.status);
            statusCode = concert.status;
        }
        
        return { text: statusText, code: statusCode };
    },

    getStatusText(status) {
        switch (status) {
            case 0: return "未开始售票";
            case 1: return "售票中";
            case 2: return "已结束";
            case 3: return "演出中";
            default: return "未知状态";
        }
    },
    
    // 创建搜索过滤UI - 修复重复创建问题
    createSearchFilterUI() {
        const concertList = document.getElementById('concertList');
        if (!concertList) return;
        
        // 检查是否已经存在过滤容器，如果存在则先移除
        const existingFilter = document.querySelector('.filter-container');
        if (existingFilter) {
            existingFilter.remove();
        }
        
        // 创建过滤容器
        const filterContainer = document.createElement('div');
        filterContainer.className = 'filter-container';
        filterContainer.innerHTML = `
            <div class="search-filter">
                <input type="text" id="concertSearchInput" placeholder="搜索演唱会名称..." class="search-input">
                <select id="statusFilter" class="status-filter">
                    <option value="all">所有状态</option>
                    <option value="0">未开始售票</option>
                    <option value="1">售票中</option>
                    <option value="3">演出中</option>
                    <option value="2">已结束</option>
                </select>
            </div>
        `;
        
        // 插入到演唱会列表标题后面
        const title = concertList.querySelector('h2');
        if (title) {
            title.after(filterContainer);
            
            // 只添加一次样式
            if (!document.getElementById('concert-filter-styles')) {
                const style = document.createElement('style');
                style.id = 'concert-filter-styles';
                style.textContent = `
                    .filter-container {
                        margin-bottom: 20px;
                    }
                    .search-filter {
                        display: flex;
                        gap: 10px;
                        margin-bottom: 20px;
                    }
                    .search-input {
                        flex: 1;
                        padding: 10px;
                        border: 1px solid #ddd;
                        border-radius: 4px;
                        font-size: 16px;
                    }
                    .status-filter {
                        padding: 10px;
                        border: 1px solid #ddd;
                        border-radius: 4px;
                        font-size: 16px;
                        min-width: 120px;
                    }
                    .no-results {
                        text-align: center;
                        padding: 40px;
                        color: #666;
                        font-size: 18px;
                    }
                `;
                document.head.appendChild(style);
            }
        }
    },
    
    // 初始化搜索过滤功能
    initSearchFilter() {
        // 创建搜索过滤区域
        this.createSearchFilterUI();
        
        // 移除旧的事件监听器
        this.removeSearchEventListeners();
        
        // 绑定搜索事件
        const searchInput = document.getElementById('concertSearchInput');
        if (searchInput) {
            this.searchInputHandler = () => this.filterConcerts();
            searchInput.addEventListener('input', this.searchInputHandler);
        }
        
        // 绑定状态过滤事件
        const statusFilter = document.getElementById('statusFilter');
        if (statusFilter) {
            this.statusFilterHandler = () => this.filterConcerts();
            statusFilter.addEventListener('change', this.statusFilterHandler);
        }
    },

    // 移除搜索事件监听器
    removeSearchEventListeners() {
        const searchInput = document.getElementById('concertSearchInput');
        const statusFilter = document.getElementById('statusFilter');
        
        if (searchInput && this.searchInputHandler) {
            searchInput.removeEventListener('input', this.searchInputHandler);
        }
        
        if (statusFilter && this.statusFilterHandler) {
            statusFilter.removeEventListener('change', this.statusFilterHandler);
        }
    },
    
    // 过滤演唱会
    filterConcerts() {
        const searchInput = document.getElementById('concertSearchInput');
        const statusFilter = document.getElementById('statusFilter');
        
        if (!searchInput || !statusFilter) return;
        
        const searchTerm = searchInput.value.trim().toLowerCase();
        const statusValue = statusFilter.value;
        
        // 过滤演唱会列表
        const filteredConcerts = this.allConcerts.filter(concert => {
            // 搜索词过滤
            const matchesSearch = concert.title.toLowerCase().includes(searchTerm);
            
            // 状态过滤
            let matchesStatus = true;
            if (statusValue !== 'all') {
                const status = this.calculateStatus(concert);
                matchesStatus = status.code.toString() === statusValue;
            }
            
            return matchesSearch && matchesStatus;
        });
        
        // 渲染过滤后的列表
        this.renderConcertList(filteredConcerts);
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    concertManager.loadConcerts();
});