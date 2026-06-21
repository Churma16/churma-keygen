const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1';

async function request(path, options = {}) {
    const token = localStorage.getItem('admin_token');
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
    });

    const contentType = res.headers.get('content-type');
    if (contentType && contentType.includes('text/plain')) {
        const text = await res.text();
        if (!res.ok) {
            throw new Error(text || 'Fetch failed');
        }
        return text;
    }

    let body = null;
    try {
        body = await res.json();
    } catch (e) {}

    if (!res.ok) {
        if (res.status === 401) {
            localStorage.removeItem('admin_token');
            localStorage.removeItem('admin_username');
            localStorage.removeItem('admin_role');
            window.location.reload();
        }
        throw new Error(body?.meta?.message || body?.error || 'API Request failed');
    }

    return body;
}

export const baseApi = {
    get: (path, options) => request(path, { method: 'GET', ...options }),
    post: (path, data, options) => request(path, { method: 'POST', body: JSON.stringify(data), ...options }),
    put: (path, data, options) => request(path, { method: 'PUT', body: JSON.stringify(data), ...options }),
    delete: (path, options) => request(path, { method: 'DELETE', ...options }),
};
