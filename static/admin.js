document.addEventListener("DOMContentLoaded", () => {
   
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

  
    function getAuthToken() {
        return localStorage.getItem('authToken');
    }


    function populateKeycapTypeSelects(types) {
        const addSelect = document.getElementById("keycapType");
        const editSelect = document.getElementById("editKeycapType");
        
  
        addSelect.innerHTML = '<option value="">Выберите тип колпачков</option>';
        editSelect.innerHTML = '<option value="">Выберите тип колпачков</option>';
        
        types.forEach(kt => {
            const option = document.createElement("option");
            option.value = kt.name;  
            option.textContent = kt.name;
            addSelect.appendChild(option);
            
            const option2 = document.createElement("option");
            option2.value = kt.name;
            option2.textContent = kt.name;
            editSelect.appendChild(option2);
        });
    }


    async function refreshKeycapTypeSelects() {
        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:1000/keycap_types", {
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
            });
            if (response.ok) {
                const types = await response.json();
                populateKeycapTypeSelects(types);
            }
        } catch (err) {
            console.error("Ошибка загрузки типов для селектов:", err);
        }
    }

    async function loadKeyboards() {
        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:1000/keyboard", {
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
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

    addKeyboardForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const keycapType = document.getElementById("keycapType").value;
        const baseType = document.getElementById("baseType").value.trim();
        const switchType = document.getElementById("switchType").value.trim();
        const color = document.getElementById("color").value.trim();

        if (!keycapType || !baseType || !switchType || !color) {
            alert("Пожалуйста, заполните все поля и выберите тип колпачков");
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:1000/keyboard", {
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

            alert("Клавиатура успешно добавлена!");
            loadKeyboards();
            addKeyboardForm.reset();
       
            document.getElementById("keycapType").selectedIndex = 0;
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    });

    window.editKeyboard = async function(id) {
        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:1000/keyboard/${id}`, {
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
            });
            
            if (response.status === 401) {
                alert("Сессия истекла. Пожалуйста, войдите снова.");
                localStorage.clear();
                window.location.href = "/index.html";
                return;
            }
            
            const keyboard = await response.json();

            document.getElementById("editId").value = keyboard.id;
    
            const editSelect = document.getElementById("editKeycapType");
            for (let i = 0; i < editSelect.options.length; i++) {
                if (editSelect.options[i].value === keyboard.keycap_type) {
                    editSelect.selectedIndex = i;
                    break;
                }
            }
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

    editKeyboardForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const id = document.getElementById("editId").value;
        const keycapType = document.getElementById("editKeycapType").value;
        const baseType = document.getElementById("editBaseType").value.trim();
        const switchType = document.getElementById("editSwitchType").value.trim();
        const color = document.getElementById("editColor").value.trim();

        if (!keycapType || !baseType || !switchType || !color) {
            alert("Пожалуйста, заполните все поля и выберите тип колпачков");
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:1000/keyboard/${id}`, {
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

            alert("Клавиатура успешно обновлена!");
            loadKeyboards();
            cancelEdit();
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    });

    window.cancelEdit = function() {
        editKeyboardForm.style.display = "none";
        editOverlay.style.display = "none";
        editKeyboardForm.reset();
     
        document.getElementById("editKeycapType").selectedIndex = 0;
    };

    window.deleteKeyboard = async function(id) {
        if (!confirm("Вы уверены, что хотите удалить эту клавиатуру?")) {
            return;
        }

        try {
            const token = getAuthToken();
            const response = await fetch(`http://localhost:1000/keyboard/${id}`, {
                method: "DELETE",
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
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

            alert("Клавиатура успешно удалена!");
            loadKeyboards();
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        }
    };


    async function loadKeycapTypes() {
        try {
            const token = getAuthToken();
            const response = await fetch("http://localhost:1000/keycap_types", {
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
            });
            if (response.status === 401) {
                alert("Сессия истекла");
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
                const actionsCell = row.insertCell(2);
                actionsCell.innerHTML = `
                    <button onclick="editKeycapType(${kt.id})">Ред.</button>
                    <button class="delete" onclick="deleteKeycapType(${kt.id})">Удалить</button>
                `;
            });

            populateKeycapTypeSelects(types);
        } catch (err) {
            console.error("Ошибка загрузки типов:", err);
        }
    }

    document.getElementById("addKeycapTypeForm")?.addEventListener("submit", async (e) => {
        e.preventDefault();
        const name = document.getElementById("keycapTypeName").value.trim();
        if (!name) return alert("Введите название");
        const token = getAuthToken();
        try {
            const res = await fetch("http://localhost:1000/keycap_types", {
                method: "POST",
                headers: { 
                    "Content-Type": "application/json", 
                    "Authorization": token ? `Bearer ${token}` : "" 
                },
                body: JSON.stringify({ name })
            });
            const data = await res.json();
            if (data.success) {
                alert("Тип добавлен");
                loadKeycapTypes();  
                document.getElementById("addKeycapTypeForm").reset();
            } else {
                alert("Ошибка: " + (data.error || "Неизвестная"));
            }
        } catch (err) {
            alert("Ошибка сети");
        }
    });

    window.editKeycapType = async function(id) {
        const token = getAuthToken();
        try {
            const res = await fetch(`http://localhost:1000/keycap_types/${id}`, {
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
            });
            const kt = await res.json();
            document.getElementById("editKeycapTypeId").value = kt.id;
            document.getElementById("editKeycapTypeName").value = kt.name;
            document.getElementById("editKeycapTypeForm").style.display = "block";
            document.getElementById("editKeycapTypeOverlay").style.display = "block";
        } catch (err) {
            alert("Ошибка загрузки данных");
        }
    };

    document.getElementById("editKeycapTypeForm")?.addEventListener("submit", async (e) => {
        e.preventDefault();
        const id = document.getElementById("editKeycapTypeId").value;
        const name = document.getElementById("editKeycapTypeName").value.trim();
        if (!name) return alert("Введите название");
        const token = getAuthToken();
        try {
            const res = await fetch(`http://localhost:1000/keycap_types/${id}`, {
                method: "PUT",
                headers: { 
                    "Content-Type": "application/json", 
                    "Authorization": token ? `Bearer ${token}` : "" 
                },
                body: JSON.stringify({ name })
            });
            const data = await res.json();
            if (data.success) {
                alert("Тип обновлён");
                loadKeycapTypes();
                cancelEditKeycapType();
            } else {
                alert("Ошибка: " + (data.error || ""));
            }
        } catch (err) {
            alert("Ошибка сети");
        }
    });

    window.cancelEditKeycapType = function() {
        document.getElementById("editKeycapTypeForm").style.display = "none";
        document.getElementById("editKeycapTypeOverlay").style.display = "none";
        document.getElementById("editKeycapTypeForm").reset();
    };

    window.deleteKeycapType = async function(id) {
        if (!confirm("Удалить этот тип кейкапа? Это может повлиять на существующие клавиатуры, так как они ссылаются на название типа.")) return;
        const token = getAuthToken();
        try {
            const res = await fetch(`http://localhost:1000/keycap_types/${id}`, {
                method: "DELETE",
                headers: { "Authorization": token ? `Bearer ${token}` : "" }
            });
            const data = await res.json();
            if (data.success) {
                alert("Тип удалён");
                loadKeycapTypes();
            } else {
                alert("Ошибка: " + (data.error || ""));
            }
        } catch (err) {
            alert("Ошибка сети");
        }
    };

   
    loadKeyboards();
    loadKeycapTypes();
});