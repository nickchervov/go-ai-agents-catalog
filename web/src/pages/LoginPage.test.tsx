import { describe, expect, it, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LoginPage } from '@/pages/LoginPage';
import { renderWithProviders } from '@/test/utils';
import { authApi } from '@/api';
import { useAuthStore } from '@/store/auth';

// Мокаем только сетевой слой — компонент и стор остаются настоящими.
vi.mock('@/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api')>();
  return {
    ...actual,
    authApi: { dummyLogin: vi.fn() },
  };
});

describe('LoginPage', () => {
  beforeEach(() => {
    useAuthStore.getState().clear();
  });

  it('логинит администратора и сохраняет токен и роль в сторе', async () => {
    const user = userEvent.setup();
    const token = 'jwt-token-admin';
    vi.mocked(authApi.dummyLogin).mockResolvedValue({
      token,
      user: { id: 'u-1', email: 'admin@catalog.local', role: 'admin' },
    });

    renderWithProviders(<LoginPage />);

    // По умолчанию выбрана роль user — переключаем на admin.
    await user.click(screen.getByRole('button', { name: /Администратор/i }));
    await user.click(screen.getByRole('button', { name: /Войти/i }));

    await waitFor(() => {
      expect(authApi.dummyLogin).toHaveBeenCalledWith('admin');
    });
    expect(useAuthStore.getState().token).toBe(token);
    expect(useAuthStore.getState().role).toBe('admin');
  });

  it('показывает ошибку при неудачном входе', async () => {
    const user = userEvent.setup();
    const { ApiError } = await import('@/api/http');
    vi.mocked(authApi.dummyLogin).mockRejectedValue(
      new ApiError('Недостаточно прав', 'FORBIDDEN', 403),
    );

    renderWithProviders(<LoginPage />);
    await user.click(screen.getByRole('button', { name: /Войти/i }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Недостаточно прав');
    expect(useAuthStore.getState().token).toBeNull();
  });
});
