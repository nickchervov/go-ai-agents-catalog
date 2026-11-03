import { useSearchParams } from 'react-router-dom';
import { useMyRuns } from '@/hooks/queries';
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

export function MyRunsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = Number(searchParams.get('page') ?? '1') || 1;
  const status = (searchParams.get('status') ?? '') as RunStatus | '';

  const query = useMyRuns({
    status: status || undefined,
    page,
    pageSize: PAGE_SIZE,
  });

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

  return (
    <div>
      <div className="title-row">
        <div className="title-row__text">
          <h1>Мои запуски</h1>
          <p className="sub">История обращений к ассистентам</p>
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
          <span className="kpi__label">Фильтр статуса</span>
          <span className="kpi__value">
            {STATUSES.find((s) => s.value === status)?.label ?? 'Все статусы'}
          </span>
        </div>
      </div>

      <div className="filters">
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

      {query.isLoading && <LoadingState label="Загружаем историю…" />}
      {query.isError && (
        <ErrorState message={toApiError(query.error).message} onRetry={() => query.refetch()} />
      )}

      {query.isSuccess && query.data.runs.length === 0 && (
        <EmptyState
          title="Запусков пока нет"
          description="Откройте ассистента из каталога и запустите его."
        />
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
