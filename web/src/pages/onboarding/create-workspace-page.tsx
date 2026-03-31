import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createWorkspace, workspaceKeys } from "@/api/workspace-client";
import { createProject, listProjectBoards } from "@/api/project-client";
import { toast } from "@/components/ui/toast";

function generateProjectKey(name: string): string {
  return name.toUpperCase().replace(/[^A-Z]/g, "").slice(0, 4) || "PRJ";
}

const ACCENT_COLORS = [
  { hex: "#6366f1", label: "Indigo" },
  { hex: "#7c3aed", label: "Violet" },
  { hex: "#059669", label: "Emerald" },
  { hex: "#d97706", label: "Amber" },
  { hex: "#e11d48", label: "Rose" },
  { hex: "#0891b2", label: "Cyan" },
  { hex: "#475569", label: "Slate" },
  { hex: "#db2777", label: "Pink" },
] as const;

export function CreateWorkspacePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [step, setStep] = useState<1 | 2 | 3>(1);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [accentColor, setAccentColor] = useState<string>(ACCENT_COLORS[0].hex);
  const [nameError, setNameError] = useState("");
  const [workspaceSlug, setWorkspaceSlug] = useState("");
  const [projectName, setProjectName] = useState("");
  const [projectKey, setProjectKey] = useState("");
  const [keyManuallyEdited, setKeyManuallyEdited] = useState(false);

  useEffect(() => {
    if (!keyManuallyEdited) {
      setProjectKey(generateProjectKey(projectName));
    }
  }, [projectName, keyManuallyEdited]);

  const mutation = useMutation({
    mutationFn: createWorkspace,
    onSuccess: (workspace) => {
      queryClient.invalidateQueries({ queryKey: workspaceKeys.list });
      setWorkspaceSlug(workspace.slug);
      setStep(3);
    },
    onError: () => toast.error("Không thể tạo workspace. Vui lòng thử lại."),
  });

  function handleNext() {
    const trimmed = name.trim();
    if (!trimmed) { setNameError("Tên workspace không được để trống"); return; }
    if (trimmed.length > 64) { setNameError("Tên không được vượt quá 64 ký tự"); return; }
    setNameError("");
    setStep(2);
  }

  function handleSubmit() {
    mutation.mutate({ name: name.trim(), description: description.trim() || undefined, accent_color: accentColor });
  }

  const keyError =
    projectKey.length > 0 && !/^[A-Z]{2,10}$/.test(projectKey)
      ? "Key phải gồm 2–10 chữ cái in hoa (A-Z)"
      : "";

  const createProjectMutation = useMutation({
    mutationFn: () =>
      createProject(workspaceSlug, {
        name: projectName.trim(),
        key: projectKey.trim(),
      }),
    onSuccess: async (project) => {
      const boards = await listProjectBoards(workspaceSlug, project.id);
      const firstBoard = boards[0];
      if (firstBoard) {
        navigate(`/w/${workspaceSlug}/b/${firstBoard.id}`, { replace: true });
      } else {
        navigate(`/w/${workspaceSlug}`, { replace: true });
      }
    },
    onError: () => toast.error("Không thể tạo project. Vui lòng thử lại."),
  });

  const initial = name.trim()[0]?.toUpperCase() ?? "K";

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-1 p-6">
      <div className="w-full max-w-md rounded-lg border border-border bg-surface-2 p-8">
        {/* Step indicator */}
        <div className="flex items-center gap-2 mb-6">
          <div className="h-2.5 w-2.5 rounded-full bg-primary" />
          <div className={`h-0.5 w-16 ${step >= 2 ? "bg-primary" : "bg-border"}`} />
          <div className={`h-2.5 w-2.5 rounded-full border-2 ${step >= 2 ? "bg-primary border-primary" : "border-border"}`} />
          <div className={`h-0.5 w-16 ${step >= 3 ? "bg-primary" : "bg-border"}`} />
          <div className={`h-2.5 w-2.5 rounded-full border-2 ${step >= 3 ? "bg-primary border-primary" : "border-border"}`} />
          <span className="ml-2 text-xs text-text-2">Bước {step} / 3</span>
        </div>

        {step === 3 ? (
          <>
            <h1 className="text-xl font-bold text-text-1">Tạo project đầu tiên</h1>
            <p className="mt-1 mb-6 text-sm text-text-2">Tổ chức công việc theo project</p>

            <div className="space-y-4">
              <div>
                <label htmlFor="proj-name" className="block text-sm font-medium text-text-1 mb-1">
                  Tên project <span className="text-rose-500">*</span>
                </label>
                <input
                  id="proj-name"
                  type="text"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  placeholder="Ví dụ: Website Redesign"
                  className="w-full rounded-md border border-border bg-surface-1 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>

              <div>
                <label htmlFor="proj-key" className="block text-sm font-medium text-text-1 mb-1">
                  Key
                </label>
                <input
                  id="proj-key"
                  type="text"
                  value={projectKey}
                  onChange={(e) => {
                    setProjectKey(e.target.value.toUpperCase());
                    setKeyManuallyEdited(true);
                  }}
                  maxLength={10}
                  placeholder="KEY"
                  className="w-full rounded-md border border-border bg-surface-1 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
                />
                {keyError && <p className="mt-1 text-xs text-rose-500">{keyError}</p>}
              </div>
            </div>

            <div className="mt-6 flex justify-between">
              <button
                type="button"
                onClick={() => navigate(`/w/${workspaceSlug}`, { replace: true })}
                className="text-sm text-text-2 hover:text-text-1"
              >
                Bỏ qua
              </button>
              <button
                type="button"
                onClick={() => createProjectMutation.mutate()}
                disabled={!projectName.trim() || !!keyError || createProjectMutation.isPending}
                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60"
              >
                {createProjectMutation.isPending ? "Đang tạo…" : "Tạo →"}
              </button>
            </div>
          </>
        ) : step === 1 ? (
          <>
            <h1 className="text-xl font-bold text-text-1">Tạo workspace mới</h1>
            <p className="mt-1 mb-6 text-sm text-text-2">Nơi team của bạn cộng tác</p>

            <div className="space-y-4">
              <div>
                <label htmlFor="ws-name" className="block text-sm font-medium text-text-1 mb-1">
                  Tên workspace <span className="text-rose-500">*</span>
                </label>
                <input
                  id="ws-name"
                  type="text"
                  value={name}
                  onChange={(e) => { setName(e.target.value); if (nameError) setNameError(""); }}
                  maxLength={64}
                  placeholder="Ví dụ: Acme Corp"
                  className="w-full rounded-md border border-border bg-surface-1 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
                />
                {nameError && <p className="mt-1 text-xs text-rose-500">{nameError}</p>}
              </div>

              <div>
                <label htmlFor="ws-desc" className="block text-sm font-medium text-text-1 mb-1">
                  Mô tả <span className="text-text-2 font-normal">(tùy chọn)</span>
                </label>
                <textarea
                  id="ws-desc"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  maxLength={200}
                  rows={3}
                  placeholder="Mô tả ngắn về workspace của bạn"
                  className="w-full resize-none rounded-md border border-border bg-surface-1 px-3 py-2 text-sm text-text-1 placeholder:text-text-2 focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>
            </div>

            <div className="mt-6 flex justify-between">
              <button
                type="button"
                onClick={() => navigate("/")}
                className="text-sm text-text-2 hover:text-text-1"
              >
                ← Hủy
              </button>
              <button
                type="button"
                onClick={handleNext}
                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:opacity-90"
              >
                Tiếp theo →
              </button>
            </div>
          </>
        ) : (
          <>
            <h1 className="text-xl font-bold text-text-1">Chọn màu nhận diện</h1>
            <p className="mt-1 mb-6 text-sm text-text-2">Màu đại diện cho workspace của bạn</p>

            {/* Color swatches */}
            <div className="flex flex-wrap gap-3 mb-6">
              {ACCENT_COLORS.map((color) => (
                <button
                  key={color.hex}
                  type="button"
                  aria-label={color.label}
                  aria-pressed={accentColor === color.hex}
                  onClick={() => setAccentColor(color.hex)}
                  className="h-8 w-8 rounded-full transition-transform hover:scale-110 focus:outline-none"
                  style={{
                    backgroundColor: color.hex,
                    boxShadow: accentColor === color.hex ? `0 0 0 3px white, 0 0 0 5px ${color.hex}` : undefined,
                  }}
                />
              ))}
            </div>

            {/* Preview card */}
            <p className="text-xs font-medium text-text-2 uppercase tracking-wide mb-2">Xem trước</p>
            <div className="flex items-start gap-3 rounded-lg border border-border bg-surface-1 p-4">
              <div
                className="flex h-10 w-10 shrink-0 items-center justify-center rounded-md text-lg font-bold text-white"
                style={{ backgroundColor: accentColor }}
              >
                {initial}
              </div>
              <div>
                <p className="font-medium text-text-1">{name.trim() || "Tên workspace"}</p>
                {description.trim() && (
                  <p className="mt-0.5 text-xs text-text-2 line-clamp-2">{description.trim()}</p>
                )}
              </div>
            </div>

            <div className="mt-6 flex justify-between">
              <button
                type="button"
                onClick={() => setStep(1)}
                className="text-sm text-text-2 hover:text-text-1"
              >
                ← Quay lại
              </button>
              <button
                type="button"
                onClick={handleSubmit}
                disabled={mutation.isPending}
                className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60"
              >
                {mutation.isPending ? "Đang tạo…" : "Tạo workspace"}
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
