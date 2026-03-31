import { useState } from "react";
import { createCard } from "@/api/board-client";
import { useBoardStore } from "@/stores/board-store";

interface AddCardFormProps {
  boardId: string;
  columnId: string;
}

export function AddCardForm({ boardId, columnId }: AddCardFormProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const addCard = useBoardStore((s) => s.addCard);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = title.trim();
    if (!trimmed || submitting) return;

    setSubmitting(true);
    try {
      const card = await createCard(boardId, {
        column_id: columnId,
        title: trimmed,
      });
      addCard(card);
      setTitle("");
      setIsOpen(false);
    } catch {
      // TODO: toast error
    } finally {
      setSubmitting(false);
    }
  };

  if (!isOpen) {
    return (
      <button
        onClick={() => setIsOpen(true)}
        className="mt-1 w-full rounded-md px-2 py-1.5 text-left text-sm text-text-2 transition-colors hover:bg-surface-2 hover:text-text-1"
      >
        + Add card
      </button>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="mt-1">
      <textarea
        autoFocus
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            handleSubmit(e);
          }
          if (e.key === "Escape") {
            setIsOpen(false);
            setTitle("");
          }
        }}
        placeholder="Card title..."
        className="w-full resize-none rounded-card border border-border bg-surface p-2 text-sm text-text-1 placeholder:text-text-3 focus:border-border-focus focus:outline-none"
        rows={2}
      />
      <div className="mt-1.5 flex gap-2">
        <button
          type="submit"
          disabled={submitting || !title.trim()}
          className="rounded-md bg-primary px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-primary-hover disabled:opacity-50"
        >
          Add
        </button>
        <button
          type="button"
          onClick={() => {
            setIsOpen(false);
            setTitle("");
          }}
          className="rounded-md px-3 py-1 text-xs text-text-2 transition-colors hover:bg-surface-2"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
