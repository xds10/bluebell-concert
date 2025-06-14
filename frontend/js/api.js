// API 基础配置
const API_BASE_URL = 'http://localhost:8081/api/v1';

// API 请求工具函数
async function apiRequest(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    const headers = {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };

    console.log(`Making request to: ${API_BASE_URL}${endpoint}`);

    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            ...options,
            headers,
            credentials: 'include'
        });

        console.log('Response status:', response.status);

        const data = await response.json();
        console.log('Response data:', data);
        return data;
    } catch (error) {
        console.error('API request failed:', error);
        throw error;
    }
}

// 用户认证相关 API
const authAPI = {
    login: async (username, password) => {
        return apiRequest('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
    },

    register: async (username, password) => {
        return apiRequest('/signup', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
    }
};

// 演唱会相关 API
const concertAPI = {
    getList: async () => {
        return apiRequest('/concert-list', {
            method: 'GET',
            headers: {
                'Cache-Control': 'no-cache'
            }
        });
    },

    getDetail: async (id) => {
        return apiRequest(`/concert/${id}`);
    },

    getSeatsInfo: async (id) => {
        return apiRequest(`/concert/${id}/seats`);
    },

    create: async (concertData) => {
        return apiRequest('/create_concert', {
            method: 'POST',
            body: JSON.stringify(concertData)
        });
    },

    buyTicket: async (concertId, seatSection) => {
        return apiRequest('/buy', {
            method: 'POST',
            body: JSON.stringify({
                concert_id: concertId,
                seat_idx: {
                    section: seatSection
                }
            })
        });
    },

    payOrder: async (orderData) => {
        return apiRequest('/pay', {
            method: 'POST',
            body: JSON.stringify(orderData)
        });
    }
};

// 订单相关 API
const orderAPI = {
    getOrderList: async (userId) => {
        return apiRequest('/order-list', {
            method: 'POST',
            body: JSON.stringify({ user_id: userId })
        });
    },

    getOrderDetail: async (orderId) => {
        return apiRequest(`/order/${orderId}`, {
            method: 'POST'
        });
    },

    payOrder: async (orderId, userId, concertId, seatId, price) => {
        return apiRequest('/pay', {
            method: 'POST',
            body: JSON.stringify({
                id: orderId,
                user_id: userId,
                concert_id: concertId,
                seat_id: seatId,
                price: price
            })
        });
    },
    
    cancelOrder: async (orderId) => {
        return apiRequest('/cancel-order', {
            method: 'POST',
            body: JSON.stringify({
                id: orderId
            })
        });
    }
}; 