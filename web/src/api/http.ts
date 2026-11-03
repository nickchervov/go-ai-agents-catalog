import axios, { AxiosError } from 'axios';
import type { ErrorCode, ErrorResponse } from '@/types/api';

/**
 * Единый HTTP-клиент. Все компоненты ходят в API только через этот модуль
 * (слой api/*), чтобы не собирать URL и заголовки вручную.
 */
// По умолчанию ходим через /api: этот префикс проксируют Vite и nginx.
// Так пути API не пересекаются с путями SPA (/assistants, /runs, /admin).
const baseURL = import.meta.env.VITE_API_BASE_URL?.trim() || '/api';

export const http = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
});

/** Ошибка API, нормализованная для отображения в UI. */
export class ApiError extends Error {
  code: ErrorCode | 'NETWORK_ERROR';
  status: number | null;

  constructor(message: string, code: ErrorCode | 'NETWORK_ERROR', status: number | null) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

/** Колбэк, вызываемый при 401 — используется для авто-разлогина. */
let onUnauthorized: (() => void) | null = null;
export function setUnauthorizedHandler(handler: (() => void) | null) {
  onUnauthorized = handler;
}

export function toApiError(error: unknown): ApiError {
  if (error instanceof ApiError) return error;

  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ErrorResponse>;
    const status = axiosError.response?.status ?? null;
    const payload = axiosError.response?.data;

    if (payload?.error?.code && payload.error.message) {
      return new ApiError(payload.error.message, payload.error.code, status);
    }
    if (axiosError.code === 'ECONNABORTED') {
      return new ApiError('Превышено время ожидания ответа сервера', 'NETWORK_ERROR', status);
    }
    if (!axiosError.response) {
      return new ApiError(
        'Не удалось связаться с сервером. Проверьте, запущен ли backend.',
        'NETWORK_ERROR',
        null,
      );
    }
    return new ApiError(axiosError.message, 'NETWORK_ERROR', status);
  }

  return new ApiError('Неизвестная ошибка', 'INTERNAL_ERROR', null);
}

http.interceptors.request.use((config) => {
  // Токен берём из единого хранилища (Zustand) через getState,
  // чтобы избежать циклических импортов между api-слоем и стором.
  const token = getStoredToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = axios.isAxiosError(error) ? error.response?.status : undefined;
    if (status === 401) {
      onUnauthorized?.();
    }
    return Promise.reject(toApiError(error));
  },
);

/**
 * Ленивый доступ к токену, чтобы не импортировать стор на этапе инициализации.
 * Стор регистрирует себя через setTokenGetter при загрузке.
 */
let tokenGetter: () => string | null = () => null;
export function setTokenGetter(getter: () => string | null) {
  tokenGetter = getter;
}
function getStoredToken(): string | null {
  return tokenGetter();
}
