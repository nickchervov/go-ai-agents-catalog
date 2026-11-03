import { useState } from 'react';
import { Link } from 'react-router-dom';
import { StatusBadge } from './Badges';
import { Button } from './Button';
import type { AssistantRun } from '@/types/api';

function formatDate(value?: string | null): string {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('ru-RU');
}

function truncate(value?: string | null, max = 72): string {
  if (!value) return '—';
  return value.length > max ? `${value.slice(0, max)}…` : value;
}

interface RunsTableProps {
  runs: AssistantRun[];
  /** Показывать ли колонку с ассистентом (для /admin/runs полезно всегда). */
  showAssistant?: boolean;
}

/** Таблица запусков, переиспользуется на /runs/my и /admin/runs. */
export function RunsTable({ runs, showAssistant = true }: RunsTableProps) {
  const [opened, setOpened] = useState<AssistantRun | null>(null);

  return (
    <>
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              {showAssistant && <th className="col-id">Ассистент</th>}
              <th>Статус</th>
              <th>Дата</th>
              <th>Промпт пользователя</th>
              <th>Ответ</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {runs.map((run) => (
              <tr key={run.id}>
                {showAssistant && (
                  <td data-label="Ассистент">
                    <span className="cell-main">
                      <Link className="cell-main__title" to={`/assistants/${run.assistantId}`}>
                        {run.assistantName ?? run.assistantId.slice(0, 8)}
                      </Link>
                      {run.categoryName && <span className="cell-main__sub">{run.categoryName}</span>}
                    </span>
                  </td>
                )}
                <td data-label="Статус">
                  <StatusBadge status={run.status} />
                </td>
                <td className="muted nowrap" data-label="Дата">
                  {formatDate(run.createdAt)}
                </td>
                <td data-label="Промпт">
                  <span className="truncate" title={run.userPrompt}>
                    {truncate(run.userPrompt)}
                  </span>
                </td>
                <td data-label="Ответ">
                  <span className="truncate" title={run.output ?? run.error ?? ''}>
                    {run.status === 'failed' ? (
                      <span style={{ color: 'var(--negative)' }}>{truncate(run.error)}</span>
                    ) : (
                      truncate(run.output)
                    )}
                  </span>
                </td>
                <td className="n">
                  <Button size="sm" onClick={() => setOpened(run)}>
                    Открыть
                  </Button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {opened && (
        <div
          className="modal"
          role="dialog"
          aria-modal="true"
          aria-label="Детали запуска"
          onClick={() => setOpened(null)}
        >
          <div className="modal__panel" onClick={(e) => e.stopPropagation()}>
            <div className="modal__head">
              <h2>Запуск {opened.assistantName ?? ''}</h2>
              <StatusBadge status={opened.status} />
            </div>

            <div className="kv">
              <div className="kv__row">
                <span className="kv__label">Модель</span>
                <span className="mono">{opened.model}</span>
              </div>
              <div className="kv__row">
                <span className="kv__label">Дата</span>
                <span>{formatDate(opened.createdAt)}</span>
              </div>
            </div>

            <h3 className="panel-title" style={{ margin: '18px 0 0', fontSize: 'var(--fs-section)' }}>
              Промпт
            </h3>
            <div className="run-output">{opened.userPrompt}</div>

            <h3 className="panel-title" style={{ margin: '18px 0 0', fontSize: 'var(--fs-section)' }}>
              {opened.status === 'failed' ? 'Ошибка' : 'Ответ'}
            </h3>
            {opened.status === 'failed' ? (
              <div className="alert alert--error" style={{ marginTop: 8 }}>
                {opened.error ?? 'Неизвестная ошибка'}
              </div>
            ) : (
              <div className="run-output run-output--boxed" data-testid="run-modal-output">
                {opened.output ?? 'Пустой ответ'}
              </div>
            )}

            <div className="form-actions" style={{ marginTop: 18 }}>
              <Button onClick={() => setOpened(null)}>Закрыть</Button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

export { formatDate };
