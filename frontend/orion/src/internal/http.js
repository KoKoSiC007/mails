import { isLogged } from "./utils";

export class APIClient {
    baseUrl = "http://localhost:8080/api/v1"

    async request(endpoint, method = 'GET', data = null, headers = {}) {
        try {
            if (isLogged()) {
                headers['Authorization'] = `Bearer ${localStorage.getItem("JWT")}`
            }
            const requestOptions = {
                method: method.toUpperCase(),
                headers: {
                    'Content-Type': 'application/json',
                    ...headers
                }
            };

            if (data) {
                requestOptions.body = JSON.stringify(data);
            }
            console.warn(`${this.baseUrl}${endpoint}`)
            const response = await fetch(`${this.baseUrl}${endpoint}`, requestOptions);
            const status = response.status;

            // Проверяем статус ответа
            if (!response.ok) {
                if (status == 401) {
                    localStorage.clear()
                }
                console.error(response)
                const error = await response.text();
                throw new Error(`HTTP error ${status}: ${error}`);
            }

            // Для статусов 204 No Content не возвращаем тело ответа
            if (status === 204) {
                return null;
            }

            // Пытаемся парсить JSON для остальных статусов
            try {
                return await response.json();
            } catch (error) {
                return response.text();
            }
        } catch (error) {
            // Обрабатываем сетевые ошибки
            if (error.name === 'TypeError') {
                console.error(error)
                throw new Error('Network error: Unable to connect to the server');
            }
            throw error;
        }
    }

    async get(endpoint, headers = {}) {
        return this.request(endpoint, 'GET', null, headers);
    }

    async post(endpoint, data, headers = {}) {
        return this.request(endpoint, 'POST', data, headers);
    }

    async put(endpoint, data, headers = {}) {
        return this.request(endpoint, 'PUT', data, headers);
    }

    async delete(endpoint, headers = {}) {
        return this.request(endpoint, 'DELETE', null, headers);
    }

    async login(endpoint) {
        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const response = await fetch(`${this.baseUrl}${endpoint}`, requestOptions);
        const status = response.status;
        if (!response.ok) {
            const error = await response.text();
            throw new Error(`HTTP error ${status}: ${error}`);
        }

        localStorage.setItem("JWT", response.headers.get("Authorization"))
    }
}