document.addEventListener('DOMContentLoaded', () => {
    const postbackInfoDiv = document.getElementById('postback-info');
    const BACKEND_URL = 'https://testpostback-timofey3498.amvera.io';

    async function fetchPostbackData() {
        console.log('Fetching data from:', `${BACKEND_URL}/api/view`);
        try {
            const response = await fetch(`${BACKEND_URL}/api/view`); 
            console.log('Response status:', response.status);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const contentType = response.headers.get("content-type");
            console.log('Content-Type:', contentType);
            if (!contentType || !contentType.includes("application/json")) {
                throw new TypeError("Ожидался JSON ответ от сервера");
            }

            const data = await response.json();
            console.log('Data received:', data);
            
            renderData(data);

        } catch (error) {
            console.error('Error fetching postback data:', error);
            postbackInfoDiv.innerHTML = `<p style="color: red;">Error loading postback data: ${error.message}</p>
            <p><small>Проверьте консоль браузера (F12) для деталей.</small></p>`;
        }
    }

    function renderData(data) {
        if (!data || (!data.path && !data.query_params)) {
            postbackInfoDiv.innerHTML = '<p>Нет данных для отображения.</p>';
            return;
        }

        let html = `
            <div class="data-section">
                <h3>Path</h3>
                <p><code>${data.path || '/'}</code></p>
            </div>
        `;

        if (data.query_params && Object.keys(data.query_params).length > 0) {
            html += `
                <div class="data-section">
                    <h3>Query Parameters</h3>
                    <ul>
            `;
            for (const [key, value] of Object.entries(data.query_params)) {
                html += `<li><strong>${key}:</strong> <code>${value}</code></li>`;
            }
            html += `
                    </ul>
                </div>
            `;
        } else {
            html += `
                <div class="data-section">
                    <h3>Query Parameters</h3>
                    <p>Параметры отсутствуют.</p>
                </div>
            `;
        }

        postbackInfoDiv.innerHTML = html;
    }

    fetchPostbackData();
    // Обновляем данные каждые 5 секунд
    setInterval(fetchPostbackData, 5000);
});
