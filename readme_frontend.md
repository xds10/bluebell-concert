# 蓝铃铛演唱会售票系统前端文档

## 项目概述

蓝铃铛(Bluebell)演唱会售票系统前端是一个基于HTML、CSS和JavaScript开发的Web应用，为用户提供演唱会信息浏览、座位选择、票务购买和订单管理等功能。前端采用响应式设计，确保在不同设备上都能提供良好的用户体验。

## 技术栈

- **HTML5**: 页面结构
- **CSS3**: 样式和布局
- **JavaScript**: 交互逻辑
- **AJAX**: 与后端API通信
- **LocalStorage**: 本地数据存储

## 项目结构

```
frontend/
├── css/                 # 样式文件目录
│   ├── main.css         # 主样式文件
│   └── ...              # 其他样式文件
├── js/                  # JavaScript文件目录
│   ├── api.js           # API请求封装
│   ├── auth.js          # 认证相关功能
│   ├── concert.js       # 演唱会相关功能
│   ├── main.js          # 主要功能和初始化
│   └── ...              # 其他JS文件
├── index.html           # 首页
├── login.html           # 登录页面
├── register.html        # 注册页面
├── concert_detail.html  # 演唱会详情页面
├── orders.html          # 订单列表页面
├── order_detail.html    # 订单详情页面
├── payment.html         # 支付页面
└── captcha_test.html    # 验证码测试页面
```

## 页面说明

### 1. 首页 (index.html)

首页展示热门演唱会列表，用户可以浏览所有可用的演唱会信息，并通过点击进入演唱会详情页面。

**主要功能**:
- 展示演唱会列表
- 提供演唱会搜索功能
- 导航到其他页面的链接
- 用户登录/注册入口

### 2. 登录页面 (login.html)

用户可以通过此页面登录系统，登录成功后获取JWT令牌，用于后续的认证请求。

**主要功能**:
- 用户名和密码输入
- 表单验证
- 错误提示
- 登录请求处理
- 跳转到注册页面的链接

### 3. 注册页面 (register.html)

新用户可以通过此页面注册账号。

**主要功能**:
- 用户信息输入（用户名、密码、邮箱等）
- 表单验证
- 错误提示
- 注册请求处理
- 跳转到登录页面的链接

### 4. 演唱会详情页面 (concert_detail.html)

展示特定演唱会的详细信息，包括演出时间、地点、票价等，用户可以选择座位并购买票。

**主要功能**:
- 演唱会基本信息展示
- 座位图展示和选择
- 票价信息
- 购票功能
- 相关演唱会推荐

### 5. 订单列表页面 (orders.html)

用户可以查看自己的所有订单，包括待支付、已支付和已取消的订单。

**主要功能**:
- 订单列表展示
- 订单状态过滤
- 订单详情链接
- 支付和取消订单操作

### 6. 订单详情页面 (order_detail.html)

展示特定订单的详细信息，包括演唱会信息、座位信息、支付状态等。

**主要功能**:
- 订单基本信息展示
- 票务信息展示
- 支付和取消订单操作
- 生成电子票功能

### 7. 支付页面 (payment.html)

用户可以选择支付方式并完成订单支付。

**主要功能**:
- 订单信息确认
- 支付方式选择
- 支付流程处理
- 支付结果展示

## API接口调用

前端通过`api.js`中封装的函数与后端API进行交互，主要包括以下几类接口：

### 1. 用户认证

```javascript
// 用户注册
async function register(username, password, rePassword, email) {
  return await apiRequest('/api/v1/signup', 'POST', {
    username,
    password,
    re_password: rePassword,
    email
  });
}

// 用户登录
async function login(username, password) {
  return await apiRequest('/api/v1/login', 'POST', {
    username,
    password
  });
}
```

### 2. 演唱会信息

```javascript
// 获取演唱会列表
async function getConcertList() {
  return await apiRequest('/api/v1/concert-list', 'GET');
}

// 获取演唱会详情
async function getConcertDetail(concertId) {
  return await apiRequest(`/api/v1/concert/${concertId}`, 'GET');
}

// 获取演唱会座位信息
async function getConcertSeats(concertId) {
  return await apiRequest(`/api/v1/concert/${concertId}/seats`, 'GET');
}
```

### 3. 票务和订单

```javascript
// 购买票
async function buyTicket(concertId, seatIds) {
  return await apiRequest('/api/v1/buy', 'POST', {
    concert_id: concertId,
    seat_ids: seatIds
  });
}

// 支付订单
async function payOrder(orderId, paymentMethod) {
  return await apiRequest('/api/v1/pay', 'POST', {
    order_id: orderId,
    payment_method: paymentMethod
  });
}

// 获取订单列表
async function getOrderList(page, size) {
  return await apiRequest('/api/v1/order-list', 'POST', {
    page,
    size
  });
}

// 获取订单详情
async function getOrderDetail(orderId) {
  return await apiRequest(`/api/v1/order/${orderId}`, 'POST');
}

// 取消订单
async function cancelOrder(orderId) {
  return await apiRequest('/api/v1/cancel-order', 'POST', {
    order_id: orderId
  });
}
```

## 认证和授权

前端使用JWT（JSON Web Token）进行用户认证，token存储在LocalStorage中，并在每次API请求时通过Authorization头部发送。

```javascript
// 保存token到LocalStorage
function saveToken(token) {
  localStorage.setItem('token', token);
}

// 从LocalStorage获取token
function getToken() {
  return localStorage.getItem('token');
}

// 清除token（登出）
function clearToken() {
  localStorage.removeItem('token');
}

// 检查用户是否已登录
function isLoggedIn() {
  return !!getToken();
}

// 在API请求中添加认证头
function apiRequest(url, method, data = null) {
  const headers = {
    'Content-Type': 'application/json'
  };
  
  const token = getToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  
  // 发送请求...
}
```

## 座位选择功能

演唱会座位选择是系统的核心功能之一，实现了交互式座位图，用户可以直观地选择座位。

```javascript
// 初始化座位图
function initSeatMap(concertId) {
  // 获取座位数据
  getConcertSeats(concertId).then(response => {
    if (response.code === 1000) {
      const seats = response.data;
      renderSeatMap(seats);
    } else {
      showError('获取座位信息失败');
    }
  });
}

// 渲染座位图
function renderSeatMap(seats) {
  const seatMap = document.getElementById('seat-map');
  seatMap.innerHTML = '';
  
  // 计算座位图尺寸
  const rows = Math.max(...seats.map(seat => seat.row)) + 1;
  const cols = Math.max(...seats.map(seat => seat.column)) + 1;
  
  // 创建座位网格
  for (let r = 0; r < rows; r++) {
    const row = document.createElement('div');
    row.className = 'seat-row';
    
    for (let c = 0; c < cols; c++) {
      const seat = seats.find(s => s.row === r && s.column === c);
      if (seat) {
        const seatEl = document.createElement('div');
        seatEl.className = `seat status-${seat.status}`;
        seatEl.dataset.seatId = seat.id;
        
        if (seat.status === 0) { // 可选座位
          seatEl.addEventListener('click', () => toggleSeatSelection(seatEl, seat.id));
        }
        
        row.appendChild(seatEl);
      } else {
        // 空位置
        const spacer = document.createElement('div');
        spacer.className = 'seat-spacer';
        row.appendChild(spacer);
      }
    }
    
    seatMap.appendChild(row);
  }
}

// 切换座位选择状态
function toggleSeatSelection(seatEl, seatId) {
  if (seatEl.classList.contains('selected')) {
    seatEl.classList.remove('selected');
    selectedSeats = selectedSeats.filter(id => id !== seatId);
  } else {
    seatEl.classList.add('selected');
    selectedSeats.push(seatId);
  }
  
  // 更新选择信息和价格
  updateSelectionInfo();
}
```

## 高并发抢票处理

前端针对高并发抢票场景进行了特殊处理，包括：

1. **倒计时功能**：在演唱会开票前显示倒计时，准确同步开票时间。
2. **防重复提交**：在提交购票请求后禁用提交按钮，防止用户重复提交。
3. **状态实时更新**：定期轮询座位状态，确保用户看到的是最新的座位可用情况。
4. **错误重试机制**：对于因网络问题导致的请求失败，实现自动重试功能。
5. **用户友好的等待提示**：在处理购票请求时显示进度指示器，提升用户体验。

```javascript
// 开始抢票
async function startBuyingTicket(concertId, seatIds) {
  // 禁用提交按钮
  const buyButton = document.getElementById('buy-button');
  buyButton.disabled = true;
  
  // 显示加载指示器
  showLoading('正在处理您的购票请求...');
  
  try {
    // 尝试购买票
    const result = await buyTicket(concertId, seatIds);
    
    if (result.code === 1000) {
      // 购票成功，跳转到订单页面
      window.location.href = `/order_detail.html?id=${result.data.order_id}`;
    } else {
      // 购票失败，显示错误信息
      showError(result.message || '购票失败，请稍后重试');
      // 重新启用按钮
      buyButton.disabled = false;
    }
  } catch (error) {
    // 网络错误，自动重试
    retryCount++;
    if (retryCount <= MAX_RETRY) {
      showMessage(`网络请求失败，正在第${retryCount}次重试...`);
      setTimeout(() => startBuyingTicket(concertId, seatIds), 1000);
    } else {
      showError('网络连接不稳定，请稍后重试');
      buyButton.disabled = false;
    }
  } finally {
    hideLoading();
  }
}
```

## 响应式设计

前端采用响应式设计，确保在不同设备（桌面、平板、手机）上都能提供良好的用户体验。

```css
/* 响应式布局 */
@media (max-width: 768px) {
  .container {
    width: 100%;
    padding: 0 15px;
  }
  
  .seat {
    width: 20px;
    height: 20px;
  }
  
  .concert-card {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .seat {
    width: 15px;
    height: 15px;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 5px;
  }
  
  .form-group input {
    width: 100%;
  }
}
```

## 性能优化

前端实现了多项性能优化措施，确保在高并发场景下也能提供流畅的用户体验：

1. **资源压缩**：压缩CSS和JavaScript文件，减少加载时间。
2. **延迟加载**：非关键资源采用延迟加载策略。
3. **缓存策略**：利用浏览器缓存减少重复请求。
4. **批量请求**：合并多个API请求，减少网络开销。
5. **DOM操作优化**：减少直接DOM操作，使用DocumentFragment进行批量更新。

```javascript
// 使用DocumentFragment优化DOM操作
function renderConcertList(concerts) {
  const container = document.getElementById('concert-list');
  const fragment = document.createDocumentFragment();
  
  concerts.forEach(concert => {
    const card = createConcertCard(concert);
    fragment.appendChild(card);
  });
  
  // 一次性更新DOM
  container.innerHTML = '';
  container.appendChild(fragment);
}
```

## 错误处理

前端实现了全局错误处理机制，捕获并展示友好的错误信息：

```javascript
// 全局错误处理
window.addEventListener('error', function(event) {
  console.error('全局错误:', event.error);
  showError('发生了一个错误，请刷新页面重试');
  return false;
});

// API错误处理
async function apiRequest(url, method, data = null) {
  try {
    // 发送请求...
    const response = await fetch(url, options);
    const result = await response.json();
    
    if (!response.ok) {
      // HTTP错误
      handleApiError(result, response.status);
      return result;
    }
    
    return result;
  } catch (error) {
    // 网络错误
    console.error('API请求错误:', error);
    return {
      code: 5000,
      message: '网络连接错误，请检查您的网络连接'
    };
  }
}

// 处理API错误
function handleApiError(error, status) {
  if (status === 401) {
    // 未授权，清除token并跳转到登录页面
    clearToken();
    showError('登录已过期，请重新登录');
    setTimeout(() => {
      window.location.href = '/login.html?redirect=' + encodeURIComponent(window.location.href);
    }, 2000);
  } else if (status === 403) {
    showError('您没有权限执行此操作');
  } else if (status === 404) {
    showError('请求的资源不存在');
  } else {
    showError(error.message || '服务器错误，请稍后重试');
  }
}
```

## 本地存储

前端使用LocalStorage存储用户信息、认证令牌和一些配置数据：

```javascript
// 用户信息存储
const userStorage = {
  // 保存用户信息
  saveUser(user) {
    localStorage.setItem('user', JSON.stringify(user));
  },
  
  // 获取用户信息
  getUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },
  
  // 清除用户信息
  clearUser() {
    localStorage.removeItem('user');
  }
};

// 配置存储
const configStorage = {
  // 保存配置
  saveConfig(key, value) {
    localStorage.setItem(`config_${key}`, JSON.stringify(value));
  },
  
  // 获取配置
  getConfig(key, defaultValue = null) {
    const valueStr = localStorage.getItem(`config_${key}`);
    return valueStr ? JSON.parse(valueStr) : defaultValue;
  }
};
```

## 安全措施

前端实现了多项安全措施，保护用户数据和防止常见的安全漏洞：

1. **输入验证**：对所有用户输入进行验证，防止XSS攻击。
2. **CSRF防护**：在表单提交中加入CSRF令牌。
3. **敏感信息保护**：不在前端存储敏感信息，如明文密码。
4. **安全的API通信**：使用HTTPS进行API通信。
5. **令牌过期处理**：检测JWT令牌过期并引导用户重新登录。

```javascript
// 输入验证
function validateInput(input, type) {
  const value = input.value.trim();
  
  switch (type) {
    case 'username':
      return /^[a-zA-Z0-9_]{4,16}$/.test(value);
    case 'password':
      return /^.{6,20}$/.test(value);
    case 'email':
      return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
    default:
      return true;
  }
}

// 防止XSS
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}
```

## 浏览器兼容性

前端支持所有主流现代浏览器，包括：

- Chrome 60+
- Firefox 60+
- Safari 11+
- Edge 16+

对于旧版浏览器，提供了优雅降级方案，确保基本功能可用。

## 未来计划

前端未来的开发计划包括：

1. **迁移到现代前端框架**：考虑使用Vue.js或React重构前端，提升开发效率和用户体验。
2. **移动应用开发**：开发原生移动应用或PWA，提供更好的移动端体验。
3. **实时通知**：集成WebSocket，实现实时座位状态更新和通知功能。
4. **多语言支持**：添加国际化支持，使系统可以支持多种语言。
5. **无障碍优化**：提升系统的无障碍性，使其符合WCAG标准。

## 开发指南

### 环境设置

1. 克隆项目代码
2. 使用现代浏览器打开HTML文件即可开始开发
3. 建议使用VS Code等编辑器，并安装Live Server插件进行本地开发

### 编码规范

- HTML: 使用语义化标签，确保无障碍性
- CSS: 使用BEM命名规范
- JavaScript: 遵循ES6+标准，使用异步/await处理异步操作

### 调试技巧

1. 使用浏览器开发者工具进行调试
2. 检查Network面板监控API请求
3. 使用Console进行日志输出和错误跟踪

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开Pull Request


