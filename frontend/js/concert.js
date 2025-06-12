// 演唱会列表管理
const concertManager = {
    async loadConcerts() {
        try {
            const response = await concertAPI.getList();
            console.log('API Response:', response);
            if (response.code === 1000 && response.data) {
                this.renderConcertList(response.data);
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
            const date = new Date(dateStr);
            return isNaN(date.getTime()) ? '无效时间' : date.toLocaleString();
        };
        const grid = document.getElementById('concertsGrid');
        if (!grid) {
            console.error('Concerts grid element not found');
            return;
        }
        grid.innerHTML = '';

        concerts.forEach(concert => {
            // 判断售票状态
            let statusText = this.getStatusText(concert.status);
            const now = new Date();
            const saleTime = new Date(concert.saletime);
            const endDate = new Date(concert.enddate);
            if (!isNaN(saleTime.getTime()) && !isNaN(endDate.getTime())) {
                if (now >= saleTime && now < endDate) {
                    statusText = '售票中';
                }
            }
            const card = document.createElement('div');
            card.className = 'concert-card';
            card.innerHTML = `
                <div class="concert-card-content">
                    <h3>${concert.title}</h3>
                    <p>开始时间：${formatDate(concert.date)}</p>
                    <p>结束时间：${formatDate(concert.enddate)}</p>
                    <p>开售时间：${formatDate(concert.saletime)}</p>
                    <p>状态：${statusText}</p>
                    <button onclick="window.location.href='concert_detail.html?concert_id=${concert.id}'">查看详情</button>
                </div>
            `;
            grid.appendChild(card);
        });
    },

    getStatusText(status) {
        switch (status) {
            case 0: return "未开始售票";
            case 1: return "售票中";
            case 2: return "已结束";
            default: return "未知状态";
        }
    },

    // showConcertDetail、showSeatSelection、selectSeat 方法可删除或保留但不再被调用
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    concertManager.loadConcerts();
}); 