import apiClient from "@/lib/api-interceptor";
import type { Board } from "@/types/board";

export interface Project {
  id: string;
  workspace_id: string;
  name: string;
  key: string;
  card_counter: number;
  created_at: string;
}

export const projectKeys = {
  list: (slug: string) => ["projects", slug] as const,
  boards: (slug: string, projectId: string) =>
    ["projects", slug, projectId, "boards"] as const,
};

export async function listProjects(
  workspaceSlug: string,
): Promise<Project[]> {
  const { data } = await apiClient.get<Project[]>(
    `/workspaces/${workspaceSlug}/projects`,
  );
  return data;
}

export async function createProject(
  workspaceSlug: string,
  input: { name: string; key: string },
): Promise<Project> {
  const { data } = await apiClient.post<Project>(
    `/workspaces/${workspaceSlug}/projects`,
    input,
  );
  return data;
}

export async function listProjectBoards(
  workspaceSlug: string,
  projectId: string,
): Promise<Board[]> {
  const { data } = await apiClient.get<Board[]>(
    `/workspaces/${workspaceSlug}/projects/${projectId}/boards`,
  );
  return data;
}
