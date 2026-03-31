import { useCallback, useEffect, useState } from "react";
import {
  getComments,
  createComment,
  updateComment,
  deleteComment,
  type Comment,
} from "@/api/comment-client";
import { useAuthStore } from "@/stores/auth-store";

interface CardDrawerCommentsProps {
  cardId: string;
}

export function CardDrawerComments({ cardId }: CardDrawerCommentsProps) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newContent, setNewContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editContent, setEditContent] = useState("");
  const currentUser = useAuthStore((s) => s.user);

  const fetchComments = useCallback(async () => {
    try {
      setIsLoading(true);
      const data = await getComments(cardId);
      setComments(data);
    } catch {
      // non-critical
    } finally {
      setIsLoading(false);
    }
  }, [cardId]);

  useEffect(() => {
    fetchComments();
  }, [fetchComments]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newContent.trim() || submitting) return;
    setSubmitting(true);
    try {
      const comment = await createComment(cardId, newContent.trim());
      setComments((prev) => [comment, ...prev]);
      setNewContent("");
    } catch {
      // silent
    } finally {
      setSubmitting(false);
    }
  };

  const handleUpdate = async (commentId: string) => {
    if (!editContent.trim()) return;
    try {
      const updated = await updateComment(commentId, editContent.trim());
      setComments((prev) =>
        prev.map((c) => (c.id === commentId ? updated : c)),
      );
      setEditingId(null);
    } catch {
      // silent
    }
  };

  const handleDelete = async (commentId: string) => {
    try {
      await deleteComment(commentId);
      setComments((prev) => prev.filter((c) => c.id !== commentId));
    } catch {
      // silent
    }
  };

  const formatTime = (dateStr: string) => {
    const d = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffMin = Math.floor(diffMs / 60000);
    if (diffMin < 1) return "just now";
    if (diffMin < 60) return `${diffMin}m ago`;
    const diffHr = Math.floor(diffMin / 60);
    if (diffHr < 24) return `${diffHr}h ago`;
    const diffDay = Math.floor(diffHr / 24);
    return `${diffDay}d ago`;
  };

  return (
    <div className="px-5 py-4">
      <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-text-2">
        Comments
      </h3>

      {/* Add comment form */}
      <form onSubmit={handleSubmit} className="mb-4">
        <textarea
          value={newContent}
          onChange={(e) => setNewContent(e.target.value)}
          placeholder="Write a comment..."
          rows={2}
          className="w-full resize-none rounded-md border border-border bg-surface-2 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:border-accent focus:outline-none"
        />
        <div className="mt-1.5 flex justify-end">
          <button
            type="submit"
            disabled={!newContent.trim() || submitting}
            className="rounded-md bg-accent px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-accent/90 disabled:opacity-50"
          >
            {submitting ? "Sending..." : "Comment"}
          </button>
        </div>
      </form>

      {/* Comment list */}
      {isLoading ? (
        <p className="text-xs text-text-2">Loading comments...</p>
      ) : comments.length === 0 ? (
        <p className="text-xs text-text-2">No comments yet</p>
      ) : (
        <div className="space-y-3">
          {comments.map((c) => (
            <div key={c.id} className="group/comment flex gap-2 text-xs">
              <div className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-2 text-[10px] font-medium text-text-1">
                {(c.author_name ?? "?")[0]?.toUpperCase()}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-text-1">
                    {c.author_name ?? "User"}
                  </span>
                  <span className="text-text-2 opacity-60">
                    {formatTime(c.created_at)}
                  </span>
                  {currentUser?.id === c.author_id && (
                    <div className="ml-auto flex gap-1 opacity-0 transition-opacity group-hover/comment:opacity-100">
                      <button
                        onClick={() => {
                          setEditingId(c.id);
                          setEditContent(c.content);
                        }}
                        className="text-text-2 hover:text-text-1"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(c.id)}
                        className="text-danger hover:text-danger/80"
                      >
                        Delete
                      </button>
                    </div>
                  )}
                </div>
                {editingId === c.id ? (
                  <div className="mt-1">
                    <textarea
                      value={editContent}
                      onChange={(e) => setEditContent(e.target.value)}
                      rows={2}
                      className="w-full resize-none rounded-md border border-border bg-surface-2 px-2 py-1 text-sm text-text-1 focus:border-accent focus:outline-none"
                    />
                    <div className="mt-1 flex gap-1.5">
                      <button
                        onClick={() => handleUpdate(c.id)}
                        className="rounded bg-accent px-2 py-0.5 text-xs text-white hover:bg-accent/90"
                      >
                        Save
                      </button>
                      <button
                        onClick={() => setEditingId(null)}
                        className="rounded bg-surface-2 px-2 py-0.5 text-xs text-text-2 hover:text-text-1"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                ) : (
                  <p className="mt-0.5 whitespace-pre-wrap text-text-1">
                    {c.content}
                  </p>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
