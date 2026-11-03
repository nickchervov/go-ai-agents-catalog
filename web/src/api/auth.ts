import { http } from './http';
import type { Role, Token } from '@/types/api';

export const authApi = {
  async dummyLogin(role: Role): Promise<Token> {
    const { data } = await http.post<Token>('/dummyLogin', { role });
    return data;
  },
};
