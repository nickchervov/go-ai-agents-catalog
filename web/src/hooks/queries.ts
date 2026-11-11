import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { assistantsApi, categoriesApi, modelsApi, runsApi } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import type {
  AdminRunsQuery,
  AssistantCreateIn,
  AssistantRunCreateIn,
  AssistantUpdateIn,
  AssistantsQuery,
  CategoryCreateIn,
  RunsQuery,
} from '@/types/api';

export function useCategories() {
  return useQuery({
    queryKey: queryKeys.categories,
    queryFn: () => categoriesApi.list(),
    staleTime: 60_000,
  });
}

export function useModels() {
  return useQuery({
    queryKey: queryKeys.models,
    queryFn: () => modelsApi.list(),
    staleTime: 60_000,
  });
}

export function useCreateCategory() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: CategoryCreateIn) => categoriesApi.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: queryKeys.categories }),
  });
}

export function useAssistants(query: AssistantsQuery) {
  return useQuery({
    queryKey: queryKeys.assistants(query),
    queryFn: () => assistantsApi.list(query),
    placeholderData: (prev) => prev,
  });
}

export function useAssistant(id: string | undefined) {
  return useQuery({
    queryKey: queryKeys.assistant(id ?? ''),
    queryFn: () => assistantsApi.get(id as string),
    enabled: Boolean(id),
  });
}

export function useCreateAssistant() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: AssistantCreateIn) => assistantsApi.create(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['assistants'] }),
  });
}

export function useUpdateAssistant() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: AssistantUpdateIn }) =>
      assistantsApi.update(id, payload),
    onSuccess: (_data, variables) => {
      qc.invalidateQueries({ queryKey: ['assistants'] });
      qc.invalidateQueries({ queryKey: queryKeys.assistant(variables.id) });
    },
  });
}

export function useRunAssistant() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: AssistantRunCreateIn }) =>
      assistantsApi.run(id, payload),
    // После запуска обновляем историю запусков.
    onSuccess: () => qc.invalidateQueries({ queryKey: ['runs'] }),
  });
}

export function useMyRuns(query: RunsQuery) {
  return useQuery({
    queryKey: queryKeys.myRuns(query),
    queryFn: () => runsApi.my(query),
    placeholderData: (prev) => prev,
  });
}

export function useAllRuns(query: AdminRunsQuery) {
  return useQuery({
    queryKey: queryKeys.allRuns(query),
    queryFn: () => runsApi.all(query),
    placeholderData: (prev) => prev,
  });
}
