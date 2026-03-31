// TODO: view tạm — cần design chính thức
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "@/hooks/use-auth";
import { listWorkspaces, workspaceKeys, type Workspace } from "@/api/workspace-client";

export function WorkspaceSelectorPage() {
  const navigate = useNavigate();
  const { user } = useAuth();

  const { data: workspaces, isLoading } = useQuery({
    queryKey: workspaceKeys.list,
    queryFn: listWorkspaces,
  });

  useEffect(() => {
    if (!workspaces) return;
    if (workspaces.length === 0) navigate("/onboarding/create-workspace", { replace: true });
    else if (workspaces.length === 1) navigate(`/w/${workspaces[0]!.slug}`, { replace: true });
  }, [workspaces, navigate]);

  if (isLoading) return <LoadingSkeleton />;
  // Single workspace case: useEffect handles redirect — render nothing while redirecting
  if (!workspaces || workspaces.length <= 1) return null;

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-1 p-6">
      <div className="w-full max-w-2xl">
        <h1 className="text-2xl font-bold text-text-1">
          Chào mừng, {user?.name}!
        </h1>
        <p className="mt-1 text-text-2">Chọn workspace để bắt đầu</p>

        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3">
          {workspaces.map((ws) => (
            <WorkspaceCard key={ws.id} workspace={ws} />
          ))}
        </div>

        <button
          onClick={() => navigate("/onboarding/create-workspace")}
          className="mt-6 text-sm text-text-2 hover:text-text-1"
        >
          + Tạo workspace mới
        </button>
      </div>
    </div>
  );
}

function WorkspaceCard({ workspace }: { workspace: Workspace }) {
  const navigate = useNavigate();
  const initial = workspace.name[0]?.toUpperCase() ?? "?";

  return (
    <button
      onClick={() => navigate(`/w/${workspace.slug}`)}
      className="flex flex-col items-start gap-3 rounded-lg border border-border bg-surface-2 p-4 text-left hover:bg-surface-3 transition-colors"
    >
      <div
        className="flex h-10 w-10 items-center justify-center rounded-md text-lg font-bold text-white"
        style={{ backgroundColor: workspace.accent_color ?? "#6366f1" }}
      >
        {initial}
      </div>
      <div>
        <p className="font-medium text-text-1">{workspace.name}</p>
        {workspace.description && (
          <p className="mt-0.5 text-xs text-text-2 line-clamp-2">{workspace.description}</p>
        )}
      </div>
    </button>
  );
}

function LoadingSkeleton() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-1 p-6">
      <div className="w-full max-w-2xl">
        <div className="h-8 w-48 animate-pulse rounded bg-surface-3" />
        <div className="mt-2 h-4 w-32 animate-pulse rounded bg-surface-3" />
        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-28 animate-pulse rounded-lg bg-surface-3" />
          ))}
        </div>
      </div>
    </div>
  );
}
