import type { Column } from "@/types/board";
import { WipLimitBar } from "@/components/board/wip-limit-bar";

interface ColumnHeaderProps {
  column: Column;
  cardCount: number;
}

export function ColumnHeader({ column, cardCount }: ColumnHeaderProps) {
  return (
    <div className="mb-2">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-text-1">{column.name}</h3>
        <span className="flex h-5 min-w-[20px] items-center justify-center rounded bg-surface-2 px-1.5 text-[11px] font-medium text-text-2">
          {cardCount}
        </span>
      </div>
      {column.wip_limit != null && column.wip_limit > 0 && (
        <div className="mt-1.5">
          <WipLimitBar count={cardCount} limit={column.wip_limit} />
        </div>
      )}
    </div>
  );
}
