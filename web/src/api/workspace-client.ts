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

export interface WorkspaceMember {
  workspace_id: string;
  user_id: string;
  role: "owner" | "admin" | "member";
  created_at: string;
  user_name?: string;
  user_email?: string;
}

export const workspaceKeys = {
  list: ["workspaces"] as const,
  detail: (slug: string) => ["workspaces", slug] as const,
  members: (slug: string) => ["workspaces", slug, "members"] as const,
};

export async function listWorkspaces(): Promise<Workspace[]> {
  const { data } = await apiClient.get<Workspace[]>("/workspaces/");
  return data;
}

export interface CreateWorkspaceInput {
  name: string;
  description?: string;
  accent_color?: string;
}

export async function createWorkspace(input: CreateWorkspaceInput): Promise<Workspace> {
  const { data } = await apiClient.post<Workspace>("/workspaces/", input);
  return data;
}

export async function getWorkspace(slug: string): Promise<Workspace> {
  const { data } = await apiClient.get<Workspace>(`/workspaces/${slug}/`);
  return data;
}

export async function listMembers(slug: string): Promise<WorkspaceMember[]> {
  const { data } = await apiClient.get<WorkspaceMember[]>(`/workspaces/${slug}/members`);
  return data;
}
