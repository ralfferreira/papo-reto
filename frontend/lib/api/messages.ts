import { api } from './client';

// Tipos para mensagens
export interface Message {
  id: string;
  content: string;
  isRead: boolean;
  isFavorite: boolean;
  isRevealed: boolean;
  senderID?: string;
  createdAt: string;
}

export interface SendMessageRequest {
  content: string;
  senderId?: string;
  revealName?: boolean;
}

/**
 * Serviço para gerenciamento de mensagens
 */
export const messagesService = {
  /**
   * Obter mensagens de um grupo
   */
  getMessages: async (groupId: string, page = 1, pageSize = 20) => {
    return api.get<{ messages: Message[] }>(`/groups/${groupId}/messages?page=${page}&pageSize=${pageSize}`);
  },

  /**
   * Atualizar status de uma mensagem
   */
  updateMessage: async (messageId: string, data: { isRead?: boolean; isFavorite?: boolean }) => {
    return api.put<{ message: string }>(`/messages/${messageId}`, data);
  },

  /**
   * Marcar mensagem como lida
   */
  markAsRead: async (messageId: string) => {
    return messagesService.updateMessage(messageId, { isRead: true });
  },

  /**
   * Alternar status de favorito
   */
  toggleFavorite: async (messageId: string, isFavorite: boolean) => {
    return messagesService.updateMessage(messageId, { isFavorite });
  },

  /**
   * Excluir uma mensagem
   */
  deleteMessage: async (messageId: string) => {
    return api.delete<{ message: string }>(`/messages/${messageId}`);
  },

  /**
   * Enviar mensagem anônima (endpoint público)
   */
  sendAnonymousMessage: async (slug: string, data: SendMessageRequest) => {
    return api.post<{ message: string }>(`/public/send/${slug}`, data, false);
  },
};
