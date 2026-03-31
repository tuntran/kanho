import { useCallback, useEffect, useState } from "react";
import { getCard, type CardDetail } from "@/api/card-client";

export function useCardDetail(cardId: string | undefined) {
  const [card, setCard] = useState<CardDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCard = useCallback(async () => {
    if (!cardId) return;
    try {
      setIsLoading(true);
      setError(null);
      const data = await getCard(cardId);
      setCard(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load card");
    } finally {
      setIsLoading(false);
    }
  }, [cardId]);

  useEffect(() => {
    fetchCard();
  }, [fetchCard]);

  return { card, isLoading, error, refetch: fetchCard, setCard };
}
