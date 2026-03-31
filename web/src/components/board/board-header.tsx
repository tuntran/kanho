import { useBoardStore } from "@/stores/board-store";
import { LiveDot } from "@/components/board/live-dot";

export function BoardHeader() {
  const board = useBoardStore((s) => s.board);
  if (!board) return null;

  return (
    <header className="flex items-center gap-3 border-b border-border bg-bg px-5 py-3">
      <h1 className="text-lg font-bold text-text-1">{board.name}</h1>
      <LiveDot />
      {board.description && (
        <span className="text-xs text-text-2">{board.description}</span>
      )}
    </header>
  );
}
