import apiClient from "@/lib/api-interceptor";

export interface Comment {
  id: string;
  card_id: string;
  author_id: string;
  content: string;
  created_at: string;
  updated_at: string;
  author_name?: string;
}

export async function getComments(
  cardId: string,
  limit = 50,
  offset = 0,
): Promise<Comment[]> {
  const { data } = await apiClient.get<Comment[]>(
    `/cards/${cardId}/comments`,
    { params: { limit, offset } },
  );
  return data;
}

export async function createComment(
  cardId: string,
  content: string,
): Promise<Comment> {
  const { data } = await apiClient.post<Comment>(
    `/cards/${cardId}/comments`,
    { content },
  );
  return data;
}

export async function updateComment(
  commentId: string,
  content: string,
): Promise<Comment> {
  const { data } = await apiClient.patch<Comment>(
    `/comments/${commentId}`,
    { content },
  );
  return data;
}

export async function deleteComment(commentId: string): Promise<void> {
  await apiClient.delete(`/comments/${commentId}`);
}
