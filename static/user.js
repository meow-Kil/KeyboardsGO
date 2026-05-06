document.addEventListener("DOMContentLoaded", () => {

    const authToken = localStorage.getItem('authToken');
    
    if (!authToken) {
        alert('Требуется авторизация');
        window.location.href = "/index.html";
        return;
    }
    

    async function loadKeyboards() {
        try {
            const response = await fetch("http://localhost:1000/keyboard");
            
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


    async function loadKeycapTypes() {
        try {
            const response = await fetch("http://localhost:1000/keycap_types");
            if (response.status === 401) {
                localStorage.clear();
                window.location.href = "/index.html";
                return;
            }
            const types = await response.json();
            const tbody = document.querySelector("#keycapTypesTable tbody");
            tbody.innerHTML = "";
            types.forEach(kt => {
                const row = tbody.insertRow();
                row.insertCell(0).innerText = kt.id;
                row.insertCell(1).innerText = kt.name;
            });
        } catch (err) {
            console.error("Ошибка загрузки типов:", err);
        }
    }

    loadKeyboards();
    loadKeycapTypes();
});