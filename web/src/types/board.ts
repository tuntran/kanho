export type Priority = "none" | "low" | "medium" | "high" | "urgent";
export type Visibility = "private" | "workspace";

export interface Board {
  id: string;
  project_id: string;
  name: string;
  slug: string;
  description?: string;
  visibility: Visibility;
  created_at: string;
  updated_at: string;
}

export interface Column {
  id: string;
  board_id: string;
  name: string;
  position: string;
  color?: string;
  wip_limit?: number | null;
  is_done: boolean;
  created_at: string;
}

export interface CardAssignee {
  card_id: string;
  user_id: string;
  created_at: string;
  user_name?: string;
}

export interface Label {
  id: string;
  workspace_id: string;
  name: string;
  color: string;
  created_at: string;
}

export interface Card {
  id: string;
  project_id: string;
  board_id: string;
  column_id: string;
  card_number: number;
  title: string;
  description?: string;
  priority: Priority;
  position: string;
  start_date?: string | null;
  due_date?: string | null;
  archived_at?: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  assignees?: CardAssignee[];
  labels?: Label[];
  readable_id?: string;
}

export interface BoardData {
  board: Board;
  columns: Column[];
}

export interface WSEvent {
  board_id: string;
  id?: string;
  action: string;
  target_type: string;
  target_id: string;
  data?: Record<string, unknown>;
}
