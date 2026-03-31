import { useNavigate } from "react-router-dom";
import type { Board } from "@/types/board";

interface BoardWithProject extends Board {
  projectKey: string;
}

interface WorkspaceBoardsGridProps {
  workspaceSlug: string;
  boards: BoardWithProject[];
  onCreateProject: () => void;
}

export function WorkspaceBoardsGrid({
  workspaceSlug,
  boards,
  onCreateProject,
}: WorkspaceBoardsGridProps) {
  const navigate = useNavigate();

  if (boards.length === 0) {
    return (
      <div className="mt-6">
        <h2 className="text-lg font-semibold text-text-1">Boards</h2>
        <div className="mt-4 flex flex-col items-center justify-center rounded-lg border border-dashed border-border py-12">
          <p className="text-sm text-text-2">Chua co board nao</p>
          <button
            onClick={onCreateProject}
            className="mt-3 rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:opacity-90"
          >
            Tao project dau tien
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="mt-6">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold text-text-1">Boards</h2>
        <button
          onClick={onCreateProject}
          className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-white hover:opacity-90"
        >
          + New
        </button>
      </div>
      <div className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        {boards.map((board) => (
          <button
            key={board.id}
            onClick={() => navigate(`/w/${workspaceSlug}/b/${board.id}`)}
            className="flex flex-col items-start gap-2 rounded-lg border border-border bg-surface-2 p-4 text-left hover:bg-surface-3 transition-colors"
          >
            <p className="font-medium text-text-1 line-clamp-1">
              {board.name}
            </p>
            <span className="inline-block rounded bg-surface-3 px-1.5 py-0.5 text-xs font-mono text-text-2">
              {board.projectKey}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
