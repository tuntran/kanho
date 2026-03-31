import apiClient from "@/lib/api-interceptor";

export interface Workspace {
  id: string;
  name: string;
  slug: string;
  description?: string;
  logo_url?: string;
  accent_color?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export const workspaceKeys = {
  list: ["workspaces"] as const,
};

export async function listWorkspaces(): Promise<Workspace[]> {
  const { data } = await apiClient.get<Workspace[]>("/workspaces/");
  return data;
}
