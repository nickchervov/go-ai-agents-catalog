import { useEffect, useState } from 'react';
import { Button } from './Button';

interface PaginationProps {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  disabled?: boolean;
}

function rangeLabel(page: number, pageSize: number, total: number): string {
  if (total === 0) return '0';
  const from = (page - 1) * pageSize + 1;
  const to = Math.min(page * pageSize, total);
  return `${from}–${to} из ${total}`;
}

/** Постраничная навигация, общая для каталога и истории запусков. */
export function Pagination({
  page,
  pageSize,
  total,
  onPageChange,
  disabled = false,
}: PaginationProps) {
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const canPrev = page > 1;
  const canNext = page < totalPages;

  return (
    <nav className="pagination" aria-label="Пагинация">
      <span className="num">
        Показано {rangeLabel(page, pageSize, total)} · страница {page} из {totalPages}
      </span>
      <span className="pagination__nav">
        <Button
          disabled={!canPrev || disabled}
          onClick={() => onPageChange(page - 1)}
          aria-label="Предыдущая страница"
        >
          Назад
        </Button>
        <Button
          disabled={!canNext || disabled}
          onClick={() => onPageChange(page + 1)}
          aria-label="Следующая страница"
        >
          Вперёд
        </Button>
      </span>
    </nav>
  );
}

/** Хук debounce для строки поиска. */
export function useDebouncedValue<T>(value: T, delay = 350): T {
  const [debounced, setDebounced] = useState(value);
  useEffect(() => {
    const id = setTimeout(() => setDebounced(value), delay);
    return () => clearTimeout(id);
  }, [value, delay]);
  return debounced;
}
