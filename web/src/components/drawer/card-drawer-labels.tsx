import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  getWorkspaceLabels,
  addCardLabel,
  removeCardLabel,
} from "@/api/label-client";
import type { Label } from "@/types/board";

interface CardDrawerLabelsProps {
  cardId: string;
  currentLabels: Label[];
  onLabelsChange: (labels: Label[]) => void;
}

export function CardDrawerLabels({
  cardId,
  currentLabels,
  onLabelsChange,
}: CardDrawerLabelsProps) {
  const { workspaceSlug } = useParams<{ workspaceSlug: string }>();
  const [allLabels, setAllLabels] = useState<Label[]>([]);
  const [isOpen, setIsOpen] = useState(false);

  const fetchLabels = useCallback(async () => {
    if (!workspaceSlug) return;
    try {
      const labels = await getWorkspaceLabels(workspaceSlug);
      setAllLabels(labels);
    } catch {
      // silent
    }
  }, [workspaceSlug]);

  useEffect(() => {
    if (isOpen) fetchLabels();
  }, [isOpen, fetchLabels]);

  const currentLabelIds = new Set(currentLabels.map((l) => l.id));

  const toggleLabel = async (label: Label) => {
    const isAssigned = currentLabelIds.has(label.id);
    try {
      if (isAssigned) {
        await removeCardLabel(cardId, label.id);
        onLabelsChange(currentLabels.filter((l) => l.id !== label.id));
      } else {
        await addCardLabel(cardId, label.id);
        onLabelsChange([...currentLabels, label]);
      }
    } catch {
      // silent
    }
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="rounded bg-surface-2 px-2 py-0.5 text-xs text-text-2 hover:text-text-1"
      >
        + Label
      </button>

      {isOpen && (
        <>
          <div
            className="fixed inset-0 z-50"
            onClick={() => setIsOpen(false)}
          />
          <div className="absolute left-0 top-full z-50 mt-1 w-48 rounded-md border border-border bg-surface p-2 shadow-lg">
            <p className="mb-1.5 text-[10px] font-semibold uppercase tracking-wider text-text-2">
              Labels
            </p>
            {allLabels.length === 0 ? (
              <p className="py-1 text-xs text-text-2">No labels yet</p>
            ) : (
              <div className="max-h-40 space-y-0.5 overflow-y-auto">
                {allLabels.map((label) => (
                  <button
                    key={label.id}
                    onClick={() => toggleLabel(label)}
                    className="flex w-full items-center gap-2 rounded px-1.5 py-1 text-xs hover:bg-surface-2"
                  >
                    <span
                      className="h-3 w-3 shrink-0 rounded-sm"
                      style={{ backgroundColor: label.color }}
                    />
                    <span className="flex-1 text-left text-text-1">
                      {label.name}
                    </span>
                    {currentLabelIds.has(label.id) && (
                      <span className="text-accent">&#10003;</span>
                    )}
                  </button>
                ))}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
