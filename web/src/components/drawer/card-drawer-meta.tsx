import type { CardDetail } from "@/api/card-client";

interface CardDrawerMetaProps {
  card: CardDetail;
}

export function CardDrawerMeta({ card }: CardDrawerMetaProps) {
  const formatDate = (dateStr: string | null | undefined) => {
    if (!dateStr) return null;
    return new Date(dateStr).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  return (
    <div className="space-y-3 border-b border-border px-5 py-4">
      {/* Assignees */}
      <div className="flex items-start gap-3">
        <span className="w-20 shrink-0 text-xs font-medium text-text-2">
          Assignees
        </span>
        <div className="flex flex-wrap gap-1">
          {card.assignees && card.assignees.length > 0 ? (
            card.assignees.map((a) => (
              <span
                key={a.user_id}
                className="flex h-6 items-center rounded-full bg-surface-2 px-2 text-xs font-medium text-text-1"
                title={a.user_name ?? a.user_id}
              >
                {a.user_name ?? a.user_id.slice(0, 8)}
              </span>
            ))
          ) : (
            <span className="text-xs text-text-2">No assignees</span>
          )}
        </div>
      </div>

      {/* Labels */}
      <div className="flex items-start gap-3">
        <span className="w-20 shrink-0 text-xs font-medium text-text-2">
          Labels
        </span>
        <div className="flex flex-wrap gap-1">
          {card.labels && card.labels.length > 0 ? (
            card.labels.map((label) => (
              <span
                key={label.id}
                className="rounded-full px-2 py-0.5 text-xs font-medium"
                style={{
                  backgroundColor: label.color + "20",
                  color: label.color,
                }}
              >
                {label.name}
              </span>
            ))
          ) : (
            <span className="text-xs text-text-2">No labels</span>
          )}
        </div>
      </div>

      {/* Due date */}
      <div className="flex items-center gap-3">
        <span className="w-20 shrink-0 text-xs font-medium text-text-2">
          Due date
        </span>
        <span className="text-xs text-text-1">
          {formatDate(card.due_date) ?? "Not set"}
        </span>
      </div>

      {/* Start date */}
      <div className="flex items-center gap-3">
        <span className="w-20 shrink-0 text-xs font-medium text-text-2">
          Start date
        </span>
        <span className="text-xs text-text-1">
          {formatDate(card.start_date) ?? "Not set"}
        </span>
      </div>
    </div>
  );
}
