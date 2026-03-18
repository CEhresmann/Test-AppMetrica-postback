document.addEventListener('DOMContentLoaded', () => {
    const postbackInfoDiv = document.getElementById('postback-info');
    const BACKEND_URL = 'https://testpostback-timofey3498.amvera.io';

    async function fetchPostbackData() {
        try {
            const response = await fetch(`${BACKEND_URL}/view`); 
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const htmlContent = await response.text();
            
            const parser = new DOMParser();
            const doc = parser.parseFromString(htmlContent, 'text/html');
            const container = doc.querySelector('.container');
            
            if (container) {
                postbackInfoDiv.innerHTML = container.innerHTML;
            } else {
                postbackInfoDiv.innerHTML = htmlContent;
            }

        } catch (error) {
            console.error('Error fetching postback data:', error);
            postbackInfoDiv.innerHTML = `<p style="color: red;">Error loading postback data: ${error.message}</p>`;
        }
    }

    fetchPostbackData();
    // Обновляем данные каждые 5 секунд
    setInterval(fetchPostbackData, 5000);
});
