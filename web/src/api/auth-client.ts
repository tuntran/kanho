import axios from "axios";

const api = axios.create({
  baseURL: "/api/v1",
  withCredentials: true,
  headers: { "Content-Type": "application/json" },
});

export interface User {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResponse {
  user: User;
  accessToken: string;
}

interface SetupStatusResponse {
  requiresSetup: boolean;
}

export async function register(
  email: string,
  password: string,
  name: string,
): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>("/auth/register", {
    email,
    password,
    name,
  });
  return data;
}

export async function login(
  email: string,
  password: string,
): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>("/auth/login", {
    email,
    password,
  });
  return data;
}

export async function refresh(): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>("/auth/refresh");
  return data;
}

export async function logout(): Promise<void> {
  await api.post("/auth/logout");
}

export async function getMe(accessToken: string): Promise<User> {
  const { data } = await api.get<User>("/auth/me", {
    headers: { Authorization: `Bearer ${accessToken}` },
  });
  return data;
}

export async function getSetupStatus(): Promise<boolean> {
  const { data } = await api.get<SetupStatusResponse>("/setup/status");
  return data.requiresSetup;
}
