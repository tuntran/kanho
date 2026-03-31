import { useState, useRef, useEffect } from "react";
import type { CardDetail } from "@/api/card-client";
import type { Priority } from "@/types/board";

const PRIORITIES: { value: Priority; label: string; color: string }[] = [
  { value: "urgent", label: "Urgent", color: "bg-danger" },
  { value: "high", label: "High", color: "bg-orange-500" },
  { value: "medium", label: "Medium", color: "bg-yellow-500" },
  { value: "low", label: "Low", color: "bg-blue-500" },
  { value: "none", label: "None", color: "bg-text-2" },
];

interface CardDrawerHeaderProps {
  card: CardDetail;
  onClose: () => void;
  onUpdateTitle: (title: string) => void;
  onUpdatePriority: (priority: Priority) => void;
}

export function CardDrawerHeader({
  card,
  onClose,
  onUpdateTitle,
  onUpdatePriority,
}: CardDrawerHeaderProps) {
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [titleDraft, setTitleDraft] = useState(card.title);
  const [showPriorityMenu, setShowPriorityMenu] = useState(false);
  const titleInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    setTitleDraft(card.title);
  }, [card.title]);

  useEffect(() => {
    if (isEditingTitle) {
      titleInputRef.current?.focus();
      titleInputRef.current?.select();
    }
  }, [isEditingTitle]);

  const handleTitleBlur = () => {
    setIsEditingTitle(false);
    const trimmed = titleDraft.trim();
    if (trimmed && trimmed !== card.title) {
      onUpdateTitle(trimmed);
    } else {
      setTitleDraft(card.title);
    }
  };

  const handleTitleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      (e.target as HTMLInputElement).blur();
    } else if (e.key === "Escape") {
      setTitleDraft(card.title);
      setIsEditingTitle(false);
    }
  };

  const currentPriority =
    PRIORITIES.find((p) => p.value === card.priority) ?? PRIORITIES[4]!;

  return (
    <div className="border-b border-border px-5 py-4">
      {/* Top row: readable ID + close button */}
      <div className="flex items-center justify-between">
        {card.readable_id && (
          <span className="font-mono text-xs font-semibold text-primary">
            {card.readable_id}
          </span>
        )}
        <button
          onClick={onClose}
          className="rounded p-1 text-text-2 transition-colors hover:bg-surface-2 hover:text-text-1"
          aria-label="Close drawer"
        >
          <svg
            width="20"
            height="20"
            viewBox="0 0 20 20"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
          >
            <path d="M5 5l10 10M15 5L5 15" />
          </svg>
        </button>
      </div>

      {/* Title (inline editable) */}
      <div className="mt-2">
        {isEditingTitle ? (
          <input
            ref={titleInputRef}
            value={titleDraft}
            onChange={(e) => setTitleDraft(e.target.value)}
            onBlur={handleTitleBlur}
            onKeyDown={handleTitleKeyDown}
            className="w-full rounded bg-transparent text-lg font-semibold text-text-1 outline-none ring-1 ring-border-focus px-1 -ml-1"
          />
        ) : (
          <h2
            onClick={() => setIsEditingTitle(true)}
            className="cursor-text text-lg font-semibold text-text-1 hover:text-primary transition-colors"
          >
            {card.title}
          </h2>
        )}
      </div>

      {/* Priority selector */}
      <div className="relative mt-3">
        <button
          onClick={() => setShowPriorityMenu(!showPriorityMenu)}
          className="flex items-center gap-1.5 rounded-md border border-border px-2.5 py-1 text-xs font-medium text-text-2 transition-colors hover:border-border-focus hover:text-text-1"
        >
          <span
            className={`inline-block h-2 w-2 rounded-full ${currentPriority.color}`}
          />
          {currentPriority.label}
          <svg
            width="12"
            height="12"
            viewBox="0 0 12 12"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
          >
            <path d="M3 5l3 3 3-3" />
          </svg>
        </button>
        {showPriorityMenu && (
          <div className="absolute left-0 top-full z-10 mt-1 w-36 rounded-lg border border-border bg-surface shadow-lg">
            {PRIORITIES.map((p) => (
              <button
                key={p.value}
                onClick={() => {
                  onUpdatePriority(p.value);
                  setShowPriorityMenu(false);
                }}
                className="flex w-full items-center gap-2 px-3 py-1.5 text-xs text-text-2 transition-colors hover:bg-surface-2 hover:text-text-1"
              >
                <span
                  className={`inline-block h-2 w-2 rounded-full ${p.color}`}
                />
                {p.label}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
