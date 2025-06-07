<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <title>演唱会列表</title>
  <style>
    body {
      font-family: sans-serif;
      margin: 2em;
    }
    h1 {
      color: #333;
    }
    .concert {
      border: 1px solid #ccc;
      border-radius: 8px;
      padding: 1em;
      margin-bottom: 1em;
    }
    .concert:hover {
      background-color: #f9f9f9;
    }
    .details {
      margin-top: 2em;
      padding: 1em;
      border: 1px solid #aaa;
      background-color: #f0f0f0;
      border-radius: 8px;
    }
  </style>
</head>
<body>

  <h1>演唱会列表</h1>
  <div id="concert-list"></div>

  <div id="concert-details" class="details" style="display: none;">
    <h2>演唱会详情</h2>
    <p id="detail-title"></p>
    <p id="detail-date"></p>
    <p id="detail-enddate"></p>
    <p id="detail-saletime"></p>
    <p id="detail-status"></p>
  </div>

  <script>
    async function loadConcertList() {
      const response = await fetch('http://localhost:8081/api/v1/concert-list');
      const json = await response.json();
      const list = document.getElementById('concert-list');

      json.data.forEach(concert => {
        const div = document.createElement('div');
        div.className = 'concert';
        div.innerHTML = `
          <strong>${concert.title}</strong><br>
          演出时间：${new Date(concert.date).toLocaleString()}<br>
          <button onclick="loadConcertDetails(${concert.id})">查看详情</button>
        `;
        list.appendChild(div);
      });
    }

    async function loadConcertDetails(id) {
      const response = await fetch(`http://localhost:8081/api/v1/concert/${id}`);
      const json = await response.json();
      const data = json.data;

      document.getElementById('detail-title').textContent = `标题：${data.title}`;
      document.getElementById('detail-date').textContent = `开始时间：${new Date(data.date).toLocaleString()}`;
      document.getElementById('detail-enddate').textContent = `结束时间：${new Date(data.enddate).toLocaleString()}`;
      document.getElementById('detail-saletime').textContent = `开售时间：${new Date(data.saletime).toLocaleString()}`;
      document.getElementById('detail-status').textContent = `状态：${getStatusText(data.status)}`;
      document.getElementById('concert-details').style.display = 'block';
    }

    function getStatusText(status) {
      switch (status) {
        case 0: return "未开始售票";
        case 1: return "售票中";
        case 2: return "已结束";
        default: return "未知状态";
      }
    }

    loadConcertList();
  </script>
</body>
</html>
