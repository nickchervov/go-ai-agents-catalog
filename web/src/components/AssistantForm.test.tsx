import { describe, expect, it, vi } from 'vitest';
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { AssistantForm } from '@/components/AssistantForm';
import { renderWithProviders } from '@/test/utils';
import { categoriesApi } from '@/api';

vi.mock('@/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/api')>();
  return {
    ...actual,
    categoriesApi: { list: vi.fn() },
  };
});

describe('AssistantForm', () => {
  it('не отправляет форму с пустыми обязательными полями (валидация)', async () => {
    vi.mocked(categoriesApi.list).mockResolvedValue([
      { id: '11111111-1111-1111-1111-111111111111', name: 'Еда' },
    ]);
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

  it('отправляет корректно заполненную форму', async () => {
    vi.mocked(categoriesApi.list).mockResolvedValue([
      { id: '11111111-1111-1111-1111-111111111111', name: 'Еда' },
    ]);
    const user = userEvent.setup();
    const onSubmit = vi.fn();

    renderWithProviders(
      <AssistantForm submitting={false} submitLabel="Создать" onSubmit={onSubmit} onCancel={() => {}} />,
    );

    await user.type(screen.getByLabelText(/Название/i), 'Повар');
    await user.type(screen.getByLabelText(/Описание/i), 'Готовит рецепты');
    await user.selectOptions(screen.getByLabelText(/Категория/i), '11111111-1111-1111-1111-111111111111');
    await user.type(screen.getByLabelText(/Системный промпт/i), 'Составь рецепт');
    await user.click(screen.getByRole('button', { name: /Создать/i }));

    expect(onSubmit).toHaveBeenCalledTimes(1);
    const values = onSubmit.mock.calls[0][0];
    expect(values).toMatchObject({
      name: 'Повар',
      description: 'Готовит рецепты',
      systemPrompt: 'Составь рецепт',
      categoryId: '11111111-1111-1111-1111-111111111111',
    });
  });
});
