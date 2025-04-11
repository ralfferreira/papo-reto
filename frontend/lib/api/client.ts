/**
 * Cliente HTTP para comunicação com a API
 */

// URL base da API
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

// Tipos de resposta da API
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  status: number;
}

// Opções para requisições
interface RequestOptions {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  headers?: Record<string, string>;
  body?: unknown;
  requiresAuth?: boolean;
}

// Função segura para acessar localStorage (evita erros no SSR)
const getLocalStorage = () => {
  if (typeof window !== 'undefined') {
    return window.localStorage;
  }
  return null;
};

/**
 * Função para fazer requisições HTTP para a API
 */
export async function apiRequest<T>(
  endpoint: string,
  options: RequestOptions
): Promise<ApiResponse<T>> {
  const { method, body, requiresAuth = true } = options;

  // Configurar headers
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Adicionar token de autenticação se necessário
  if (requiresAuth) {
    const localStorage = getLocalStorage();
    const token = localStorage?.getItem('token');
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  try {
    // Configuração da requisição
    const requestConfig: RequestInit = {
      method,
      headers,
      credentials: 'include',
      mode: 'cors',
    };

    // Adicionar body apenas se não for GET
    if (body && method !== 'GET') {
      requestConfig.body = JSON.stringify(body);
    }

    // Fazer a requisição
    const response = await fetch(`${API_BASE_URL}${endpoint}`, requestConfig);

    // Tentar obter dados da resposta como JSON
    let data;
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      data = await response.json().catch(() => ({}));
    } else {
      data = {};
    }

    // Retornar resposta formatada
    return {
      data: response.ok ? data : undefined,
      error: response.ok ? undefined : data.error || `Erro ${response.status}: ${response.statusText}`,
      status: response.status,
    };
  } catch (error) {
    // Tratar erros de rede
    console.error('API request error:', error);
    return {
      error: error instanceof Error ? error.message : 'Erro de rede',
      status: 0,
    };
  }
}

/**
 * Funções auxiliares para os diferentes métodos HTTP
 */
export const api = {
  get: <T>(endpoint: string, requiresAuth = true) =>
    apiRequest<T>(endpoint, { method: 'GET', requiresAuth }),

  post: <T>(endpoint: string, body: unknown, requiresAuth = true) =>
    apiRequest<T>(endpoint, { method: 'POST', body, requiresAuth }),

  put: <T>(endpoint: string, body: unknown, requiresAuth = true) =>
    apiRequest<T>(endpoint, { method: 'PUT', body, requiresAuth }),

  delete: <T>(endpoint: string, requiresAuth = true) =>
    apiRequest<T>(endpoint, { method: 'DELETE', requiresAuth }),
};