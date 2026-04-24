document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("loginForm");
    const errorContainer = document.getElementById("error");
    const submitBtn = document.getElementById("submitBtn");

    function showError(message) {
        errorContainer.textContent = message;
        errorContainer.style.display = 'block';
    }

    function clearError() {
        errorContainer.textContent = '';
        errorContainer.style.display = 'none';
    }

    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        clearError();

        const login = document.getElementById("login").value.trim();
        const password = document.getElementById("password").value;
        const loginType = document.getElementById("loginType").value;

        if (!login || !password) {
            showError('Пожалуйста, заполните все поля');
            return;
        }

        submitBtn.disabled = true;
        submitBtn.textContent = 'Вход...';

        try {
            const response = await fetch("http://localhost:9000/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ 
                    login, 
                    password,
                    login_type: loginType
                }),
            });

            const responseText = await response.text();
            console.log("Ответ сервера:", responseText);

            if (!response.ok) {
                try {
                    const errorData = JSON.parse(responseText);
                    showError(errorData.error || 'Ошибка авторизации');
                } catch {
                    showError(responseText || 'Ошибка авторизации');
                }
                return;
            }

            let result;
            try {
                result = JSON.parse(responseText);
            } catch (error) {
                console.error("Ошибка парсинга JSON:", error);
                showError("Ошибка формата ответа сервера");
                return;
            }

            if (!result.success) {
                showError(result.error || 'Неверный логин или пароль');
                return;
            }

  
            if (result.token) {
                localStorage.setItem('authToken', result.token);
                localStorage.setItem('userId', result.user_id);
                localStorage.setItem('userLogin', result.login);
                localStorage.setItem('isAdmin', result.is_admin);
                localStorage.setItem('authTime', Date.now().toString());
            }

        я
            if (result.redirect) {
                window.location.href = result.redirect;
            } else if (result.is_admin && loginType === 'admin') {
                window.location.href = "/admin.html";
            } else {
                window.location.href = "/user.html";
            }

        } catch (error) {
            console.error("Ошибка сети:", error);
            showError(`Ошибка сети: ${error.message}`);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Войти';
        }
    });


    function checkExistingAuth() {
        const authToken = localStorage.getItem('authToken');
        const authTime = localStorage.getItem('authTime');
        const isAdmin = localStorage.getItem('isAdmin');
        
        if (authToken && authTime) {
           
            const timeDiff = Date.now() - parseInt(authTime);
            if (timeDiff < 24 * 60 * 60 * 1000) {
                if (isAdmin === 'true') {
                    window.location.href = "/admin.html";
                } else {
                    window.location.href = "/user.html";
                }
            } else {
              
                localStorage.clear();
            }
        }
    }

    checkExistingAuth();
});