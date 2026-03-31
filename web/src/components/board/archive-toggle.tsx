interface ArchiveToggleProps {
  showArchived: boolean;
  onToggle: () => void;
}

export function ArchiveToggle({ showArchived, onToggle }: ArchiveToggleProps) {
  return (
    <button
      onClick={onToggle}
      className="flex items-center gap-1.5 rounded-md border border-border bg-surface px-2.5 py-1 text-xs font-medium text-text-2 transition-colors hover:border-border-focus hover:text-text-1"
    >
      <svg
        className="h-3.5 w-3.5"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        strokeWidth={2}
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
        />
      </svg>
      {showArchived ? "Hide archived" : "Show archived"}
    </button>
  );
}
