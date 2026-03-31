import type { Workspace, WorkspaceMember } from "@/api/workspace-client";

interface WorkspaceDashboardHeaderProps {
  workspace: Workspace;
  members: WorkspaceMember[];
  boardCount: number;
}

export function WorkspaceDashboardHeader({
  workspace,
  members,
  boardCount,
}: WorkspaceDashboardHeaderProps) {
  const initial = workspace.name[0]?.toUpperCase() ?? "?";

  return (
    <div className="flex items-start justify-between">
      <div className="flex items-start gap-4">
        <div
          className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg text-xl font-bold text-white"
          style={{ backgroundColor: workspace.accent_color || "#6366f1" }}
        >
          {initial}
        </div>
        <div>
          <h1 className="text-xl font-bold text-text-1">{workspace.name}</h1>
          {workspace.description && (
            <p className="mt-0.5 text-sm text-text-2">{workspace.description}</p>
          )}
          <div className="mt-1.5 flex items-center gap-4 text-xs text-text-2">
            <span>{members.length} thanh vien</span>
            <span>{boardCount} board</span>
          </div>
        </div>
      </div>
      {/* TODO: workspace settings */}
      <button
        className="rounded-md p-2 text-text-2 hover:bg-surface-3 hover:text-text-1"
        title="Cai dat workspace"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
          />
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
          />
        </svg>
      </button>
    </div>
  );
}
