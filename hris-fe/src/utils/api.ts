// hris-fe/src/utils/api.ts

import axios from 'axios';

// Ambil base URL dari environment variable.
// Jika tidak ada (saat pengembangan lokal), gunakan localhost:6969
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:6969';

// Buat instance Axios dengan base URL ini
const api = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000, // Timeout request dalam milidetik (10 detik), opsional
    headers: {
        'Content-Type': 'application/json',
    },
});

// Anda bisa menambahkan interceptor di sini jika diperlukan (misalnya untuk token autentikasi)
// api.interceptors.request.use(
//   (config) => {
//     const token = localStorage.getItem('token'); // Contoh: ambil token dari localStorage
//     if (token) {
//       config.headers.Authorization = `Bearer ${token}`;
//     }
//     return config;
//   },
//   (error) => {
//     return Promise.reject(error);
//   }
// );

export default api;