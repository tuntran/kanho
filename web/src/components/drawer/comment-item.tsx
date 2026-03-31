import { useState } from "react";
import type { Comment } from "@/api/comment-client";
import { formatDistanceToNow } from "date-fns";

interface CommentItemProps {
  comment: Comment;
  isOwner: boolean;
  onEdit: (commentId: string, content: string) => Promise<void>;
  onDelete: (commentId: string) => Promise<void>;
}

export function CommentItem({ comment, isOwner, onEdit, onDelete }: CommentItemProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editContent, setEditContent] = useState(comment.content);
  const [saving, setSaving] = useState(false);

  const handleSave = async () => {
    if (!editContent.trim() || saving) return;
    setSaving(true);
    try {
      await onEdit(comment.id, editContent.trim());
      setIsEditing(false);
    } finally {
      setSaving(false);
    }
  };

  const initial = (comment.author_name ?? "?")[0]?.toUpperCase() ?? "?";
  const timeAgo = formatDistanceToNow(new Date(comment.created_at), { addSuffix: true });

  return (
    <div className="group/comment flex gap-2 text-xs">
      <div className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-surface-2 text-[10px] font-medium text-text-1">
        {initial}
      </div>
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className="font-medium text-text-1">
            {comment.author_name ?? "User"}
          </span>
          <span className="text-text-2 opacity-60">{timeAgo}</span>
          {isOwner && !isEditing && (
            <div className="ml-auto flex gap-1 opacity-0 transition-opacity group-hover/comment:opacity-100">
              <button
                onClick={() => {
                  setEditContent(comment.content);
                  setIsEditing(true);
                }}
                className="text-text-2 hover:text-text-1"
              >
                Edit
              </button>
              <button
                onClick={() => onDelete(comment.id)}
                className="text-danger hover:text-danger/80"
              >
                Delete
              </button>
            </div>
          )}
        </div>
        {isEditing ? (
          <div className="mt-1">
            <textarea
              value={editContent}
              onChange={(e) => setEditContent(e.target.value)}
              rows={2}
              className="w-full resize-none rounded-md border border-border bg-surface-2 px-2 py-1 text-sm text-text-1 focus:border-accent focus:outline-none"
            />
            <div className="mt-1 flex gap-1.5">
              <button
                onClick={handleSave}
                disabled={saving}
                className="rounded bg-accent px-2 py-0.5 text-xs text-white hover:bg-accent/90 disabled:opacity-50"
              >
                {saving ? "Saving..." : "Save"}
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="rounded bg-surface-2 px-2 py-0.5 text-xs text-text-2 hover:text-text-1"
              >
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <p className="mt-0.5 whitespace-pre-wrap text-text-1">
            {comment.content}
          </p>
        )}
      </div>
    </div>
  );
}
