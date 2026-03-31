import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery, useQueries } from "@tanstack/react-query";
import {
  getWorkspace,
  listMembers,
  workspaceKeys,
} from "@/api/workspace-client";
import {
  listProjects,
  listProjectBoards,
  projectKeys,
} from "@/api/project-client";
import { WorkspaceDashboardHeader } from "./workspace-dashboard-header";
import { WorkspaceBoardsGrid } from "./workspace-boards-grid";
import { CreateProjectModal } from "./create-project-modal";

export function WorkspaceDashboardPage() {
  const { workspaceSlug } = useParams<{ workspaceSlug: string }>();
  const navigate = useNavigate();
  const [modalOpen, setModalOpen] = useState(false);

  const slug = workspaceSlug!;

  const workspaceQuery = useQuery({
    queryKey: workspaceKeys.detail(slug),
    queryFn: () => getWorkspace(slug),
  });

  const membersQuery = useQuery({
    queryKey: workspaceKeys.members(slug),
    queryFn: () => listMembers(slug),
  });

  const projectsQuery = useQuery({
    queryKey: projectKeys.list(slug),
    queryFn: () => listProjects(slug),
  });

  const projects = projectsQuery.data ?? [];

  const boardQueries = useQueries({
    queries: projects.map((project) => ({
      queryKey: projectKeys.boards(slug, project.id),
      queryFn: () => listProjectBoards(slug, project.id),
      enabled: projects.length > 0,
    })),
  });

  // Flatten boards with project key
  const allBoards = projects.flatMap((project, i) => {
    const boards = boardQueries[i]?.data ?? [];
    return boards.map((board) => ({ ...board, projectKey: project.key }));
  });

  // Error: workspace not found
  if (workspaceQuery.error) {
    navigate("/", { replace: true });
    return null;
  }

  // Loading state
  if (workspaceQuery.isLoading || membersQuery.isLoading) {
    return <DashboardSkeleton />;
  }

  const workspace = workspaceQuery.data;
  if (!workspace) return null;

  return (
    <div className="min-h-screen bg-surface-1">
      <div className="mx-auto max-w-4xl px-6 py-8">
        <WorkspaceDashboardHeader
          workspace={workspace}
          members={membersQuery.data ?? []}
          boardCount={allBoards.length}
        />
        <WorkspaceBoardsGrid
          workspaceSlug={slug}
          boards={allBoards}
          onCreateProject={() => setModalOpen(true)}
        />
        <CreateProjectModal
          workspaceSlug={slug}
          isOpen={modalOpen}
          onClose={() => setModalOpen(false)}
        />
      </div>
    </div>
  );
}

function DashboardSkeleton() {
  return (
    <div className="min-h-screen bg-surface-1">
      <div className="mx-auto max-w-4xl px-6 py-8">
        <div className="flex items-start gap-4">
          <div className="h-12 w-12 animate-pulse rounded-lg bg-surface-3" />
          <div>
            <div className="h-6 w-48 animate-pulse rounded bg-surface-3" />
            <div className="mt-2 h-4 w-32 animate-pulse rounded bg-surface-3" />
            <div className="mt-2 h-3 w-24 animate-pulse rounded bg-surface-3" />
          </div>
        </div>
        <div className="mt-8 h-5 w-20 animate-pulse rounded bg-surface-3" />
        <div className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-20 animate-pulse rounded-lg bg-surface-3"
            />
          ))}
        </div>
      </div>
    </div>
  );
}
