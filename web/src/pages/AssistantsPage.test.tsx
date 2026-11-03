import { describe, expect, it, vi, beforeEach } from 'vitest';
import { screen } from '@testing-library/react';
import { AssistantsPage } from '@/pages/AssistantsPage';
import { renderWithProviders } from '@/test/utils';
import { assistantsApi, categoriesApi } from '@/api';
import { useAuthStore } from '@/store/auth';

vi.mock('@/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api')>();
  return {
    ...actual,
    assistantsApi: { list: vi.fn() },
    categoriesApi: { list: vi.fn() },
  };
});

const category = { id: 'cat-1', name: 'Еда' };
const assistant = {
  id: 'a-1',
  categoryId: 'cat-1',
  categoryName: 'Еда',
  name: 'Повар',
  description: 'Готовит рецепты',
  model: 'mock-smart',
  isActive: true,
};

describe('AssistantsPage', () => {
  beforeEach(() => {
    useAuthStore.getState().clear();
    vi.mocked(categoriesApi.list).mockResolvedValue([category]);
  });

  it('отображает карточки ассистентов из API', async () => {
    vi.mocked(assistantsApi.list).mockResolvedValue({
      assistants: [assistant],
      pagination: { page: 1, pageSize: 10, total: 1 },
    });

    renderWithProviders(<AssistantsPage />);

    expect(await screen.findByText('Повар')).toBeInTheDocument();
    expect(screen.getByText('Готовит рецепты')).toBeInTheDocument();
    // Ссылка ведёт на карточку ассистента.
    expect(screen.getByRole('link', { name: /Повар/ })).toHaveAttribute(
      'href',
      '/assistants/a-1',
    );
  });

  it('показывает пустое состояние, когда ассистентов нет', async () => {
    vi.mocked(assistantsApi.list).mockResolvedValue({
      assistants: [],
      pagination: { page: 1, pageSize: 10, total: 0 },
    });

    renderWithProviders(<AssistantsPage />);
    expect(await screen.findByTestId('empty-state')).toBeInTheDocument();
  });

  it('передаёт includeInactive только для администратора', async () => {
    useAuthStore.getState().setSession('tok', {
      id: 'u',
      email: 'admin@catalog.local',
      role: 'admin',
    });
    vi.mocked(assistantsApi.list).mockResolvedValue({
      assistants: [],
      pagination: { page: 1, pageSize: 10, total: 0 },
    });

    renderWithProviders(<AssistantsPage />, { route: '/assistants?includeInactive=true' });

    // Чекбокс доступен админу.
    expect(await screen.findByLabelText(/Показать неактивных/i)).toBeChecked();
  });
});
