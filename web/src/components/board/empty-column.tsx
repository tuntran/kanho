interface EmptyColumnProps {
  onAddCard: () => void;
}

export function EmptyColumn({ onAddCard }: EmptyColumnProps) {
  return (
    <div className="flex flex-col items-center justify-center py-8 text-center">
      <div className="mb-2 text-2xl opacity-30">&#128196;</div>
      <p className="text-xs text-text-2">No cards yet</p>
      <button
        onClick={onAddCard}
        className="mt-2 rounded-md bg-surface-2 px-3 py-1 text-xs font-medium text-text-2 transition-colors hover:bg-surface-2/80 hover:text-text-1"
      >
        + Add a card
      </button>
    </div>
  );
}
