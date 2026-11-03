import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useAssistants, useCategories } from '@/hooks/queries';
import { useAuthStore } from '@/store/auth';
import { LoadingState, ErrorState, EmptyState } from '@/components/States';
import { ActiveBadge } from '@/components/Badges';
import { Pagination, useDebouncedValue } from '@/components/Pagination';
import { Button } from '@/components/Button';
import { toApiError } from '@/api';

const PAGE_SIZE = 10;

export function AssistantsPage() {
  const isAdmin = useAuthStore((s) => s.user?.role === 'admin');
  // Фильтры синхронизированы с query-параметрами URL.
  const [searchParams, setSearchParams] = useSearchParams();

  const page = Number(searchParams.get('page') ?? '1') || 1;
  const categoryId = searchParams.get('categoryId') ?? '';
  const includeInactive = searchParams.get('includeInactive') === 'true';
  const q = searchParams.get('q') ?? '';

  const [searchInput, setSearchInput] = useState(q);
  const debouncedQ = useDebouncedValue(searchInput);

  // При появлении внешнего q (например, из ссылки) — синхронизируем поле ввода.
  useEffect(() => {
    setSearchInput(q);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [q]);

  // Debounced-поиск пишем в URL отдельным эффектом.
  useEffect(() => {
    if (debouncedQ !== q) patch({ q: debouncedQ, page: 1 });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [debouncedQ]);

  const categoriesQuery = useCategories();
  const assistantsQuery = useAssistants({
    categoryId: categoryId || undefined,
    q: debouncedQ || undefined,
    includeInactive: isAdmin && includeInactive ? true : undefined,
    page,
    pageSize: PAGE_SIZE,
  });

  function patch(next: Record<string, string | number | boolean | undefined>) {
    const updated = new URLSearchParams(searchParams);
    Object.entries(next).forEach(([key, value]) => {
      if (value === undefined || value === '' || value === false) updated.delete(key);
      else updated.set(key, String(value));
    });
    setSearchParams(updated, { replace: true });
  }

  const categories = categoriesQuery.data ?? [];

  return (
    <div>
      <div className="title-row">
        <div className="title-row__text">
          <h1>Ассистенты</h1>
          <p className="sub">Каталог AI-ассистентов компании</p>
        </div>
        <div className="title-row__actions">
          {isAdmin && (
            <Link className="btn btn--primary" to="/admin/assistants/new">
              Новый ассистент
            </Link>
          )}
        </div>
      </div>

      <div className="filters">
        <div className="field field--search">
          <label className="field__label" htmlFor="search">
            Поиск
          </label>
          <input
            id="search"
            className="input"
            placeholder="Название или описание"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
          />
        </div>

        <div className="field">
          <label className="field__label" htmlFor="category">
            Категория
          </label>
          <select
            id="category"
            className="input"
            value={categoryId}
            onChange={(e) => patch({ categoryId: e.target.value, page: 1 })}
            disabled={categoriesQuery.isLoading}
          >
            <option value="">Все категории</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </div>

        {isAdmin && (
          <label className="check">
            <input
              type="checkbox"
              checked={includeInactive}
              onChange={(e) => patch({ includeInactive: e.target.checked, page: 1 })}
            />
            Показать неактивных
          </label>
        )}

        <div className="filters__spacer" />

        <div className="field">
          <span className="field__label" aria-hidden="true">
            &nbsp;
          </span>
          <Button variant="ghost" size="sm" onClick={() => assistantsQuery.refetch()}>
            Обновить
          </Button>
        </div>
      </div>

      {categoriesQuery.isError && (
        <p className="muted">Не удалось загрузить категории для фильтра.</p>
      )}

      {assistantsQuery.isLoading && <LoadingState label="Загружаем ассистентов…" />}

      {assistantsQuery.isError && (
        <ErrorState
          message={toApiError(assistantsQuery.error).message}
          onRetry={() => assistantsQuery.refetch()}
        />
      )}

      {assistantsQuery.isSuccess && assistantsQuery.data.assistants.length === 0 && (
        <EmptyState
          title="Ассистенты не найдены"
          description="Попробуйте изменить фильтры или поисковый запрос."
        />
      )}

      {assistantsQuery.isSuccess && assistantsQuery.data.assistants.length > 0 && (
        <>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th className="col-id">Ассистент</th>
                  <th>Категория</th>
                  <th>Модель</th>
                  {isAdmin && <th>Статус</th>}
                </tr>
              </thead>
              <tbody>
                {assistantsQuery.data.assistants.map((assistant) => (
                  <tr key={assistant.id}>
                    <td className="col-id">
                      <span className="cell-main">
                        <Link className="cell-main__title" to={`/assistants/${assistant.id}`}>
                          {assistant.name}
                        </Link>
                        <span className="cell-main__sub">{assistant.description}</span>
                      </span>
                    </td>
                    <td data-label="Категория">{assistant.categoryName ?? 'Без категории'}</td>
                    <td className="mono" data-label="Модель">
                      {assistant.model}
                    </td>
                    {isAdmin && (
                      <td data-label="Статус">
                        <ActiveBadge isActive={assistant.isActive} />
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <Pagination
            page={assistantsQuery.data.pagination.page}
            pageSize={assistantsQuery.data.pagination.pageSize}
            total={assistantsQuery.data.pagination.total}
            onPageChange={(p) => patch({ page: p })}
          />
        </>
      )}
    </div>
  );
}
