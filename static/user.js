document.addEventListener("DOMContentLoaded", () => {
    // Проверка авторизации
    const authToken = localStorage.getItem('authToken');
    
    if (!authToken) {
        alert('Требуется авторизация');
        window.location.href = "/index.html";
        return;
    }
    
    // Загрузка клавиатур
    async function loadKeyboards() {
        try {
            const response = await fetch("http://localhost:9000/keyboard");
            
            if (response.status === 401) {
                alert("Ошибка авторизации");
                localStorage.clear();
                window.location.href = "/index.html";
                return;
            }
            
            const keyboards = await response.json();
            const tableBody = document.querySelector("#keyboardsTable tbody");
            tableBody.innerHTML = "";
            
            keyboards.forEach(keyboard => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${keyboard.id}</td>
                    <td>${keyboard.keycap_type}</td>
                    <td>${keyboard.base_type}</td>
                    <td>${keyboard.switch_type}</td>
                    <td>${keyboard.color}</td>
                `;
                tableBody.appendChild(row);
            });
        } catch (error) {
            console.error("Ошибка загрузки клавиатур:", error);
            alert("Ошибка загрузки клавиатур. Попробуйте позже.");
        }
    }

    // Загрузка клавиатур при открытии страницы
    loadKeyboards();
});