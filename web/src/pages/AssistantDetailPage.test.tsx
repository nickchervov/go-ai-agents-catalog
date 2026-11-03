import { describe, expect, it, vi, beforeEach } from 'vitest';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Route, Routes } from 'react-router-dom';
import { AssistantDetailPage } from '@/pages/AssistantDetailPage';
import { renderWithProviders } from '@/test/utils';
import { assistantsApi } from '@/api';
import { useAuthStore } from '@/store/auth';

vi.mock('@/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api')>();
  return {
    ...actual,
    assistantsApi: { get: vi.fn(), run: vi.fn() },
  };
});

const assistant = {
  id: 'a-1',
  categoryId: 'cat-1',
  categoryName: 'Еда',
  name: 'Повар',
  description: 'Готовит рецепты',
  model: 'mock-smart',
  systemPrompt: 'Составь рецепт',
  exampleUserPrompt: 'курица, рис',
  isActive: true,
};

function renderDetail() {
  return renderWithProviders(
    <Routes>
      <Route path="/assistants/:id" element={<AssistantDetailPage />} />
    </Routes>,
    { route: '/assistants/a-1' },
  );
}

describe('AssistantDetailPage', () => {
  beforeEach(() => {
    useAuthStore.getState().clear();
    vi.mocked(assistantsApi.get).mockResolvedValue(assistant);
  });

  it('скрывает системный промпт от обычного пользователя', async () => {
    useAuthStore.getState().setSession('tok', {
      id: 'u',
      email: 'user@catalog.local',
      role: 'user',
    });

    renderDetail();

    expect(await screen.findByRole('heading', { name: /Повар/ })).toBeInTheDocument();
    expect(screen.queryByText(/Системный промпт/)).not.toBeInTheDocument();
  });

  it('показывает системный промпт администратору', async () => {
    useAuthStore.getState().setSession('tok', {
      id: 'u',
      email: 'admin@catalog.local',
      role: 'admin',
    });

    renderDetail();

    expect(await screen.findByText(/Системный промпт/)).toBeInTheDocument();
    expect(screen.getByText('Составь рецепт')).toBeInTheDocument();
  });

  it('запускает ассистента и отображает ответ', async () => {
    const user = userEvent.setup();
    useAuthStore.getState().setSession('tok', {
      id: 'u',
      email: 'user@catalog.local',
      role: 'user',
    });
    vi.mocked(assistantsApi.run).mockResolvedValue({
      id: 'run-1',
      assistantId: 'a-1',
      userId: 'u',
      model: 'mock-smart',
      userPrompt: 'курица, рис',
      output: 'Домашняя шаверма за 20 минут',
      status: 'success',
      error: null,
    });

    renderDetail();

    const textarea = await screen.findByLabelText(/Ваш контекст/i);
    await user.type(textarea, 'курица, рис');
    await user.click(screen.getByRole('button', { name: /Запустить/i }));

    await waitFor(() => {
      expect(assistantsApi.run).toHaveBeenCalledWith('a-1', { userPrompt: 'курица, рис' });
    });
    expect(await screen.findByTestId('run-output')).toHaveTextContent(
      'Домашняя шаверма за 20 минут',
    );
  });
});
