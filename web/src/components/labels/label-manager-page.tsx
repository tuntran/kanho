import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getWorkspaceLabels } from "@/api/label-client";
import type { Label } from "@/types/board";
import { LabelColorPicker } from "./label-color-picker";
import { LabelChip } from "@/components/ui/label-chip";
import { EmptyState } from "@/components/ui/empty-state";
import apiClient from "@/lib/api-interceptor";

export function LabelManagerPage() {
  const { workspaceSlug } = useParams<{ workspaceSlug: string }>();
  const [labels, setLabels] = useState<Label[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newName, setNewName] = useState("");
  const [newColor, setNewColor] = useState("#3b82f6");
  const [creating, setCreating] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [editColor, setEditColor] = useState("");

  const fetchLabels = useCallback(async () => {
    if (!workspaceSlug) return;
    try {
      setIsLoading(true);
      const data = await getWorkspaceLabels(workspaceSlug);
      setLabels(data);
    } catch {
      // non-critical
    } finally {
      setIsLoading(false);
    }
  }, [workspaceSlug]);

  useEffect(() => {
    fetchLabels();
  }, [fetchLabels]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newName.trim() || !workspaceSlug || creating) return;
    setCreating(true);
    try {
      const { data } = await apiClient.post<Label>(
        `/workspaces/${workspaceSlug}/labels`,
        { name: newName.trim(), color: newColor },
      );
      setLabels((prev) => [...prev, data]);
      setNewName("");
    } catch {
      // silent
    } finally {
      setCreating(false);
    }
  };

  const handleUpdate = async (labelId: string) => {
    if (!editName.trim() || !workspaceSlug) return;
    try {
      const { data } = await apiClient.patch<Label>(
        `/workspaces/${workspaceSlug}/labels/${labelId}`,
        { name: editName.trim(), color: editColor },
      );
      setLabels((prev) => prev.map((l) => (l.id === labelId ? data : l)));
      setEditingId(null);
    } catch {
      // silent
    }
  };

  const handleDelete = async (labelId: string) => {
    if (!workspaceSlug) return;
    if (!window.confirm("Delete this label?")) return;
    try {
      await apiClient.delete(`/workspaces/${workspaceSlug}/labels/${labelId}`);
      setLabels((prev) => prev.filter((l) => l.id !== labelId));
    } catch {
      // silent
    }
  };

  const startEditing = (label: Label) => {
    setEditingId(label.id);
    setEditName(label.name);
    setEditColor(label.color);
  };

  return (
    <div className="mx-auto max-w-xl px-5 py-8">
      <h2 className="text-lg font-bold text-text-1">Workspace Labels</h2>
      <p className="mt-1 text-xs text-text-2">
        Manage labels for this workspace. Labels can be added to cards.
      </p>

      {/* Create form */}
      <form onSubmit={handleCreate} className="mt-6 space-y-3 rounded-lg border border-border bg-surface p-4">
        <h3 className="text-xs font-semibold uppercase tracking-wider text-text-2">
          New Label
        </h3>
        <input
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder="Label name"
          className="w-full rounded-md border border-border bg-surface-2 px-3 py-1.5 text-sm text-text-1 placeholder:text-text-2 focus:border-accent focus:outline-none"
        />
        <LabelColorPicker value={newColor} onChange={setNewColor} />
        {newName.trim() && (
          <div>
            <span className="text-xs text-text-2">Preview: </span>
            <LabelChip name={newName} color={newColor} />
          </div>
        )}
        <button
          type="submit"
          disabled={!newName.trim() || creating}
          className="rounded-md bg-accent px-3 py-1.5 text-xs font-medium text-white hover:bg-accent/90 disabled:opacity-50"
        >
          {creating ? "Creating..." : "Create Label"}
        </button>
      </form>

      {/* Label list */}
      <div className="mt-6">
        {isLoading ? (
          <p className="text-xs text-text-2">Loading labels...</p>
        ) : labels.length === 0 ? (
          <EmptyState
            title="No labels yet"
            description="Create your first label above"
          />
        ) : (
          <div className="space-y-2">
            {labels.map((label) => (
              <LabelRow
                key={label.id}
                label={label}
                isEditing={editingId === label.id}
                editName={editName}
                editColor={editColor}
                onEditName={setEditName}
                onEditColor={setEditColor}
                onStartEdit={() => startEditing(label)}
                onSave={() => handleUpdate(label.id)}
                onCancel={() => setEditingId(null)}
                onDelete={() => handleDelete(label.id)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

interface LabelRowProps {
  label: Label;
  isEditing: boolean;
  editName: string;
  editColor: string;
  onEditName: (name: string) => void;
  onEditColor: (color: string) => void;
  onStartEdit: () => void;
  onSave: () => void;
  onCancel: () => void;
  onDelete: () => void;
}

function LabelRow({
  label,
  isEditing,
  editName,
  editColor,
  onEditName,
  onEditColor,
  onStartEdit,
  onSave,
  onCancel,
  onDelete,
}: LabelRowProps) {
  if (isEditing) {
    return (
      <div className="space-y-2 rounded-lg border border-accent/40 bg-surface p-3">
        <input
          value={editName}
          onChange={(e) => onEditName(e.target.value)}
          className="w-full rounded-md border border-border bg-surface-2 px-2 py-1 text-sm text-text-1 focus:border-accent focus:outline-none"
        />
        <LabelColorPicker value={editColor} onChange={onEditColor} />
        <div className="flex gap-1.5">
          <button
            onClick={onSave}
            className="rounded bg-accent px-2 py-0.5 text-xs text-white hover:bg-accent/90"
          >
            Save
          </button>
          <button
            onClick={onCancel}
            className="rounded bg-surface-2 px-2 py-0.5 text-xs text-text-2 hover:text-text-1"
          >
            Cancel
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="group flex items-center justify-between rounded-lg border border-border bg-surface px-3 py-2">
      <button onClick={onStartEdit} className="flex items-center gap-2">
        <LabelChip name={label.name} color={label.color} />
      </button>
      <div className="flex gap-1 opacity-0 transition-opacity group-hover:opacity-100">
        <button
          onClick={onStartEdit}
          className="text-xs text-text-2 hover:text-text-1"
        >
          Edit
        </button>
        <button
          onClick={onDelete}
          className="text-xs text-danger hover:text-danger/80"
        >
          Delete
        </button>
      </div>
    </div>
  );
}
