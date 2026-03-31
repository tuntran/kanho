import { useEffect } from "react";
import { useParams, Outlet } from "react-router-dom";
import { useBoardData } from "@/hooks/use-board-data";
import { useBoardWs } from "@/hooks/use-board-ws";
import { useBoardStore } from "@/stores/board-store";
import { BoardHeader } from "@/components/board/board-header";
import { BoardView } from "@/components/board/board-view";

export function BoardPage() {
  const { boardId } = useParams<{ boardId: string }>();
  const { isLoading, error, refetch } = useBoardData(boardId);
  const reset = useBoardStore((s) => s.reset);

  useBoardWs(boardId ?? null, refetch);

  // Cleanup on unmount
  useEffect(() => {
    return () => reset();
  }, [reset]);

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-bg">
        <div className="flex flex-col items-center gap-3">
          <svg
            className="h-8 w-8 animate-spin text-primary"
            viewBox="0 0 24 24"
            fill="none"
          >
            <circle
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="3"
              strokeLinecap="round"
              className="opacity-25"
            />
            <path
              d="M4 12a8 8 0 018-8"
              stroke="currentColor"
              strokeWidth="3"
              strokeLinecap="round"
            />
          </svg>
          <p className="text-sm text-text-2">Loading board...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-screen items-center justify-center bg-bg">
        <div className="text-center">
          <p className="text-sm text-danger">{error}</p>
          <button
            onClick={refetch}
            className="mt-3 rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:bg-primary-hover"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-screen flex-col bg-bg">
      <BoardHeader />
      <BoardView />
      <Outlet />
    </div>
  );
}
