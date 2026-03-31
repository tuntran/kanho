import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { searchCards, type SearchResult } from "@/api/label-client";

interface SearchDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export function SearchDialog({ isOpen, onClose }: SearchDialogProps) {
  const navigate = useNavigate();
  const { workspaceSlug } = useParams<{ workspaceSlug: string }>();
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedIdx, setSelectedIdx] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);
  const timerRef = useRef<ReturnType<typeof setTimeout>>(undefined);

  // Focus input when opened
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => inputRef.current?.focus(), 50);
      setQuery("");
      setResults([]);
      setSelectedIdx(0);
    }
  }, [isOpen]);

  const doSearch = useCallback(
    async (q: string) => {
      if (!workspaceSlug || q.length < 2) {
        setResults([]);
        return;
      }
      setIsSearching(true);
      try {
        const data = await searchCards(workspaceSlug, q);
        setResults(data ?? []);
        setSelectedIdx(0);
      } catch {
        setResults([]);
      } finally {
        setIsSearching(false);
      }
    },
    [workspaceSlug],
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setQuery(val);
    if (timerRef.current) clearTimeout(timerRef.current);
    timerRef.current = setTimeout(() => doSearch(val), 300);
  };

  const selectResult = (result: SearchResult) => {
    onClose();
    navigate(
      `/w/${workspaceSlug}/b/${result.board_id}/cards/${result.readable_id}`,
    );
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Escape") {
      onClose();
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      setSelectedIdx((i) => Math.min(i + 1, results.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setSelectedIdx((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter" && results[selectedIdx]) {
      selectResult(results[selectedIdx]);
    }
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 z-[60] bg-black/50"
        onClick={onClose}
      />

      {/* Dialog */}
      <div className="fixed left-1/2 top-[20%] z-[61] w-full max-w-lg -translate-x-1/2 rounded-lg border border-border bg-surface shadow-2xl">
        {/* Input */}
        <div className="flex items-center border-b border-border px-4">
          <svg
            className="mr-2 h-4 w-4 shrink-0 text-text-2"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            strokeWidth={2}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            ref={inputRef}
            value={query}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            placeholder="Search cards..."
            className="flex-1 bg-transparent py-3 text-sm text-text-1 placeholder:text-text-2 focus:outline-none"
          />
          <kbd className="ml-2 rounded border border-border px-1.5 py-0.5 text-[10px] text-text-2">
            ESC
          </kbd>
        </div>

        {/* Results */}
        <div className="max-h-64 overflow-y-auto p-2">
          {isSearching && (
            <p className="px-2 py-3 text-center text-xs text-text-2">
              Searching...
            </p>
          )}
          {!isSearching && query.length >= 2 && results.length === 0 && (
            <p className="px-2 py-3 text-center text-xs text-text-2">
              No results found
            </p>
          )}
          {!isSearching && query.length > 0 && query.length < 2 && (
            <p className="px-2 py-3 text-center text-xs text-text-2">
              Type at least 2 characters
            </p>
          )}
          {results.map((r, idx) => (
            <button
              key={r.card_id}
              onClick={() => selectResult(r)}
              className={`flex w-full items-center gap-3 rounded-md px-3 py-2 text-left text-sm transition-colors ${
                idx === selectedIdx
                  ? "bg-accent/10 text-accent"
                  : "text-text-1 hover:bg-surface-2"
              }`}
            >
              <span className="shrink-0 font-mono text-xs font-medium text-text-2">
                {r.readable_id}
              </span>
              <span className="flex-1 truncate">{r.title}</span>
              <span className="shrink-0 text-[10px] text-text-2">
                {r.board_name} &middot; {r.column_name}
              </span>
            </button>
          ))}
        </div>
      </div>
    </>
  );
}
