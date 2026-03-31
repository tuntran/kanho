import { useCallback, useState } from "react";
import {
  DndContext,
  DragOverlay,
  closestCorners,
  PointerSensor,
  KeyboardSensor,
  useSensor,
  useSensors,
  type DragStartEvent,
  type DragEndEvent,
  type DragOverEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  horizontalListSortingStrategy,
} from "@dnd-kit/sortable";
import { useBoardStore } from "@/stores/board-store";
import { moveCard as apiMoveCard } from "@/api/board-client";
import { BoardColumn } from "@/components/board/board-column";
import { CardDragOverlay } from "@/components/board/card-drag-overlay";
import { ColumnDragOverlay } from "@/components/board/column-drag-overlay";
import type { Card, Column } from "@/types/board";

type ActiveDrag =
  | { type: "card"; card: Card }
  | { type: "column"; column: Column }
  | null;

export function BoardView() {
  const columns = useBoardStore((s) => s.columns);
  const cards = useBoardStore((s) => s.cards);
  const cardsByColumn = useBoardStore((s) => s.cardsByColumn);
  const moveCardOptimistic = useBoardStore((s) => s.moveCardOptimistic);
  const confirmOp = useBoardStore((s) => s.confirmOp);
  const rollbackOp = useBoardStore((s) => s.rollbackOp);

  const [activeDrag, setActiveDrag] = useState<ActiveDrag>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(KeyboardSensor),
  );

  const columnIds = columns.map((c) => c.id);

  const handleDragStart = useCallback(
    (event: DragStartEvent) => {
      const { active } = event;
      const data = active.data.current;
      if (data?.type === "card") {
        setActiveDrag({ type: "card", card: data.card as Card });
      } else if (data?.type === "column") {
        setActiveDrag({ type: "column", column: data.column as Column });
      }
    },
    [],
  );

  const handleDragOver = useCallback(
    (event: DragOverEvent) => {
      const { active, over } = event;
      if (!over) return;

      const activeData = active.data.current;
      if (activeData?.type !== "card") return;

      const activeCard = cards[active.id as string];
      if (!activeCard) return;

      // Determine the target column
      let targetColumnId: string | null = null;
      const overData = over.data.current;

      if (overData?.type === "card") {
        const overCard = cards[over.id as string];
        if (overCard) targetColumnId = overCard.column_id;
      } else if (overData?.type === "column-drop") {
        targetColumnId = overData.columnId as string;
      } else if (overData?.type === "column") {
        targetColumnId = (overData.column as Column).id;
      }

      if (!targetColumnId || targetColumnId === activeCard.column_id) return;

      // Move card to new column optimistically (at the end)
      moveCardOptimistic(active.id as string, targetColumnId, null, null);
    },
    [cards, moveCardOptimistic],
  );

  const handleDragEnd = useCallback(
    async (event: DragEndEvent) => {
      setActiveDrag(null);
      const { active, over } = event;
      if (!over) return;

      const activeData = active.data.current;
      if (activeData?.type !== "card") return;

      const cardId = active.id as string;
      const activeCard = cards[cardId];
      if (!activeCard) return;

      // Determine final target column
      let targetColumnId = activeCard.column_id;
      const overData = over.data.current;

      if (overData?.type === "card") {
        const overCard = cards[over.id as string];
        if (overCard) targetColumnId = overCard.column_id;
      } else if (overData?.type === "column-drop") {
        targetColumnId = overData.columnId as string;
      } else if (overData?.type === "column") {
        targetColumnId = (overData.column as Column).id;
      }

      // Find where the card ended up in the sorted list
      const colCards = cardsByColumn[targetColumnId] ?? [];
      const idx = colCards.indexOf(cardId);

      const afterCardId = idx > 0 ? (colCards[idx - 1] ?? null) : null;
      const beforeCardId =
        idx < colCards.length - 1 ? (colCards[idx + 1] ?? null) : null;

      // If card didn't actually move, skip
      if (active.id === over.id && !afterCardId && !beforeCardId) return;

      const opId = moveCardOptimistic(
        cardId,
        targetColumnId,
        afterCardId,
        beforeCardId,
      );
      if (!opId) return;

      try {
        const moveInput: {
          column_id: string;
          after_card_id?: string;
          before_card_id?: string;
        } = { column_id: targetColumnId };
        if (afterCardId) moveInput.after_card_id = afterCardId;
        if (beforeCardId) moveInput.before_card_id = beforeCardId;
        await apiMoveCard(cardId, moveInput);
        confirmOp(opId);
      } catch {
        rollbackOp(opId);
      }
    },
    [cards, cardsByColumn, moveCardOptimistic, confirmOp, rollbackOp],
  );

  return (
    <div className="flex-1 overflow-x-auto overflow-y-hidden p-5">
      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragOver={handleDragOver}
        onDragEnd={handleDragEnd}
      >
        <SortableContext
          items={columnIds}
          strategy={horizontalListSortingStrategy}
        >
          <div className="flex gap-4">
            {columns.map((col) => (
              <BoardColumn key={col.id} column={col} />
            ))}
          </div>
        </SortableContext>

        <DragOverlay dropAnimation={null}>
          {activeDrag?.type === "card" && (
            <CardDragOverlay card={activeDrag.card} />
          )}
          {activeDrag?.type === "column" && (
            <ColumnDragOverlay column={activeDrag.column} />
          )}
        </DragOverlay>
      </DndContext>
    </div>
  );
}
