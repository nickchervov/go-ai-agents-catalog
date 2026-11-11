// Типы, зеркалящие схемы из api.yaml.

export type Role = 'admin' | 'user';

export type RunStatus = 'pending' | 'success' | 'failed';

export type ErrorCode =
  | 'INVALID_REQUEST'
  | 'UNAUTHORIZED'
  | 'FORBIDDEN'
  | 'NOT_FOUND'
  | 'CATEGORY_NOT_FOUND'
  | 'ASSISTANT_NOT_FOUND'
  | 'ASSISTANT_INACTIVE'
  | 'LLM_PROVIDER_ERROR'
  | 'INTERNAL_ERROR';

export interface ErrorResponse {
  error: {
    code: ErrorCode;
    message: string;
  };
}

export interface User {
  id: string;
  email: string;
  role: Role;
  createdAt?: string | null;
}

export interface Token {
  token: string;
  user: User;
}

export interface Category {
  id: string;
  name: string;
  description?: string | null;
  createdAt?: string | null;
}

export interface CategoryCreateIn {
  name: string;
  description?: string | null;
}

export interface Assistant {
  id: string;
  categoryId: string;
  categoryName?: string | null;
  name: string;
  description: string;
  model: string;
  systemPrompt?: string | null;
  exampleUserPrompt?: string | null;
  isActive: boolean;
  createdAt?: string | null;
  updatedAt?: string | null;
}

export interface AssistantCreateIn {
  categoryId: string;
  name: string;
  description: string;
  model: string;
  systemPrompt: string;
  exampleUserPrompt?: string | null;
  isActive?: boolean;
}

export interface AssistantUpdateIn {
  categoryId: string;
  name: string;
  description: string;
  model: string;
  systemPrompt: string;
  exampleUserPrompt?: string | null;
  isActive: boolean;
}

export interface AssistantRun {
  id: string;
  assistantId: string;
  assistantName?: string | null;
  categoryId?: string | null;
  categoryName?: string | null;
  userId: string;
  model: string;
  userPrompt: string;
  output?: string | null;
  status: RunStatus;
  error?: string | null;
  createdAt?: string | null;
}

export interface AssistantRunCreateIn {
  userPrompt: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
}

export interface Paginated<T> {
  items: T[];
  pagination: Pagination;
}

export interface AssistantsResponse {
  assistants: Assistant[];
  pagination: Pagination;
}

export interface CategoriesResponse {
  categories: Category[];
}

export interface ModelsResponse {
  models: string[];
}

export interface RunsResponse {
  runs: AssistantRun[];
  pagination: Pagination;
}

export interface AssistantsQuery {
  categoryId?: string;
  q?: string;
  includeInactive?: boolean;
  page?: number;
  pageSize?: number;
}

export interface RunsQuery {
  status?: RunStatus;
  page?: number;
  pageSize?: number;
}

export interface AdminRunsQuery extends RunsQuery {
  assistantId?: string;
}
