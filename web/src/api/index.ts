import { http, setTokenGetter, setUnauthorizedHandler } from './http';
import { useAuthStore } from '@/store/auth';

/**
 * Связывает HTTP-клиент со стором авторизации:
 * - клиент читает актуальный токен перед каждым запросом;
 * - при 401 стор очищается (авто-разлогин).
 *
 * Вызывается один раз при старте приложения (см. main.tsx).
 */
export function wireHttpAuth() {
  setTokenGetter(() => useAuthStore.getState().token);
  setUnauthorizedHandler(() => useAuthStore.getState().clear());
}

export { http };
export { authApi } from './auth';
export { categoriesApi } from './categories';
export { assistantsApi } from './assistants';
export { runsApi } from './runs';
export { ApiError, toApiError } from './http';
