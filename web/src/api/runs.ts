import { http } from './http';
import { normalizeListResponse, normalizeRun } from './normalize';
import type { AdminRunsQuery, RunsQuery, RunsResponse } from '@/types/api';

export const runsApi = {
  async my(query: RunsQuery = {}): Promise<RunsResponse> {
    const { data } = await http.get<unknown>('/runs/my', { params: cleanParams(query) });
    const { items, pagination } = normalizeListResponse(data, 'runs', normalizeRun);
    return { runs: items, pagination };
  },

  async all(query: AdminRunsQuery = {}): Promise<RunsResponse> {
    const { data } = await http.get<unknown>('/admin/runs', { params: cleanParams(query) });
    const { items, pagination } = normalizeListResponse(data, 'runs', normalizeRun);
    return { runs: items, pagination };
  },
};

function cleanParams(params: object): Record<string, unknown> {
  return Object.fromEntries(
    Object.entries(params).filter(
      ([, value]) => value !== undefined && value !== '' && value !== null,
    ),
  );
}
