import { useCallback, useEffect, useState } from "react";
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
import { useParams } from "react-router-dom";
import { useBoardStore } from "@/stores/board-store";
import { moveCard as apiMoveCard } from "@/api/board-client";
import { getWorkspaceLabels } from "@/api/label-client";
import { BoardColumn } from "@/components/board/board-column";
import { CardDragOverlay } from "@/components/board/card-drag-overlay";
import { ColumnDragOverlay } from "@/components/board/column-drag-overlay";
import { ArchiveToggle } from "@/components/board/archive-toggle";
import { LabelChip } from "@/components/ui/label-chip";
import type { Card, Column, Label } from "@/types/board";

type ActiveDrag =
  | { type: "card"; card: Card }
  | { type: "column"; column: Column }
  | null;

export function BoardView() {
  const { workspaceSlug } = useParams<{ workspaceSlug: string }>();
  const columns = useBoardStore((s) => s.columns);
  const cards = useBoardStore((s) => s.cards);
  const cardsByColumn = useBoardStore((s) => s.cardsByColumn);
  const moveCardOptimistic = useBoardStore((s) => s.moveCardOptimistic);
  const confirmOp = useBoardStore((s) => s.confirmOp);
  const rollbackOp = useBoardStore((s) => s.rollbackOp);
  const activeFilters = useBoardStore((s) => s.activeFilters);
  const setLabelFilter = useBoardStore((s) => s.setLabelFilter);
  const toggleShowArchived = useBoardStore((s) => s.toggleShowArchived);
  const filteredCardsByColumn = useBoardStore((s) => s.filteredCardsByColumn);

  const [activeDrag, setActiveDrag] = useState<ActiveDrag>(null);
  const [wsLabels, setWsLabels] = useState<Label[]>([]);

  useEffect(() => {
    if (!workspaceSlug) return;
    getWorkspaceLabels(workspaceSlug).then(setWsLabels).catch(() => {});
  }, [workspaceSlug]);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
    useSensor(KeyboardSensor),
  );

  const columnIds = columns.map((c) => c.id);
  const filtered = filteredCardsByColumn();

  const toggleLabelId = (labelId: string) => {
    const current = activeFilters.labelIds;
    if (current.includes(labelId)) {
      setLabelFilter(current.filter((id) => id !== labelId));
    } else {
      setLabelFilter([...current, labelId]);
    }
  };

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

      const colCards = cardsByColumn[targetColumnId] ?? [];
      const idx = colCards.indexOf(cardId);

      const afterCardId = idx > 0 ? (colCards[idx - 1] ?? null) : null;
      const beforeCardId =
        idx < colCards.length - 1 ? (colCards[idx + 1] ?? null) : null;

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
    <div className="flex flex-1 flex-col overflow-hidden">
      {/* Filter bar */}
      <BoardFilterBar
        labels={wsLabels}
        activeLabelIds={activeFilters.labelIds}
        showArchived={activeFilters.showArchived}
        onToggleLabel={toggleLabelId}
        onToggleArchived={toggleShowArchived}
      />

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
                <BoardColumn
                  key={col.id}
                  column={col}
                  filteredCardIds={filtered[col.id]}
                />
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
    </div>
  );
}

function BoardFilterBar({
  labels,
  activeLabelIds,
  showArchived,
  onToggleLabel,
  onToggleArchived,
}: {
  labels: Label[];
  activeLabelIds: string[];
  showArchived: boolean;
  onToggleLabel: (id: string) => void;
  onToggleArchived: () => void;
}) {
  if (labels.length === 0) {
    return (
      <div className="flex items-center gap-2 border-b border-border px-5 py-2">
        <ArchiveToggle showArchived={showArchived} onToggle={onToggleArchived} />
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2 border-b border-border px-5 py-2">
      <ArchiveToggle showArchived={showArchived} onToggle={onToggleArchived} />
      <div className="mx-2 h-4 w-px bg-border" />
      <span className="text-xs text-text-2">Labels:</span>
      <div className="flex flex-wrap gap-1">
        {labels.map((label) => {
          const isActive = activeLabelIds.includes(label.id);
          return (
            <button
              key={label.id}
              onClick={() => onToggleLabel(label.id)}
              className={`transition-opacity ${isActive ? "" : "opacity-40 hover:opacity-70"}`}
            >
              <LabelChip name={label.name} color={label.color} />
            </button>
          );
        })}
      </div>
      {activeLabelIds.length > 0 && (
        <button
          onClick={() => {
            // Clear all label filters by toggling each active one off
            activeLabelIds.forEach((id) => onToggleLabel(id));
          }}
          className="text-xs text-text-2 hover:text-text-1"
          title="Clear label filter"
        >
          Clear
        </button>
      )}
    </div>
  );
}
