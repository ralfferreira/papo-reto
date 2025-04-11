import { api } from './client';

// Função segura para acessar localStorage (evita erros no SSR)
const getLocalStorage = () => {
  if (typeof window !== 'undefined') {
    return window.localStorage;
  }
  return null;
};

// Tipos para autenticação
export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token?: string;
  id?: string;
  email?: string;
  name?: string;
  isVerified?: boolean;
  plan?: string;
  messageCount?: number;
  activeGroups?: number;
  createdAt?: string;
}

export interface NotificationSettings {
  [key: string]: unknown;
}

/**
 * Serviço de autenticação
 */
export const authService = {
  /**
   * Registrar um novo usuário
   */
  register: async (data: RegisterRequest) => {
    return api.post<AuthResponse>('/auth/register', data, false);
  },

  /**
   * Fazer login
   */
  login: async (data: LoginRequest) => {
    const response = await api.post<{ token: string }>('/auth/login', data, false);
    
    // Se o login for bem-sucedido, salvar o token
    if (response.data?.token) {
      const localStorage = getLocalStorage();
      localStorage?.setItem('token', response.data.token);
    }
    
    return response;
  },

  /**
   * Obter perfil do usuário
   */
  getProfile: async () => {
    return api.get<AuthResponse>('/user/profile');
  },

  /**
   * Atualizar perfil do usuário
   */
  updateProfile: async (data: { name: string; avatarURL?: string }) => {
    return api.put<{ message: string }>('/user/profile', data);
  },

  /**
   * Atualizar senha do usuário
   */
  updatePassword: async (data: { currentPassword: string; newPassword: string }) => {
    return api.put<{ message: string }>('/user/password', data);
  },

  /**
   * Atualizar configurações de notificação
   */
  updateNotifications: async (settings: NotificationSettings) => {
    return api.put<{ message: string }>('/user/notifications', settings);
  },

  /**
   * Fazer logout
   */
  logout: () => {
    const localStorage = getLocalStorage();
    localStorage?.removeItem('token');
  },

  /**
   * Verificar se o usuário está autenticado
   */
  isAuthenticated: () => {
    const localStorage = getLocalStorage();
    return !!localStorage?.getItem('token');
  },
};