import type { AdminRunsQuery, AssistantsQuery, RunsQuery } from '@/types/api';

/** Централизованные ключи запросов TanStack Query. */
export const queryKeys = {
  categories: ['categories'] as const,
  models: ['models'] as const,
  assistants: (query: AssistantsQuery) => ['assistants', query] as const,
  assistant: (id: string) => ['assistant', id] as const,
  myRuns: (query: RunsQuery) => ['runs', 'my', query] as const,
  allRuns: (query: AdminRunsQuery) => ['runs', 'all', query] as const,
};
