import type { Comment } from "@/api/comment-client";
import { useAuthStore } from "@/stores/auth-store";
import { CommentItem } from "./comment-item";

interface CommentListProps {
  comments: Comment[];
  onEdit: (commentId: string, content: string) => Promise<void>;
  onDelete: (commentId: string) => Promise<void>;
}

export function CommentList({ comments, onEdit, onDelete }: CommentListProps) {
  const currentUserId = useAuthStore((s) => s.user?.id);

  if (comments.length === 0) {
    return <p className="text-xs text-text-2">No comments yet</p>;
  }

  return (
    <div className="space-y-3">
      {comments.map((c) => (
        <CommentItem
          key={c.id}
          comment={c}
          isOwner={currentUserId === c.author_id}
          onEdit={onEdit}
          onDelete={onDelete}
        />
      ))}
    </div>
  );
}
