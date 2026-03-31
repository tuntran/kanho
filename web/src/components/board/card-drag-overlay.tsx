import type { Card } from "@/types/board";
import { PriorityBadge } from "@/components/ui/priority-badge";

export function CardDragOverlay({ card }: { card: Card }) {
  return (
    <div className="w-[256px] rotate-1 scale-[1.02] rounded-card border border-primary bg-surface p-3 shadow-card-drag">
      <div className="flex items-center justify-between gap-2">
        <PriorityBadge priority={card.priority} />
        {card.readable_id && (
          <span className="font-mono text-xs font-medium text-text-2">
            {card.readable_id}
          </span>
        )}
      </div>
      <p className="mt-1.5 text-sm font-medium leading-snug text-text-1">
        {card.title}
      </p>
    </div>
  );
}
