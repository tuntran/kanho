import { useCallback, useEffect, useState } from "react";
import { getActivities, type Activity } from "@/api/card-client";

interface CardDrawerActivityProps {
  cardId: string;
}

export function CardDrawerActivity({ cardId }: CardDrawerActivityProps) {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchActivities = useCallback(async () => {
    try {
      setIsLoading(true);
      const data = await getActivities(cardId);
      setActivities(data);
    } catch {
      // Silently fail — activity is non-critical
    } finally {
      setIsLoading(false);
    }
  }, [cardId]);

  useEffect(() => {
    fetchActivities();
  }, [fetchActivities]);

  const formatTime = (dateStr: string) => {
    const d = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffMin = Math.floor(diffMs / 60000);
    if (diffMin < 1) return "just now";
    if (diffMin < 60) return `${diffMin}m ago`;
    const diffHr = Math.floor(diffMin / 60);
    if (diffHr < 24) return `${diffHr}h ago`;
    const diffDay = Math.floor(diffHr / 24);
    return `${diffDay}d ago`;
  };

  const describeAction = (a: Activity) => {
    if (a.field) {
      return `changed ${a.field}`;
    }
    return a.action;
  };

  return (
    <div className="px-5 py-4">
      <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-text-2">
        Activity
      </h3>
      {isLoading ? (
        <p className="text-xs text-text-2">Loading activity...</p>
      ) : activities.length === 0 ? (
        <p className="text-xs text-text-2">No activity yet</p>
      ) : (
        <div className="space-y-2.5">
          {activities.map((a) => (
            <div key={a.id} className="flex items-start gap-2 text-xs">
              <div className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-surface-2 text-[10px] font-medium text-text-1">
                {(a.user_name ?? "?")[0]?.toUpperCase()}
              </div>
              <div className="min-w-0">
                <span className="font-medium text-text-1">
                  {a.user_name ?? "User"}
                </span>{" "}
                <span className="text-text-2">{describeAction(a)}</span>
                <span className="ml-1.5 text-text-2 opacity-60">
                  {formatTime(a.created_at)}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
