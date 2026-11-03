import { useEffect, useState } from 'react';
import { Link, NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import { useTheme } from '@/hooks/useTheme';
import { Button } from './Button';

function Icon({ path }: { path: string }) {
  return (
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d={path} />
    </svg>
  );
}

const ICONS = {
  catalog: 'M4 5h16M4 12h16M4 19h10',
  runs: 'M12 8v4l3 2M12 3a9 9 0 1 0 0 18 9 9 0 0 0 0-18Z',
  all: 'M4 6h16M4 12h16M4 18h16',
  plus: 'M12 5v14M5 12h14',
  tag: 'M20 12l-8 8-8-8V4h8l8 8ZM8 8h.01',
} as const;

export function Layout() {
  const user = useAuthStore((s) => s.user);
  const clear = useAuthStore((s) => s.clear);
  const navigate = useNavigate();
  const location = useLocation();
  const isAdmin = user?.role === 'admin';
  const isRuns = location.pathname.startsWith('/runs') || location.pathname.startsWith('/admin/runs');
  const { theme, toggle } = useTheme();

  const [menuOpen, setMenuOpen] = useState(false);

  // Закрываем мобильное меню при смене маршрута.
  useEffect(() => {
    setMenuOpen(false);
  }, [location.pathname]);

  const logout = () => {
    clear();
    navigate('/login', { replace: true });
  };

  return (
    <div className="app">
      <aside className={`sidebar${menuOpen ? ' open' : ''}`} id="app-nav">
        <div className="brand">
          <span className="brand__mark">
            <Icon path="M12 2v3M7 8h10a2 2 0 0 1 2 2v6a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2v-6a2 2 0 0 1 2-2ZM9 13h.01M15 13h.01" />
          </span>
          <span>
            <span className="brand__name">AI-каталог</span>
            <span className="brand__scope">assistants-catalog</span>
          </span>
        </div>

        <span className="nav-group">Каталог</span>
        <NavLink to="/assistants" className="nav-item">
          <Icon path={ICONS.catalog} />
          Ассистенты
        </NavLink>
        <NavLink to="/runs/my" className="nav-item">
          <Icon path={ICONS.runs} />
          Мои запуски
        </NavLink>

        {isAdmin && (
          <>
            <span className="nav-group">Администрирование</span>
            <NavLink to="/admin/runs" className="nav-item">
              <Icon path={ICONS.all} />
              Все запуски
            </NavLink>
            <NavLink to="/admin/assistants/new" className="nav-item">
              <Icon path={ICONS.plus} />
              Новый ассистент
            </NavLink>
            <NavLink to="/admin/categories/new" className="nav-item">
              <Icon path={ICONS.tag} />
              Новая категория
            </NavLink>
          </>
        )}

        <div className="sidebar__foot">
          <span className="nav-group" style={{ margin: '0 0 8px' }}>
            Служебное
          </span>
          <Button variant="ghost" size="sm" onClick={logout}>
            Выйти
          </Button>
        </div>
      </aside>

      {menuOpen && <div className="scrim" onClick={() => setMenuOpen(false)} />}

      <div className="main">
        <header className="topbar">
          <button
            type="button"
            className="menu-toggle"
            aria-label="Меню"
            aria-controls="app-nav"
            aria-expanded={menuOpen}
            onClick={() => setMenuOpen((v) => !v)}
          >
            <Icon path="M4 6h16M4 12h16M4 18h16" />
          </button>

          <nav className="crumbs" aria-label="Хлебные крошки">
            <Link to="/assistants" className="crumbs__root">
              Каталог
            </Link>
            <span className="crumbs__sep" aria-hidden="true">
              /
            </span>
            <span className="crumbs__current">{isRuns ? 'Запуски' : 'Ассистенты'}</span>
          </nav>

          <div className="topbar__account">
            {user && (
              <>
                <span className="topbar__email" title={user.email}>
                  {user.email}
                </span>
                <span className="pill pill--accent">{isAdmin ? 'Администратор' : 'Пользователь'}</span>
              </>
            )}            <button
              type="button"
              className="theme-toggle"
              onClick={toggle}
              aria-label={theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'}
              title={theme === 'dark' ? 'Светлая тема' : 'Тёмная тема'}
            >
              <Icon
                path={
                  theme === 'dark'
                    ? 'M12 3v2M12 19v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M3 12h2M19 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8Z'
                    : 'M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z'
                }
              />
            </button>
            <Button variant="ghost" size="sm" onClick={logout}>
              Выйти
            </Button>
          </div>
        </header>

        <div className="sheet">
          <main className="canvas">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}
