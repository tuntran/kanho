import { useCallback, useEffect, useState } from "react";
import { getBoard, getCards } from "@/api/board-client";
import { useBoardStore } from "@/stores/board-store";

export function useBoardData(boardId: string | undefined) {
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const initBoard = useBoardStore((s) => s.initBoard);

  const fetchBoard = useCallback(async () => {
    if (!boardId) return;
    try {
      setIsLoading(true);
      setError(null);
      const [boardData, cards] = await Promise.all([
        getBoard(boardId),
        getCards(boardId),
      ]);
      initBoard(boardData.board, boardData.columns, cards);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load board");
    } finally {
      setIsLoading(false);
    }
  }, [boardId, initBoard]);

  useEffect(() => {
    fetchBoard();
  }, [fetchBoard]);

  return { isLoading, error, refetch: fetchBoard };
}
