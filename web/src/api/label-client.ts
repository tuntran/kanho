import apiClient from "@/lib/api-interceptor";
import type { Label } from "@/types/board";

export async function getWorkspaceLabels(
  workspaceSlug: string,
): Promise<Label[]> {
  const { data } = await apiClient.get<Label[]>(
    `/workspaces/${workspaceSlug}/labels`,
  );
  return data;
}

export async function addCardLabel(
  cardId: string,
  labelId: string,
): Promise<void> {
  await apiClient.patch(`/cards/${cardId}/labels`, {
    add: [labelId],
    remove: [],
  });
}

export async function removeCardLabel(
  cardId: string,
  labelId: string,
): Promise<void> {
  await apiClient.patch(`/cards/${cardId}/labels`, {
    add: [],
    remove: [labelId],
  });
}

export async function searchCards(
  workspaceSlug: string,
  query: string,
): Promise<SearchResult[]> {
  const { data } = await apiClient.get<SearchResult[]>(
    `/workspaces/${workspaceSlug}/search`,
    { params: { q: query } },
  );
  return data;
}

export interface SearchResult {
  card_id: string;
  card_number: number;
  readable_id: string;
  title: string;
  board_id: string;
  board_name: string;
  column_name: string;
}
