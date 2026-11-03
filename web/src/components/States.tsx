import type { ReactNode } from 'react';

/** Плоские серые полосы реального размера — без мерцания и без сдвига макета. */
export function Skeleton({ width = '100%', height = 14 }: { width?: string | number; height?: number }) {
  return <div className="skeleton" style={{ width, height }} aria-hidden="true" />;
}

export function LoadingState({ label = 'Загрузка…' }: { label?: string }) {
  return (
    <div className="skeleton-set" data-testid="loading-state" role="status" aria-label={label}>
      <Skeleton width="42%" height={18} />
      <Skeleton width="100%" height={40} />
      <Skeleton width="100%" height={40} />
      <Skeleton width="100%" height={40} />
      <Skeleton width="100%" height={40} />
    </div>
  );
}

export function ErrorState({
  message = 'Произошла ошибка',
  onRetry,
}: {
  message?: string;
  onRetry?: () => void;
}) {
  return (
    <div className="empty" data-testid="error-state">
      <p className="empty__title">Не удалось загрузить данные</p>
      <p>{message}</p>
      {onRetry && (
        <div className="empty__actions">
          <button type="button" className="btn" onClick={onRetry}>
            Повторить
          </button>
        </div>
      )}
    </div>
  );
}

export function EmptyState({
  title = 'Ничего не найдено',
  description,
  children,
}: {
  title?: string;
  description?: string;
  children?: ReactNode;
}) {
  return (
    <div className="empty" data-testid="empty-state">
      <p className="empty__title">{title}</p>
      {description && <p>{description}</p>}
      {children && <div className="empty__actions">{children}</div>}
    </div>
  );
}

export function Alert({
  variant = 'error',
  children,
}: {
  variant?: 'error' | 'success';
  children: ReactNode;
}) {
  return (
    <div className={`alert alert--${variant}`} role={variant === 'error' ? 'alert' : 'status'}>
      {children}
    </div>
  );
}
