import apiClient from "@/lib/api-interceptor";
import type { BoardData, Card, Column } from "@/types/board";

export async function getBoard(boardId: string): Promise<BoardData> {
  const { data } = await apiClient.get<BoardData>(`/boards/${boardId}`);
  return data;
}

export async function getCards(boardId: string): Promise<Card[]> {
  const { data } = await apiClient.get<Card[]>(`/boards/${boardId}/cards`);
  return data;
}

export async function createCard(
  boardId: string,
  input: { column_id: string; title: string; priority?: string },
): Promise<Card> {
  const { data } = await apiClient.post<Card>(`/boards/${boardId}/cards`, input);
  return data;
}

export async function updateCard(
  cardId: string,
  input: { title?: string; description?: string; priority?: string },
): Promise<Card> {
  const { data } = await apiClient.patch<Card>(`/cards/${cardId}`, input);
  return data;
}

export async function moveCard(
  cardId: string,
  input: { column_id: string; after_card_id?: string; before_card_id?: string },
): Promise<void> {
  await apiClient.post(`/cards/${cardId}/move`, input);
}

export async function createColumn(
  boardId: string,
  input: { name: string; color?: string; wip_limit?: number; is_done?: boolean },
): Promise<Column> {
  const { data } = await apiClient.post<Column>(`/boards/${boardId}/columns`, input);
  return data;
}

export async function reorderColumns(
  items: Array<{ id: string; position: string }>,
): Promise<void> {
  await apiClient.patch("/columns/reorder", items);
}
