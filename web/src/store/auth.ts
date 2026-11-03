import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Role, User } from '@/types/api';

export const AUTH_STORAGE_KEY = 'ai-catalog-auth';

interface AuthState {
  token: string | null;
  user: User | null;
  /** Роль из токена; дублируется в user.role, оставлено для удобных проверок. */
  role: Role | null;
  setSession: (token: string, user: User) => void;
  clear: () => void;
  isAuthenticated: () => boolean;
}

/**
 * Единое хранилище токена и роли.
 * persist() пишет состояние в localStorage под ключом AUTH_STORAGE_KEY,
 * это единственное место, где хранится JWT на frontend.
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      role: null,
      setSession: (token, user) => set({ token, user, role: user.role }),
      clear: () => set({ token: null, user: null, role: null }),
      isAuthenticated: () => Boolean(get().token),
    }),
    {
      name: AUTH_STORAGE_KEY,
      partialize: (state) => ({ token: state.token, user: state.user, role: state.role }),
    },
  ),
);
