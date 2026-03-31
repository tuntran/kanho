import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createProject,
  listProjectBoards,
  projectKeys,
} from "@/api/project-client";
import { toast } from "@/components/ui/toast";

function generateProjectKey(name: string): string {
  return name.toUpperCase().replace(/[^A-Z]/g, "").slice(0, 4) || "PRJ";
}

const KEY_PATTERN = /^[A-Z]{2,10}$/;

interface CreateProjectModalProps {
  workspaceSlug: string;
  isOpen: boolean;
  onClose: () => void;
}

export function CreateProjectModal({
  workspaceSlug,
  isOpen,
  onClose,
}: CreateProjectModalProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [name, setName] = useState("");
  const [key, setKey] = useState("");
  const [keyTouched, setKeyTouched] = useState(false);
  const [keyError, setKeyError] = useState("");

  useEffect(() => {
    if (!keyTouched) {
      setKey(generateProjectKey(name));
    }
  }, [name, keyTouched]);

  const mutation = useMutation({
    mutationFn: (input: { name: string; key: string }) =>
      createProject(workspaceSlug, input),
    onSuccess: async (project) => {
      queryClient.invalidateQueries({
        queryKey: projectKeys.list(workspaceSlug),
      });
      try {
        const boards = await listProjectBoards(workspaceSlug, project.id);
        if (boards.length > 0) {
          navigate(`/w/${workspaceSlug}/b/${boards[0]!.id}`);
        }
      } catch {
        // Board auto-created by backend — stay on dashboard
      }
      onClose();
      resetForm();
    },
    onError: () => toast.error("Khong the tao project. Vui long thu lai."),
  });

  function resetForm() {
    setName("");
    setKey("");
    setKeyTouched(false);
    setKeyError("");
  }

  function handleKeyChange(value: string) {
    const upper = value.toUpperCase();
    setKey(upper);
    setKeyTouched(true);
    setKeyError(
      upper && !KEY_PATTERN.test(upper)
        ? "Key phai tu 2-10 ky tu in hoa (A-Z)"
        : "",
    );
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName || !KEY_PATTERN.test(key)) return;
    mutation.mutate({ name: trimmedName, key });
  }

  function handleClose() {
    onClose();
    resetForm();
  }

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={handleClose}
    >
      <div
        className="w-full max-w-md rounded-lg border border-border bg-surface-2 p-6"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-lg font-bold text-text-1">Tao project moi</h2>
        <p className="mt-1 text-sm text-text-2">
          Project giup nhom cac board lien quan
        </p>

        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label
              htmlFor="project-name"
              className="block text-sm font-medium text-text-1 mb-1"
            >
              Ten project <span className="text-rose-500">*</span>
            </label>
            <input
              id="project-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              maxLength={64}
              placeholder="Vi du: Marketing Campaign"
              className="w-full rounded-md border border-border bg-surface-1 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>

          <div>
            <label
              htmlFor="project-key"
              className="block text-sm font-medium text-text-1 mb-1"
            >
              Key
            </label>
            <input
              id="project-key"
              type="text"
              value={key}
              onChange={(e) => handleKeyChange(e.target.value)}
              maxLength={10}
              placeholder="VD: MKTG"
              className="w-full rounded-md border border-border bg-surface-1 px-3 py-2 text-sm font-mono text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
            />
            {keyError && (
              <p className="mt-1 text-xs text-rose-500">{keyError}</p>
            )}
          </div>

          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={handleClose}
              className="rounded-md px-4 py-2 text-sm text-text-2 hover:text-text-1"
            >
              Huy
            </button>
            <button
              type="submit"
              disabled={
                !name.trim() || !KEY_PATTERN.test(key) || mutation.isPending
              }
              className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60"
            >
              {mutation.isPending ? "Dang tao..." : "Tao project"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
