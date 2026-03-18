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
        if (!data || !Array.isArray(data) || data.length === 0) {
            postbackInfoDiv.innerHTML = '<p>Нет данных для отображения.</p>';
            return;
        }

        let html = '';
        
        data.forEach((item, index) => {
            html += `<div class="postback-item" style="border: 1px solid #ccc; padding: 15px; margin-bottom: 20px; border-radius: 8px; background-color: #f9f9f9;">
                <div style="display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #eee; margin-bottom: 10px; padding-bottom: 5px;">
                    <strong style="color: #555;">Запрос #${data.length - index}</strong>
                    <span style="color: #888; font-size: 0.9em;">${item.received_at ? new Date(item.received_at).toLocaleString('ru-RU') : 'Время не указано'}</span>
                </div>
                
                <div class="data-section" style="margin-bottom: 10px;">
                    <h3 style="margin: 0 0 5px 0; font-size: 1.1em;">Path</h3>
                    <p style="margin: 0;"><code>${item.path || '/'}</code></p>
                </div>
            `;

            if (item.query_params && Object.keys(item.query_params).length > 0) {
                html += `
                    <div class="data-section">
                        <h3 style="margin: 0 0 5px 0; font-size: 1.1em;">Query Parameters</h3>
                        <ul style="margin: 0; padding-left: 20px;">
                `;
                for (const [key, value] of Object.entries(item.query_params)) {
                    html += `<li><strong>${key}:</strong> <code>${value}</code></li>`;
                }
                html += `
                        </ul>
                    </div>
                `;
            } else {
                html += `
                    <div class="data-section">
                        <h3 style="margin: 0 0 5px 0; font-size: 1.1em;">Query Parameters</h3>
                        <p style="margin: 0; color: #777;">Параметры отсутствуют.</p>
                    </div>
                `;
            }
            
            html += `</div>`;
        });

        postbackInfoDiv.innerHTML = html;
    }

    fetchPostbackData();
    setInterval(fetchPostbackData, 500);
});
