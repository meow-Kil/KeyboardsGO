document.addEventListener("DOMContentLoaded", () => {
    // Проверка авторизации
    const authToken = localStorage.getItem('authToken');
    const isAdmin = localStorage.getItem('isAdmin');
    
    if (!authToken || isAdmin !== 'true') {
        alert('Требуется авторизация администратора');
        window.location.href = "/index.html";
        return;
    }
    
    const addKeyboardForm = document.getElementById("addKeyboardForm");
    const editKeyboardForm = document.getElementById("editKeyboardForm");
    const editOverlay = document.getElementById("editOverlay");

    // Получение токена авторизации
    function getAuthToken() {
        return localStorage.getItem('authToken');
    }

    // Загрузка клавиатур
    async function loadKeyboards() {
        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:9000/keyboard", {
                headers: {
                    "Authorization": token ? `Bearer ${token}` : ""
                }
            });
            
            if (response.status === 401) {
                alert("Сессия истекла. Пожалуйста, войдите снова.");
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
                    <td class="actions">
                        <button onclick="editKeyboard(${keyboard.id})">Редактировать</button>
                        <button class="delete" onclick="deleteKeyboard(${keyboard.id})">Удалить</button>
                    </td>
                `;
                tableBody.appendChild(row);
            });
        } catch (error) {
            console.error("Ошибка загрузки клавиатур:", error);
            alert("Ошибка загрузки клавиатур. Попробуйте позже.");
        }
    }

    // Добавление клавиатуры
    addKeyboardForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const keycapType = document.getElementById("keycapType").value.trim();
        const baseType = document.getElementById("baseType").value.trim();
        const switchType = document.getElementById("switchType").value.trim();
        const color = document.getElementById("color").value.trim();

        if (!keycapType || !baseType || !switchType || !color) {
            alert("Пожалуйста, заполните все поля");
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:9000/keyboard", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": token ? `Bearer ${token}` : ""
                },
                body: JSON.stringify({
                    keycap_type: keycapType,
                    base_type: baseType,
                    switch_type: switchType,
                    color: color,
                }),
            });

            const responseText = await response.text();
            
            if (!response.ok) {
                try {
                    const errorData = JSON.parse(responseText);
                    alert(`Ошибка: ${errorData.error || "Неизвестная ошибка"}`);
                } catch {
                    alert(`Ошибка: ${responseText || "Неизвестная ошибка"}`);
                }
                return;
            }

            const result = JSON.parse(responseText);
            alert("Клавиатура успешно добавлена!");
            loadKeyboards();
            addKeyboardForm.reset();
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    });

    // Редактирование клавиатуры
    window.editKeyboard = async function(id) {
        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:9000/keyboard/${id}`, {
                headers: {
                    "Authorization": token ? `Bearer ${token}` : ""
                }
            });
            
            if (response.status === 401) {
                alert("Сессия истекла. Пожалуйста, войдите снова.");
                localStorage.clear();
                window.location.href = "/index.html";
                return;
            }
            
            const keyboard = await response.json();

            document.getElementById("editId").value = keyboard.id;
            document.getElementById("editKeycapType").value = keyboard.keycap_type;
            document.getElementById("editBaseType").value = keyboard.base_type;
            document.getElementById("editSwitchType").value = keyboard.switch_type;
            document.getElementById("editColor").value = keyboard.color;

            editKeyboardForm.style.display = "block";
            editOverlay.style.display = "block";
        } catch (error) {
            console.error("Ошибка:", error);
            alert("Ошибка загрузки данных клавиатуры");
        }
    };

    // Сохранение изменений
    editKeyboardForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const id = document.getElementById("editId").value;
        const keycapType = document.getElementById("editKeycapType").value.trim();
        const baseType = document.getElementById("editBaseType").value.trim();
        const switchType = document.getElementById("editSwitchType").value.trim();
        const color = document.getElementById("editColor").value.trim();

        if (!keycapType || !baseType || !switchType || !color) {
            alert("Пожалуйста, заполните все поля");
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:9000/keyboard/${id}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": token ? `Bearer ${token}` : ""
                },
                body: JSON.stringify({
                    keycap_type: keycapType,
                    base_type: baseType,
                    switch_type: switchType,
                    color: color,
                }),
            });

            const responseText = await response.text();
            
            if (!response.ok) {
                try {
                    const errorData = JSON.parse(responseText);
                    alert(`Ошибка: ${errorData.error || "Неизвестная ошибка"}`);
                } catch {
                    alert(`Ошибка: ${responseText || "Неизвестная ошибка"}`);
                }
                return;
            }

            const result = JSON.parse(responseText);
            alert("Клавиатура успешно обновлена!");
            loadKeyboards();
            cancelEdit();
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    });

    // Отмена редактирования
    window.cancelEdit = function() {
        editKeyboardForm.style.display = "none";
        editOverlay.style.display = "none";
        editKeyboardForm.reset();
    };

    // Удаление клавиатуры
    window.deleteKeyboard = async function(id) {
        if (!confirm("Вы уверены, что хотите удалить эту клавиатуру?")) {
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:9000/keyboard/${id}`, {
                method: "DELETE",
                headers: {
                    "Authorization": token ? `Bearer ${token}` : ""
                }
            });

            const responseText = await response.text();
            
            if (!response.ok) {
                try {
                    const errorData = JSON.parse(responseText);
                    alert(`Ошибка: ${errorData.error || "Неизвестная ошибка"}`);
                } catch {
                    alert(`Ошибка: ${responseText || "Неизвестная ошибка"}`);
                }
                return;
            }

            const result = JSON.parse(responseText);
            alert("Клавиатура успешно удалена!");
            loadKeyboards();
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    };

    // Загрузка клавиатур при открытии страницы
    loadKeyboards();
});