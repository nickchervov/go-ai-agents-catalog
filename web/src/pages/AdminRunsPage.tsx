import { useSearchParams } from 'react-router-dom';
import { useAllRuns, useAssistants } from '@/hooks/queries';
import { toApiError } from '@/api';
import { EmptyState, ErrorState, LoadingState } from '@/components/States';
import { Pagination } from '@/components/Pagination';
import { RunsTable } from '@/components/RunsTable';
import type { RunStatus } from '@/types/api';

const PAGE_SIZE = 20;
const STATUSES: { value: RunStatus | ''; label: string }[] = [
  { value: '', label: 'Все статусы' },
  { value: 'pending', label: 'В ожидании' },
  { value: 'success', label: 'Успешно' },
  { value: 'failed', label: 'Ошибка' },
];

export function AdminRunsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = Number(searchParams.get('page') ?? '1') || 1;
  const status = (searchParams.get('status') ?? '') as RunStatus | '';
  const assistantId = searchParams.get('assistantId') ?? '';

  const query = useAllRuns({
    status: status || undefined,
    assistantId: assistantId || undefined,
    page,
    pageSize: PAGE_SIZE,
  });

  // Для фильтра по ассистенту подгружаем список (включая неактивных).
  const assistantsQuery = useAssistants({ includeInactive: true, pageSize: 100, page: 1 });

  const patch = (next: Record<string, string | number | undefined>) => {
    const updated = new URLSearchParams(searchParams);
    Object.entries(next).forEach(([key, value]) => {
      if (value === undefined || value === '') updated.delete(key);
      else updated.set(key, String(value));
    });
    setSearchParams(updated, { replace: true });
  };

  const total = query.data?.pagination.total ?? 0;
  const onPage = query.data?.runs.length ?? 0;
  const assistantLabel =
    (assistantsQuery.data?.assistants ?? []).find((a) => a.id === assistantId)?.name ?? 'Все ассистенты';

  return (
    <div>
      <div className="title-row">
        <div className="title-row__text">
          <h1>Запуски ассистентов</h1>
          <p className="sub">Все обращения пользователей к ассистентам</p>
        </div>
      </div>

      <div className="kpis">
        <div className="kpi">
          <span className="kpi__label">Всего запусков</span>
          <span className="kpi__value num">{total}</span>
        </div>
        <div className="kpi">
          <span className="kpi__label">На странице</span>
          <span className="kpi__value num">{onPage}</span>
        </div>
        <div className="kpi">
          <span className="kpi__label">Ассистент</span>
          <span className="kpi__value">{assistantLabel}</span>
        </div>
      </div>

      <div className="filters">
        <div className="field">
          <label className="field__label" htmlFor="assistantId">
            Ассистент
          </label>
          <select
            id="assistantId"
            className="input"
            value={assistantId}
            onChange={(e) => patch({ assistantId: e.target.value, page: 1 })}
            disabled={assistantsQuery.isLoading}
          >
            <option value="">Все ассистенты</option>
            {(assistantsQuery.data?.assistants ?? []).map((a) => (
              <option key={a.id} value={a.id}>
                {a.name}
              </option>
            ))}
          </select>
        </div>

        <div className="field">
          <label className="field__label" htmlFor="status">
            Статус
          </label>
          <select
            id="status"
            className="input"
            value={status}
            onChange={(e) => patch({ status: e.target.value, page: 1 })}
          >
            {STATUSES.map((s) => (
              <option key={s.value} value={s.value}>
                {s.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {query.isLoading && <LoadingState label="Загружаем запуски…" />}
      {query.isError && (
        <ErrorState message={toApiError(query.error).message} onRetry={() => query.refetch()} />
      )}

      {query.isSuccess && query.data.runs.length === 0 && (
        <EmptyState title="Запусков не найдено" description="Измените фильтры." />
      )}

      {query.isSuccess && query.data.runs.length > 0 && (
        <>
          <RunsTable runs={query.data.runs} showAssistant />
          <Pagination
            page={query.data.pagination.page}
            pageSize={query.data.pagination.pageSize}
            total={query.data.pagination.total}
            onPageChange={(p) => patch({ page: p })}
          />
        </>
      )}
    </div>
  );
}
