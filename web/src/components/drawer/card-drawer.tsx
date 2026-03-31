import { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { motion, AnimatePresence } from "framer-motion";
import { useCardDetail } from "@/hooks/use-card-detail";
import { updateCard as apiUpdateCard } from "@/api/card-client";
import { useBoardStore } from "@/stores/board-store";
import { CardDrawerHeader } from "./card-drawer-header";
import { CardDrawerMeta } from "./card-drawer-meta";
import { CardDrawerActivity } from "./card-drawer-activity";
import { CardDrawerComments } from "./card-drawer-comments";
import { RichTextEditor } from "@/components/editor/rich-text-editor";
import type { Priority } from "@/types/board";

type DrawerTab = "description" | "comments" | "activity";

export function CardDrawer() {
  const { cardNumber } = useParams<{ cardNumber: string }>();
  const navigate = useNavigate();
  const storeCards = useBoardStore((s) => s.cards);
  const storeUpdateCard = useBoardStore((s) => s.updateCard);
  const [activeTab, setActiveTab] = useState<DrawerTab>("description");

  // Find card ID from store by readable_id or card_number
  const storeCard = Object.values(storeCards).find(
    (c) =>
      c.readable_id === cardNumber ||
      String(c.card_number) === cardNumber,
  );

  const cardId = storeCard?.id;
  const { card, isLoading, error, setCard } = useCardDetail(cardId);

  const handleClose = useCallback(() => {
    navigate(-1);
  }, [navigate]);

  // Esc key to close
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") handleClose();
    };
    document.addEventListener("keydown", handler);
    return () => document.removeEventListener("keydown", handler);
  }, [handleClose]);

  const handleUpdateTitle = useCallback(
    async (title: string) => {
      if (!cardId || !card) return;
      try {
        const updated = await apiUpdateCard(cardId, { title });
        setCard(updated);
        storeUpdateCard({ ...storeCard!, title });
      } catch {
        // revert handled by card state
      }
    },
    [cardId, card, setCard, storeCard, storeUpdateCard],
  );

  const handleUpdatePriority = useCallback(
    async (priority: Priority) => {
      if (!cardId || !card) return;
      try {
        const updated = await apiUpdateCard(cardId, { priority });
        setCard(updated);
        storeUpdateCard({ ...storeCard!, priority });
      } catch {
        // revert handled by card state
      }
    },
    [cardId, card, setCard, storeCard, storeUpdateCard],
  );

  const handleSaveDescription = useCallback(
    async (description: string) => {
      if (!cardId) return;
      try {
        const updated = await apiUpdateCard(cardId, { description });
        setCard(updated);
      } catch {
        // silent fail for auto-save
      }
    },
    [cardId, setCard],
  );

  return (
    <AnimatePresence>
      {/* Backdrop */}
      <motion.div
        key="backdrop"
        className="fixed inset-0 z-40 bg-black/60"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        exit={{ opacity: 0 }}
        onClick={handleClose}
      />

      {/* Drawer panel */}
      <motion.div
        key="drawer"
        className="fixed right-0 top-0 z-50 flex h-screen w-[560px] max-w-full flex-col border-l border-border bg-surface shadow-2xl"
        initial={{ x: "100%", opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        exit={{ x: "100%", opacity: 0 }}
        transition={{ type: "spring", stiffness: 300, damping: 30 }}
      >
        {isLoading && (
          <div className="flex flex-1 items-center justify-center">
            <p className="text-sm text-text-2">Loading card...</p>
          </div>
        )}

        {error && (
          <div className="flex flex-1 items-center justify-center">
            <p className="text-sm text-danger">{error}</p>
          </div>
        )}

        {card && !isLoading && (
          <>
            <CardDrawerHeader
              card={card}
              onClose={handleClose}
              onUpdateTitle={handleUpdateTitle}
              onUpdatePriority={handleUpdatePriority}
            />

            <div className="flex-1 overflow-y-auto">
              <CardDrawerMeta card={card} />

              {/* Tabs */}
              <div className="flex border-b border-border px-5">
                {(["description", "comments", "activity"] as DrawerTab[]).map(
                  (tab) => (
                    <button
                      key={tab}
                      onClick={() => setActiveTab(tab)}
                      className={`mr-4 border-b-2 pb-2 pt-3 text-xs font-medium capitalize transition-colors ${
                        activeTab === tab
                          ? "border-accent text-accent"
                          : "border-transparent text-text-2 hover:text-text-1"
                      }`}
                    >
                      {tab}
                    </button>
                  ),
                )}
              </div>

              {/* Tab content */}
              {activeTab === "description" && (
                <div className="px-5 py-4">
                  <RichTextEditor
                    initialContent={card.description ?? ""}
                    onSave={handleSaveDescription}
                  />
                </div>
              )}
              {activeTab === "comments" && (
                <CardDrawerComments cardId={card.id} />
              )}
              {activeTab === "activity" && (
                <CardDrawerActivity cardId={card.id} />
              )}
            </div>
          </>
        )}
      </motion.div>
    </AnimatePresence>
  );
}
