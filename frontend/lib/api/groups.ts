import { api } from './client';

// Tipos para grupos de mensagens
export interface MessageGroup {
  id: string;
  name: string;
  slug: string;
  description: string;
  isPublic: boolean;
  isArchived: boolean;
  settings?: {
    icebreakers?: string[];
    bannedWords?: string[];
    [key: string]: any;
  };
  createdAt: string;
}

export interface CreateGroupRequest {
  name: string;
  description: string;
  isPublic: boolean;
  settings?: {
    icebreakers?: string[];
    bannedWords?: string[];
    [key: string]: any;
  };
}

/**
 * Serviço para gerenciamento de grupos de mensagens
 */
export const groupsService = {
  /**
   * Obter todos os grupos do usuário
   */
  getGroups: async (includeArchived = false) => {
    return api.get<{ groups: MessageGroup[] }>(`/groups?includeArchived=${includeArchived}`);
  },

  /**
   * Criar um novo grupo
   */
  createGroup: async (data: CreateGroupRequest) => {
    return api.post<MessageGroup>('/groups', data);
  },

  /**
   * Obter detalhes de um grupo
   */
  getGroup: async (id: string) => {
    return api.get<MessageGroup>(`/groups/${id}`);
  },

  /**
   * Atualizar um grupo
   */
  updateGroup: async (id: string, data: CreateGroupRequest) => {
    return api.put<{ message: string }>(`/groups/${id}`, data);
  },

  /**
   * Arquivar um grupo
   */
  archiveGroup: async (id: string) => {
    return api.delete<{ message: string }>(`/groups/${id}`);
  },

  /**
   * Desarquivar um grupo
   */
  unarchiveGroup: async (id: string) => {
    return api.post<{ message: string }>(`/groups/${id}/unarchive`, {});
  },

  /**
   * Gerar link para compartilhamento
   */
  createSharedAccess: async (groupId: string, email: string, expiresAt?: Date) => {
    return api.post<{ id: string; token: string }>(`/groups/${groupId}/share`, {
      email,
      expiresAt: expiresAt?.toISOString(),
    });
  },

  /**
   * Obter acessos compartilhados
   */
  getSharedAccess: async (groupId: string) => {
    return api.get<{ sharedAccess: Array<{ id: string; email: string; token: string; isActive: boolean; expiresAt?: string; createdAt: string }> }>(`/groups/${groupId}/shared`);
  },

  /**
   * Revogar acesso compartilhado
   */
  revokeSharedAccess: async (groupId: string, shareId: string) => {
    return api.delete<{ message: string }>(`/groups/${groupId}/share/${shareId}`);
  },
};
