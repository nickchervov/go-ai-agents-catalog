import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { authApi, toApiError } from '@/api';
import { useAuthStore } from '@/store/auth';
import { Alert } from '@/components/States';
import { Button } from '@/components/Button';
import type { Role } from '@/types/api';

const ROLES: { value: Role; name: string; note: string }[] = [
  {
    value: 'user',
    name: 'Пользователь',
    note: 'Каталог, запуск ассистентов и своя история запусков.',
  },
  {
    value: 'admin',
    name: 'Администратор',
    note: 'Те же права, плюс категории, ассистенты и все запуски.',
  },
];

export function LoginPage() {
  const [role, setRole] = useState<Role>('user');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const setSession = useAuthStore((s) => s.setSession);
  const navigate = useNavigate();
  const location = useLocation();
  const from = (location.state as { from?: string } | null)?.from;

  const handleLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      const { token, user } = await authApi.dummyLogin(role);
      setSession(token, user);
      const fallback = from ?? '/assistants';
      navigate(fallback, { replace: true });
    } catch (err) {
      setError(toApiError(err).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login">
      <div className="login__panel">
        <div className="login__brand">
          <span className="brand__mark" aria-hidden="true">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor">
              <rect x="4" y="7" width="16" height="12" rx="3" />
              <path d="M12 3v4" />
              <circle cx="9" cy="13" r="1.2" fill="currentColor" stroke="none" />
              <circle cx="15" cy="13" r="1.2" fill="currentColor" stroke="none" />
            </svg>
          </span>
          <span className="brand__name">AI-каталог</span>
        </div>

        <h1>Вход в каталог</h1>
        <p className="sub">
          Тестовый вход. Выберите роль, приложение получит JWT через /dummyLogin и сохранит сессию.
        </p>

        {error && <Alert variant="error">{error}</Alert>}

        <div className="role-picker" aria-label="Роль пользователя">
          {ROLES.map((r) => (
            <button
              key={r.value}
              type="button"
              aria-pressed={role === r.value}
              className={`role-option ${role === r.value ? 'role-option--active' : ''}`}
              onClick={() => setRole(r.value)}
            >
              <span className="role-option__name">{r.name}</span>
              <span className="role-option__note">{r.note}</span>
            </button>
          ))}
        </div>

        <Button variant="primary" block onClick={handleLogin} disabled={loading}>
          {loading ? 'Входим…' : 'Войти'}
        </Button>
      </div>
    </div>
  );
}
