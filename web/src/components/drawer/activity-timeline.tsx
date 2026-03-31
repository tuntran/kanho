import { useMemo } from "react";
import { formatDistanceToNow } from "date-fns";
import type { Activity } from "@/api/card-client";
import type { Comment } from "@/api/comment-client";

type TimelineEntry =
  | { kind: "activity"; data: Activity }
  | { kind: "comment"; data: Comment };

interface ActivityTimelineProps {
  activities: Activity[];
  comments: Comment[];
}

function describeAction(a: Activity): string {
  if (a.field === "column_id" && a.old_value && a.new_value) {
    return `moved card from "${a.old_value}" to "${a.new_value}"`;
  }
  if (a.field) {
    return `changed ${a.field}`;
  }
  return a.action;
}

export function ActivityTimeline({ activities, comments }: ActivityTimelineProps) {
  const entries = useMemo(() => {
    const merged: TimelineEntry[] = [
      ...activities.map((a) => ({ kind: "activity" as const, data: a })),
      ...comments.map((c) => ({ kind: "comment" as const, data: c })),
    ];
    merged.sort(
      (a, b) =>
        new Date(b.data.created_at).getTime() -
        new Date(a.data.created_at).getTime(),
    );
    return merged;
  }, [activities, comments]);

  if (entries.length === 0) {
    return <p className="text-xs text-text-2">No activity yet</p>;
  }

  return (
    <div className="space-y-2.5">
      {entries.map((entry) => {
        if (entry.kind === "activity") {
          const a = entry.data;
          const initial = (a.user_name ?? "?")[0]?.toUpperCase() ?? "?";
          const timeAgo = formatDistanceToNow(new Date(a.created_at), {
            addSuffix: true,
          });
          return (
            <div key={`a-${a.id}`} className="flex items-start gap-2 text-xs">
              <div className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-surface-2 text-[10px] font-medium text-text-1">
                {initial}
              </div>
              <div className="min-w-0">
                <span className="font-medium text-text-1">
                  {a.user_name ?? "User"}
                </span>{" "}
                <span className="text-text-2">{describeAction(a)}</span>
                <span className="ml-1.5 text-text-2 opacity-60">{timeAgo}</span>
              </div>
            </div>
          );
        }

        const c = entry.data;
        const initial = (c.author_name ?? "?")[0]?.toUpperCase() ?? "?";
        const timeAgo = formatDistanceToNow(new Date(c.created_at), {
          addSuffix: true,
        });
        return (
          <div key={`c-${c.id}`} className="flex items-start gap-2 text-xs">
            <div className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-accent/20 text-[10px] font-medium text-accent">
              {initial}
            </div>
            <div className="min-w-0 flex-1">
              <span className="font-medium text-text-1">
                {c.author_name ?? "User"}
              </span>{" "}
              <span className="text-text-2">commented</span>
              <span className="ml-1.5 text-text-2 opacity-60">{timeAgo}</span>
              <p className="mt-0.5 whitespace-pre-wrap rounded border border-border bg-surface-2 px-2 py-1 text-text-1">
                {c.content}
              </p>
            </div>
          </div>
        );
      })}
    </div>
  );
}
