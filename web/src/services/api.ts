const BASE_URL = '/proxy/api';

export interface UserInfo {
  open_id: string;
  subscribed: boolean;
  created_at: string;
}

export const getToken = () => localStorage.getItem('token');
export const setToken = (token: string) => localStorage.setItem('token', token);
export const removeToken = () => localStorage.removeItem('token');

const customFetch = async (url: string, options: RequestInit = {}) => {
  const token = getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
};

export const login = (code: string) => {
  return customFetch('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ code }),
  });
};

export const getUser = () => {
  return customFetch('/auth/user', {
    method: 'GET',
  });
};

export const logout = () => {
  return customFetch('/auth/logout', {
    method: 'POST',
  });
};
