import axios from "axios";

const BASE_URL = import.meta.env.VITE_API_URL;

export const api = axios.create({
    baseURL: BASE_URL,
    withCredentials: true,
});

let isRefreshing = false;
let failedQueue: {
    resolve: (value?: unknown) => void;
    reject: (reason?: unknown) => void;
}[] = [];

const processQueue = (error: unknown) => {
    failedQueue.forEach((prom) => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve();
        }
    });

    failedQueue = [];
};

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // Only attempt refresh on 401 errors, not on the refresh endpoint itself,
        // and not if we've already retried this request
        if (
            error.response?.status === 401 &&
            !originalRequest._retry &&
            !originalRequest.url?.includes("/refresh")
        ) {
            if (isRefreshing) {
                // Another refresh is already in progress — queue this request
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                })
                    .then(() => api(originalRequest))
                    .catch((err) => Promise.reject(err));
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                // Call the refresh endpoint
                await axios.post(
                    `${BASE_URL}/refresh`,
                    {},
                    { withCredentials: true }
                );

                // Refresh succeeded
                processQueue(null);

                // Retry the original request
                return api(originalRequest);
            } catch (refreshError) {
                // Refresh failed — reject all queued requests
                processQueue(refreshError);
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }

        return Promise.reject(error);
    }
);

export default api;