import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth';
import type { Role } from '@/types/api';

/** Пускает только авторизованных; иначе — на /login. */
export function RequireAuth() {
  const token = useAuthStore((s) => s.token);
  const location = useLocation();
  if (!token) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }
  return <Outlet />;
}

/** Пускает только указанные роли; иначе — на доступную страницу. */
export function RequireRole({ roles }: { roles: Role[] }) {
  const user = useAuthStore((s) => s.user);
  if (!user || !roles.includes(user.role)) {
    return <Navigate to="/assistants" replace />;
  }
  return <Outlet />;
}
