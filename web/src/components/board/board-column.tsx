import {
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { useDroppable } from "@dnd-kit/core";
import { CSS } from "@dnd-kit/utilities";
import type { Column } from "@/types/board";
import { useBoardStore } from "@/stores/board-store";
import { ColumnHeader } from "@/components/board/column-header";
import { KanbanCard } from "@/components/board/kanban-card";
import { AddCardForm } from "@/components/board/add-card-form";

interface BoardColumnProps {
  column: Column;
  filteredCardIds?: string[];
}

export function BoardColumn({ column, filteredCardIds }: BoardColumnProps) {
  const cards = useBoardStore((s) => s.cards);
  const archivedCards = useBoardStore((s) => s.archivedCards);
  const storeCardIds = useBoardStore((s) => s.cardsByColumn[column.id] ?? []);
  const cardIds = filteredCardIds ?? storeCardIds;
  const boardId = useBoardStore((s) => s.boardId);

  const {
    attributes,
    listeners,
    setNodeRef: setSortableRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: column.id,
    data: { type: "column", column },
  });

  const { setNodeRef: setDroppableRef } = useDroppable({
    id: `column-drop-${column.id}`,
    data: { type: "column-drop", columnId: column.id },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setSortableRef}
      style={style}
      className="flex w-[280px] min-w-[280px] flex-col rounded-panel border border-border bg-surface/60 p-3"
    >
      {/* Drag handle = column header area */}
      <div {...attributes} {...listeners} className="cursor-grab active:cursor-grabbing">
        <ColumnHeader column={column} cardCount={cardIds.length} />
      </div>

      {/* Card list — droppable zone */}
      <div
        ref={setDroppableRef}
        className="flex min-h-[40px] flex-1 flex-col gap-2"
      >
        <SortableContext
          items={cardIds}
          strategy={verticalListSortingStrategy}
        >
          {cardIds.map((id) => {
            const card = cards[id] ?? archivedCards[id];
            if (!card) return null;
            const isArchived = !!card.archived_at;
            return (
              <div
                key={id}
                className={isArchived ? "opacity-50" : ""}
              >
                <KanbanCard card={card} isArchived={isArchived} />
              </div>
            );
          })}
        </SortableContext>
      </div>

      {boardId && <AddCardForm boardId={boardId} columnId={column.id} />}
    </div>
  );
}
