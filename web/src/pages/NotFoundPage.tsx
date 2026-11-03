import { Link } from 'react-router-dom';
import { EmptyState } from '@/components/States';

export function NotFoundPage() {
  return (
    <EmptyState title="Страница не найдена" description="Проверьте адрес или вернитесь в каталог.">
      <Link className="btn btn--primary" to="/assistants">
        В каталог
      </Link>
    </EmptyState>
  );
}
