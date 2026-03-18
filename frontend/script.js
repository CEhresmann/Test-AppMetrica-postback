document.addEventListener('DOMContentLoaded', () => {
    const postbackInfoDiv = document.getElementById('postback-info');

    async function fetchPostbackData() {
        try {
            const response = await fetch('/view'); 
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const htmlContent = await response.text();
            
            postbackInfoDiv.innerHTML = htmlContent;

        } catch (error) {
            console.error('Error fetching postback data:', error);
            postbackInfoDiv.innerHTML = `<p style="color: red;">Error loading postback data: ${error.message}</p>`;
        }
    }

    fetchPostbackData();
});
