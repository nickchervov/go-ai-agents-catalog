import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Layout } from '@/components/Layout';
import { RequireAuth, RequireRole } from '@/components/Guards';
import { LoginPage } from '@/pages/LoginPage';
import { AssistantsPage } from '@/pages/AssistantsPage';
import { AssistantDetailPage } from '@/pages/AssistantDetailPage';
import { MyRunsPage } from '@/pages/MyRunsPage';
import { AdminRunsPage } from '@/pages/AdminRunsPage';
import { NewCategoryPage } from '@/pages/NewCategoryPage';
import { NewAssistantPage } from '@/pages/NewAssistantPage';
import { EditAssistantPage } from '@/pages/EditAssistantPage';
import { NotFoundPage } from '@/pages/NotFoundPage';

export const router = createBrowserRouter(
  [
    { path: '/login', element: <LoginPage /> },
    {
      element: <RequireAuth />,
      children: [
        {
          element: <Layout />,
          children: [
            { path: '/', element: <Navigate to="/assistants" replace /> },
            { path: '/assistants', element: <AssistantsPage /> },
            { path: '/assistants/:id', element: <AssistantDetailPage /> },
            { path: '/runs/my', element: <MyRunsPage /> },
            {
              element: <RequireRole roles={['admin']} />,
              children: [
                { path: '/admin/runs', element: <AdminRunsPage /> },
                { path: '/admin/categories/new', element: <NewCategoryPage /> },
                { path: '/admin/assistants/new', element: <NewAssistantPage /> },
                { path: '/admin/assistants/:id/edit', element: <EditAssistantPage /> },
              ],
            },
            { path: '*', element: <NotFoundPage /> },
          ],
        },
      ],
    },
  ],
  {
    future: {
      v7_relativeSplatPath: true,
    },
  },
);
