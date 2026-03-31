import { useCallback, useEffect, useState } from "react";
import { getActivities, type Activity } from "@/api/card-client";
import { getComments, type Comment } from "@/api/comment-client";
import { ActivityTimeline } from "./activity-timeline";

interface CardDrawerActivityProps {
  cardId: string;
}

export function CardDrawerActivity({ cardId }: CardDrawerActivityProps) {
  const [activities, setActivities] = useState<Activity[]>([]);
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      setIsLoading(true);
      const [acts, cmts] = await Promise.all([
        getActivities(cardId),
        getComments(cardId),
      ]);
      setActivities(acts);
      setComments(cmts);
    } catch {
      // non-critical
    } finally {
      setIsLoading(false);
    }
  }, [cardId]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div className="px-5 py-4">
      <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-text-2">
        Activity
      </h3>
      {isLoading ? (
        <p className="text-xs text-text-2">Loading activity...</p>
      ) : (
        <ActivityTimeline activities={activities} comments={comments} />
      )}
    </div>
  );
}
