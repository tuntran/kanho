import { useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import type { Card } from "@/types/board";
import { PriorityBadge } from "@/components/ui/priority-badge";
import { LabelChip } from "@/components/ui/label-chip";

interface KanbanCardProps {
  card: Card;
}

export function KanbanCard({ card }: KanbanCardProps) {
  const navigate = useNavigate();
  const { workspaceSlug, boardId } = useParams<{
    workspaceSlug: string;
    boardId: string;
  }>();
  const pointerDownPos = useRef<{ x: number; y: number } | null>(null);

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: card.id,
    data: { type: "card", card },
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const dueStr = card.due_date
    ? new Date(card.due_date).toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
      })
    : null;

  const handlePointerDown = (e: React.PointerEvent) => {
    pointerDownPos.current = { x: e.clientX, y: e.clientY };
    listeners?.onPointerDown?.(e as never);
  };

  const handleClick = (e: React.MouseEvent) => {
    // Only navigate if pointer didn't move (not a drag)
    if (pointerDownPos.current) {
      const dx = Math.abs(e.clientX - pointerDownPos.current.x);
      const dy = Math.abs(e.clientY - pointerDownPos.current.y);
      if (dx < 5 && dy < 5) {
        const cardRef = card.readable_id ?? String(card.card_number);
        navigate(`/w/${workspaceSlug}/b/${boardId}/cards/${cardRef}`);
      }
    }
    pointerDownPos.current = null;
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onPointerDown={handlePointerDown}
      onClick={handleClick}
      className="group cursor-grab rounded-card border border-border bg-surface p-3 shadow-card transition-all duration-150 hover:-translate-y-0.5 hover:border-border-focus hover:shadow-card-hover active:cursor-grabbing"
    >
      {/* Top row: priority badge + readable ID */}
      <div className="flex items-center justify-between gap-2">
        <PriorityBadge priority={card.priority} />
        {card.readable_id && (
          <span className="font-mono text-xs font-medium text-text-2">
            {card.readable_id}
          </span>
        )}
      </div>

      {/* Title */}
      <p className="mt-1.5 text-sm font-medium leading-snug text-text-1">
        {card.title}
      </p>

      {/* Labels */}
      {card.labels && card.labels.length > 0 && (
        <div className="mt-2 flex flex-wrap gap-1">
          {card.labels.map((label) => (
            <LabelChip key={label.id} name={label.name} color={label.color} />
          ))}
        </div>
      )}

      {/* Bottom row: assignees + due date */}
      {(dueStr || (card.assignees && card.assignees.length > 0)) && (
        <div className="mt-2 flex items-center justify-between text-xs text-text-2">
          {/* Assignee avatars */}
          <div className="flex -space-x-1">
            {card.assignees?.slice(0, 3).map((a) => (
              <div
                key={a.user_id}
                className="flex h-5 w-5 items-center justify-center rounded-full bg-surface-2 text-[10px] font-medium text-text-1 ring-1 ring-surface"
                title={a.user_name ?? a.user_id}
              >
                {(a.user_name ?? "?")[0]?.toUpperCase()}
              </div>
            ))}
          </div>
          {dueStr && <span>{dueStr}</span>}
        </div>
      )}
    </div>
  );
}
