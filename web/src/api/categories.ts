import { http } from './http';
import { normalizeCategory, normalizeListResponse } from './normalize';
import type { Category, CategoryCreateIn } from '@/types/api';

export const categoriesApi = {
  async list(): Promise<Category[]> {
    const { data } = await http.get<unknown>('/categories');
    const { items } = normalizeListResponse(data, 'categories', normalizeCategory);
    return items;
  },

  async create(payload: CategoryCreateIn): Promise<Category> {
    const { data } = await http.post<unknown>('/categories', payload);
    return normalizeCategory(data);
  },
};
