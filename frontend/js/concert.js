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
                    <button onclick="concertManager.showConcertDetail(${concert.id})">查看详情</button>
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

    async showConcertDetail(id) {
        try {
            const response = await concertAPI.getDetail(id);
            if (response.code !== 1000 || !response.data) {
                throw new Error('Invalid response format');
            }
            const concert = response.data;
            
            const detailContainer = document.getElementById('concertDetail');
            const infoContainer = detailContainer.querySelector('.concert-info');
            
            infoContainer.innerHTML = `
                <h3>${concert.title}</h3>
                <p>开始时间：${new Date(concert.date).toLocaleString()}</p>
                <p>结束时间：${new Date(concert.enddate).toLocaleString()}</p>
                <p>开售时间：${new Date(concert.saletime).toLocaleString()}</p>
                <p>状态：${this.getStatusText(concert.status)}</p>
            `;

            document.getElementById('concertList').style.display = 'none';
            detailContainer.style.display = 'block';

            // 前端判断售票时间，只有在售票期间才显示选座，否则显示未开始售票
            const now = new Date();
            const saleTime = new Date(concert.saletime);
            const endDate = new Date(concert.enddate);
            const seatMap = document.querySelector('.seat-map');
            if (!isNaN(saleTime.getTime()) && !isNaN(endDate.getTime()) && now >= saleTime && now < endDate) {
                this.showSeatSelection(id);
            } else {
                seatMap.innerHTML = '<div style="padding: 32px; text-align: center; color: #888; font-size: 20px;">未开始售票</div>';
            }
        } catch (error) {
            console.error('Failed to load concert detail:', error);
            alert('加载演唱会详情失败，请稍后重试');
        }
    },

    async showSeatSelection(concertId) {
        const seatMap = document.querySelector('.seat-map');
        seatMap.innerHTML = '';

        // 只显示区域选择
        const sections = ['VIP区','A区', 'B区', 'C区'];
        sections.forEach(section => {
            const sectionDiv = document.createElement('div');
            sectionDiv.className = 'seat-section';
            sectionDiv.innerHTML = `<button class="section-btn">${section}区</button>`;
            sectionDiv.querySelector('button').onclick = () => this.selectSeat(concertId, section);
            seatMap.appendChild(sectionDiv);
        });
    },

    async selectSeat(concertId, section) {
        if (!auth.isAuthenticated) {
            alert('请先登录');
            return;
        }

        try {
            const response = await concertAPI.buyTicket(concertId, section);
            if (response.code === 1000) {
                // 选座成功后跳转到支付页面
                const orderData = {
                    id: response.data.id,
                    user_id: response.data.user_id,
                    concert_id: response.data.concert_id,
                    seat_id: response.data.seat_idx.id,
                    price: response.data.seat_idx.price,
                    status: 0,
                    create_time: new Date().toISOString()
                };
                // 将订单信息存储到localStorage，供支付页面读取
                localStorage.setItem('currentOrder', JSON.stringify(orderData));
                // 跳转到支付页面
                window.location.href = '/payment.html';
            } else {
                alert('选座失败：' + (response.msg || '未知错误'));
            }
        } catch (error) {
            console.error('Failed to select seat:', error);
            alert('选座失败，请稍后重试');
        }
    }
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
    concertManager.loadConcerts();
}); 