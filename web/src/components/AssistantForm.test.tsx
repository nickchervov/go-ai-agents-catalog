import { beforeEach, describe, expect, it, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AssistantForm } from '@/components/AssistantForm';
import { renderWithProviders } from '@/test/utils';
import { categoriesApi, modelsApi } from '@/api';

vi.mock('@/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api')>();
  return {
    ...actual,
    categoriesApi: { list: vi.fn() },
    modelsApi: { list: vi.fn() },
  };
});

describe('AssistantForm', () => {
  beforeEach(() => {
    vi.mocked(categoriesApi.list).mockResolvedValue([
      { id: '11111111-1111-1111-1111-111111111111', name: 'Еда' },
    ]);
    vi.mocked(modelsApi.list).mockResolvedValue(['mock-smart', 'gpt-4o-mini', 'o1']);
  });
  it('не отправляет форму с пустыми обязательными полями (валидация)', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    renderWithProviders(
      <AssistantForm submitting={false} submitLabel="Создать" onSubmit={onSubmit} onCancel={() => {}} />,
    );

    await user.click(screen.getByRole('button', { name: /Создать/i }));

    expect(onSubmit).not.toHaveBeenCalled();
    // Показаны сообщения об обязательных полях.
    expect(screen.getAllByText(/обязательно/i).length).toBeGreaterThan(0);
  });

  it('фильтрует список моделей по названию', async () => {
    const user = userEvent.setup();

    renderWithProviders(
      <AssistantForm submitting={false} submitLabel="Создать" onSubmit={vi.fn()} onCancel={() => {}} />,
    );

    await screen.findByRole('option', { name: 'gpt-4o-mini' });
    await user.type(screen.getByLabelText(/Поиск модели/i), 'gpt');

    expect(screen.getByRole('option', { name: 'gpt-4o-mini' })).toBeInTheDocument();
    expect(screen.queryByRole('option', { name: /^o1$/ })).not.toBeInTheDocument();
    // Текущая выбранная модель остаётся в списке, даже если не совпала с поиском.
    expect(screen.getByRole('option', { name: 'mock-smart' })).toBeInTheDocument();
  });

  it('отправляет корректно заполненную форму с выбранной моделью', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    renderWithProviders(
      <AssistantForm submitting={false} submitLabel="Создать" onSubmit={onSubmit} onCancel={() => {}} />,
    );

    await user.type(screen.getByLabelText(/Название/i), 'Повар');
    await user.type(screen.getByLabelText(/Описание/i), 'Готовит рецепты');
    await user.selectOptions(screen.getByLabelText(/Категория/i), '11111111-1111-1111-1111-111111111111');
    // Список моделей приходит из GET /models.
    expect(await screen.findByRole('option', { name: 'gpt-4o-mini' })).toBeInTheDocument();
    await user.selectOptions(screen.getByLabelText(/Модель/i), 'gpt-4o-mini');
    await user.type(screen.getByLabelText(/Системный промпт/i), 'Составь рецепт');
    await user.click(screen.getByRole('button', { name: /Создать/i }));

    expect(onSubmit).toHaveBeenCalledTimes(1);
    const values = onSubmit.mock.calls[0][0];
    expect(values).toMatchObject({
      name: 'Повар',
      description: 'Готовит рецепты',
      model: 'gpt-4o-mini',
      systemPrompt: 'Составь рецепт',
      categoryId: '11111111-1111-1111-1111-111111111111',
    });
  });
});
