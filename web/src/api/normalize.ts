import type { Assistant, AssistantRun, Category, Pagination, RunStatus } from '@/types/api';

/**
 * Терпимая нормализация ответов API.
 *
 * Спецификация (api.yaml) — источник правды: каталог возвращает
 * `{ assistants, pagination }`, категории — `{ categories }`, запуски — `{ runs, pagination }`,
 * а поля в camelCase.
 *
 * Однако во время параллельной разработки backend может отдавать голый массив
 * и/или snake_case. Чтобы клиент не падал с белым экраном, здесь мы приводим
 * оба варианта к каноническому виду. Когда backend полностью соответствует
 * api.yaml, эти функции просто становятся тождественными и не мешают.
 */

type Raw = Record<string, unknown>;

function asRecord(value: unknown): Raw {
  return value && typeof value === 'object' ? (value as Raw) : {};
}

/** Берёт первое непустое значение по списку возможных ключей (camelCase и snake_case). */
function pick<T = unknown>(obj: Raw, ...keys: string[]): T | undefined {
  for (const key of keys) {
    const value = obj[key];
    if (value !== undefined && value !== null) return value as T;
  }
  return undefined;
}

function toArray(value: unknown): unknown[] {
  return Array.isArray(value) ? value : [];
}

export function normalizeAssistant(raw: unknown): Assistant {
  const o = asRecord(raw);
  return {
    id: String(pick(o, 'id') ?? ''),
    categoryId: String(pick(o, 'categoryId', 'category_id') ?? ''),
    categoryName: (pick<string>(o, 'categoryName', 'category_name') ?? null) as string | null,
    name: String(pick(o, 'name') ?? ''),
    description: String(pick(o, 'description') ?? ''),
    model: String(pick(o, 'model') ?? ''),
    systemPrompt: (pick<string>(o, 'systemPrompt', 'system_prompt') ?? null) as string | null,
    exampleUserPrompt: (pick<string>(
      o,
      'exampleUserPrompt',
      'example_user_prompt',
      'ex_user_prompt',
    ) ?? null) as string | null,
    isActive: Boolean(pick(o, 'isActive', 'is_active', 'is_activated') ?? false),
    createdAt: (pick<string>(o, 'createdAt', 'created_at') ?? null) as string | null,
    updatedAt: (pick<string>(o, 'updatedAt', 'updated_at') ?? null) as string | null,
  };
}

export function normalizeCategory(raw: unknown): Category {
  const o = asRecord(raw);
  return {
    id: String(pick(o, 'id') ?? ''),
    name: String(pick(o, 'name', 'category') ?? ''),
    description: (pick<string>(o, 'description') ?? null) as string | null,
    createdAt: (pick<string>(o, 'createdAt', 'created_at') ?? null) as string | null,
  };
}

export function normalizeRun(raw: unknown): AssistantRun {
  const o = asRecord(raw);
  return {
    id: String(pick(o, 'id') ?? ''),
    assistantId: String(pick(o, 'assistantId', 'assistant_id') ?? ''),
    assistantName: (pick<string>(o, 'assistantName', 'assistant_name') ?? null) as string | null,
    categoryId: (pick<string>(o, 'categoryId', 'category_id') ?? null) as string | null,
    categoryName: (pick<string>(o, 'categoryName', 'category_name') ?? null) as string | null,
    userId: String(pick(o, 'userId', 'user_id') ?? ''),
    model: String(pick(o, 'model') ?? ''),
    userPrompt: String(pick(o, 'userPrompt', 'user_prompt') ?? ''),
    output: (pick<string>(o, 'output', 'result') ?? null) as string | null,
    status: (pick<RunStatus>(o, 'status') ?? 'pending') as RunStatus,
    error: (pick<string>(o, 'error') ?? null) as string | null,
    createdAt: (pick<string>(o, 'createdAt', 'created_at') ?? null) as string | null,
  };
}

export function normalizePagination(raw: unknown, fallbackCount: number): Pagination {
  const o = asRecord(raw);
  return {
    page: Number(pick(o, 'page') ?? 1),
    pageSize: Number(pick(o, 'pageSize', 'page_size') ?? (fallbackCount || 10)),
    total: Number(pick(o, 'total') ?? fallbackCount),
  };
}

/**
 * Достаёт список и пагинацию из ответа:
 * поддержаны и `{ [key]: [...] , pagination }`, и голый массив.
 */
export function normalizeListResponse<T>(
  raw: unknown,
  key: string,
  mapItem: (item: unknown) => T,
): { items: T[]; pagination: Pagination } {
  const wrapped = asRecord(raw);
  const list = Array.isArray(raw) ? raw : toArray(wrapped[key]);
  const items = list.map(mapItem);
  return { items, pagination: normalizePagination(wrapped.pagination, items.length) };
}
