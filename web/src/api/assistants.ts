import { http } from './http';
import { normalizeAssistant, normalizeListResponse, normalizeRun } from './normalize';
import type {
  Assistant,
  AssistantCreateIn,
  AssistantRun,
  AssistantRunCreateIn,
  AssistantsQuery,
  AssistantsResponse,
  AssistantUpdateIn,
} from '@/types/api';

export const assistantsApi = {
  async list(query: AssistantsQuery = {}): Promise<AssistantsResponse> {
    const { data } = await http.get<unknown>('/assistants', {
      params: cleanParams(query),
    });
    const { items, pagination } = normalizeListResponse(data, 'assistants', normalizeAssistant);
    return { assistants: items, pagination };
  },

  async get(id: string): Promise<Assistant> {
    const { data } = await http.get<unknown>(`/assistants/${id}`);
    return normalizeAssistant(data);
  },

  async create(payload: AssistantCreateIn): Promise<Assistant> {
    const { data } = await http.post<unknown>('/assistants', payload);
    return normalizeAssistant(data);
  },

  async update(id: string, payload: AssistantUpdateIn): Promise<Assistant> {
    const { data } = await http.put<unknown>(`/assistants/${id}`, payload);
    return normalizeAssistant(data);
  },

  async run(id: string, payload: AssistantRunCreateIn): Promise<AssistantRun> {
    const { data } = await http.post<unknown>(`/assistants/${id}/run`, payload);
    return normalizeRun(data);
  },
};

/** Убирает пустые значения, чтобы не отправлять `?q=` и подобное. */
function cleanParams(params: object): Record<string, unknown> {
  return Object.fromEntries(
    Object.entries(params).filter(
      ([, value]) => value !== undefined && value !== '' && value !== null,
    ),
  );
}
