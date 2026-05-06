document.addEventListener("DOMContentLoaded", () => {
    const registerForm = document.getElementById("registerForm");
    const submitBtn = registerForm.querySelector('button[type="submit"]');

    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const login = document.getElementById("login").value.trim();
        const password = document.getElementById("password").value;

        if (!login || !password) {
            alert('Пожалуйста, заполните все поля');
            return;
        }

        if (password.length < 6) {
            alert('Пароль должен содержать минимум 6 символов');
            return;
        }

        submitBtn.disabled = true;
        submitBtn.textContent = 'Регистрация...';

        try {
            const response = await fetch("http://localhost:1000/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ 
                    login, 
                    password, 
                    is_admin: false 
                }),
            });

            const responseText = await response.text();
            console.log("Ответ сервера:", responseText);

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
            alert("Регистрация успешна! Теперь вы можете авторизоваться.");
            window.location.href = "/index.html";
        } catch (error) {
            console.error("Ошибка:", error);
            alert(`Ошибка сети: ${error.message}`);
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Зарегистрироваться';
        }
    });
});