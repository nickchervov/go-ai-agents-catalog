import { useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useAssistant, useRunAssistant } from '@/hooks/queries';
import { toApiError } from '@/api';
import { useAuthStore } from '@/store/auth';
import { Alert, ErrorState, LoadingState } from '@/components/States';
import { ActiveBadge, StatusBadge } from '@/components/Badges';
import { Button } from '@/components/Button';
import type { AssistantRun } from '@/types/api';

export function AssistantDetailPage() {
  const { id } = useParams<{ id: string }>();
  const isAdmin = useAuthStore((s) => s.user?.role === 'admin');
  const assistantsNavigate = useNavigate();

  const query = useAssistant(id);
  const runMutation = useRunAssistant();

  const [userPrompt, setUserPrompt] = useState('');
  const [runResult, setRunResult] = useState<AssistantRun | null>(null);
  const [runError, setRunError] = useState<string | null>(null);

  if (query.isLoading) return <LoadingState label="Загружаем ассистента…" />;
  if (query.isError) {
    return <ErrorState message={toApiError(query.error).message} onRetry={() => query.refetch()} />;
  }
  if (!query.data) return null;

  const assistant = query.data;
  const canRun = assistant.isActive;

  const handleRun = async () => {
    if (!id) return;
    setRunError(null);
    setRunResult(null);
    try {
      const result = await runMutation.mutateAsync({ id, payload: { userPrompt } });
      setRunResult(result);
    } catch (err) {
      setRunError(toApiError(err).message);
    }
  };

  return (
    <div>
      <div className="title-row">
        <div className="title-row__text">
          <h1>
            {assistant.name} {isAdmin && <ActiveBadge isActive={assistant.isActive} />}
          </h1>
          <p className="sub">{assistant.description}</p>
        </div>
        <div className="title-row__actions">
          {isAdmin && (
            <Link className="btn" to={`/admin/assistants/${assistant.id}/edit`}>
              Редактировать
            </Link>
          )}
          <Button variant="ghost" onClick={() => assistantsNavigate('/assistants')}>
            К каталогу
          </Button>
        </div>
      </div>

      <div className="grid">
        <section className="panel grid__col-8">
          <div className="panel__header">
            <h2>Запуск ассистента</h2>
          </div>
          <p className="muted">
            Передайте свой контекст - backend добавит системный промпт и обратится к LLM.
          </p>

          {!canRun && <Alert variant="error">Ассистент выключен, запуск недоступен.</Alert>}

          <div className="field">
            <label className="field__label" htmlFor="userPrompt">
              Ваш контекст
            </label>
            <textarea
              id="userPrompt"
              className="input"
              value={userPrompt}
              placeholder={assistant.exampleUserPrompt ?? 'Введите запрос…'}
              onChange={(e) => setUserPrompt(e.target.value)}
            />
          </div>

          {runError && <Alert variant="error">{runError}</Alert>}

          <div className="form-actions">
            <Button
              variant="primary"
              disabled={!canRun || userPrompt.trim().length === 0 || runMutation.isPending}
              onClick={handleRun}
            >
              {runMutation.isPending ? 'Выполняется…' : 'Запустить'}
            </Button>
          </div>

          {runResult && (
            <>
              <div className="panel__header" style={{ marginTop: 'var(--gap)' }}>
                <h2>Результат</h2>
                <StatusBadge status={runResult.status} />
              </div>
              {runResult.status === 'failed' ? (
                <Alert variant="error">{runResult.error ?? 'Не удалось выполнить запуск.'}</Alert>
              ) : (
                <div className="run-output run-output--boxed" data-testid="run-output">
                  {runResult.output ?? 'Пустой ответ'}
                </div>
              )}
              <p className="muted" style={{ marginTop: 10 }}>
                Запуск сохранён в <Link to="/runs/my">истории</Link>.
              </p>
            </>
          )}
        </section>

        <aside className="panel grid__col-4">
          <div className="panel__header">
            <h2>Свойства</h2>
          </div>
          <div className="kv">
            <div className="kv__row">
              <span className="kv__label">Категория</span>
              <span>{assistant.categoryName ?? '—'}</span>
            </div>
            <div className="kv__row">
              <span className="kv__label">Модель</span>
              <span className="mono">{assistant.model}</span>
            </div>
            {assistant.exampleUserPrompt && (
              <div className="kv__row">
                <span className="kv__label">Пример контекста</span>
                <span className="muted">{assistant.exampleUserPrompt}</span>
              </div>
            )}
            {/*
              Системный промпт — внутренняя настройка ассистента.
              Показываем только администратору (обосновано в README).
            */}
            {isAdmin && assistant.systemPrompt && (
              <div className="kv__row">
                <span className="kv__label">Системный промпт</span>
                <span className="muted">{assistant.systemPrompt}</span>
              </div>
            )}
          </div>
        </aside>
      </div>
    </div>
  );
}
