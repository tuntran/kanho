import { useCallback, useEffect, useState } from "react";
import {
  getComments,
  createComment,
  updateComment,
  deleteComment,
  type Comment,
} from "@/api/comment-client";
import { CommentEditor } from "./comment-editor";
import { CommentList } from "./comment-list";

interface CardDrawerCommentsProps {
  cardId: string;
}

export function CardDrawerComments({ cardId }: CardDrawerCommentsProps) {
  const [comments, setComments] = useState<Comment[]>([]);
  const [isLoading, setIsLoading] = useState(true);

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

  const handleSubmit = async (content: string) => {
    const comment = await createComment(cardId, content);
    setComments((prev) => [comment, ...prev]);
  };

  const handleEdit = async (commentId: string, content: string) => {
    const updated = await updateComment(commentId, content);
    setComments((prev) =>
      prev.map((c) => (c.id === commentId ? updated : c)),
    );
  };

  const handleDelete = async (commentId: string) => {
    await deleteComment(commentId);
    setComments((prev) => prev.filter((c) => c.id !== commentId));
  };

  return (
    <div className="px-5 py-4">
      <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-text-2">
        Comments
      </h3>
      <CommentEditor onSubmit={handleSubmit} />
      {isLoading ? (
        <p className="text-xs text-text-2">Loading comments...</p>
      ) : (
        <CommentList
          comments={comments}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      )}
    </div>
  );
}
