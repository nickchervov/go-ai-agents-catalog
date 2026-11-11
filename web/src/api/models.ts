import { http } from './http';
import type { ModelsResponse } from '@/types/api';

export const modelsApi = {
  /** Список моделей, доступных у LLM-провайдера (только для администратора). */
  async list(): Promise<string[]> {
    const { data } = await http.get<ModelsResponse>('/models');
    return Array.isArray(data?.models) ? data.models : [];
  },
};
