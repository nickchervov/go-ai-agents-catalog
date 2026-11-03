import type { RunStatus } from '@/types/api';

const LABELS: Record<RunStatus, string> = {
  pending: 'В ожидании',
  success: 'Успешно',
  failed: 'Ошибка',
};

const CLASS: Record<RunStatus, string> = {
  pending: 'pill pill--pending',
  success: 'pill pill--success',
  failed: 'pill pill--failed',
};

export function StatusBadge({ status }: { status: RunStatus }) {
  return (
    <span className={CLASS[status] ?? 'pill'}>
      <span className="pill__dot" />
      {LABELS[status] ?? status}
    </span>
  );
}

export function ActiveBadge({ isActive }: { isActive: boolean }) {
  return (
    <span className={`pill ${isActive ? 'pill--active' : 'pill--inactive'}`}>
      <span className="pill__dot" />
      {isActive ? 'Активен' : 'Выключен'}
    </span>
  );
}
