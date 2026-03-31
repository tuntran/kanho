import type { Column } from "@/types/board";

export function ColumnDragOverlay({ column }: { column: Column }) {
  return (
    <div className="w-[280px] rotate-1 scale-[1.02] rounded-panel border-2 border-dashed border-primary bg-surface/80 p-3 shadow-card-drag">
      <h3 className="text-sm font-semibold text-text-1">{column.name}</h3>
      <p className="mt-1 text-xs text-text-2">Dragging column...</p>
    </div>
  );
}
