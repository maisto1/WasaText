import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

instance.interceptors.request.use(
    (config) => {
        const token = sessionStorage.getItem('authToken');
        if (token) {
            config.headers.Authorization = token;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

// Response interceptor per gestire gli errori di autorizzazione
instance.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        if (error.response && error.response.status === 401) {
            // Reindirizza al login se riceviamo un errore di autorizzazione
            window.location.href = '/#/login';
        }
        return Promise.reject(error);
    }
);

export default instance;
