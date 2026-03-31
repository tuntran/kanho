import apiClient from "@/lib/api-interceptor";
import type { Card } from "@/types/board";

export interface CardDetail extends Card {
  description?: string;
}

export async function getCard(cardId: string): Promise<CardDetail> {
  const { data } = await apiClient.get<CardDetail>(`/cards/${cardId}`);
  return data;
}

export async function updateCard(
  cardId: string,
  input: { title?: string; description?: string; priority?: string },
): Promise<CardDetail> {
  const { data } = await apiClient.patch<CardDetail>(`/cards/${cardId}`, input);
  return data;
}

export interface PresignResponse {
  upload_url: string;
  file_key: string;
  public_url: string;
}

export async function getUploadUrl(input: {
  filename: string;
  content_type: string;
  card_id?: string;
}): Promise<PresignResponse> {
  const { data } = await apiClient.post<PresignResponse>(
    "/attachments/upload-url",
    input,
  );
  return data;
}

export async function uploadFileToPresignedUrl(
  url: string,
  file: File,
): Promise<void> {
  await fetch(url, {
    method: "PUT",
    body: file,
    headers: { "Content-Type": file.type },
  });
}

export interface Activity {
  id: string;
  card_id: string;
  user_id: string;
  action: string;
  field?: string;
  old_value?: string;
  new_value?: string;
  created_at: string;
  user_name?: string;
}

export async function getActivities(
  cardId: string,
  limit = 50,
  offset = 0,
): Promise<Activity[]> {
  const { data } = await apiClient.get<Activity[]>(
    `/cards/${cardId}/activity`,
    { params: { limit, offset } },
  );
  return data;
}
