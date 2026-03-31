import { useState } from "react";

interface CommentEditorProps {
  onSubmit: (content: string) => Promise<void>;
}

export function CommentEditor({ onSubmit }: CommentEditorProps) {
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim() || submitting) return;
    setSubmitting(true);
    try {
      await onSubmit(content.trim());
      setContent("");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-4">
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Write a comment..."
        rows={2}
        className="w-full resize-none rounded-md border border-border bg-surface-2 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:border-accent focus:outline-none"
      />
      <div className="mt-1.5 flex justify-end">
        <button
          type="submit"
          disabled={!content.trim() || submitting}
          className="rounded-md bg-accent px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-accent/90 disabled:opacity-50"
        >
          {submitting ? "Sending..." : "Comment"}
        </button>
      </div>
    </form>
  );
}
