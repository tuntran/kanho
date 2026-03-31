# 2026-03-31 — Workspace Dashboard & Onboarding Step 3

## Commits
- `288e825` feat(workspace): add workspace dashboard with project/board management
- `01fa577` feat(onboarding): add project creation step and workspace dashboard route

## What Changed

### New route: `/w/:workspaceSlug`
- `workspace-dashboard-page.tsx` — parallel fetch of workspace + projects + boards on mount
- `workspace-dashboard-header.tsx` — workspace avatar, member count, board count
- `workspace-boards-grid.tsx` — board cards grid with empty state CTA

### Onboarding: step 3 added
- `create-workspace-page.tsx` — extended wizard with "Create first project" step
- `create-project-modal.tsx` — project creation form with auto-generated project key
- `App.tsx` — registered `/w/:workspaceSlug` route (fixed "No routes matched" bug after workspace creation)

### API layer
- `project-client.ts` (new) — `listProjects`, `createProject`, `listProjectBoards`, `projectKeys`
- `workspace-client.ts` — added `getWorkspace`, `listMembers`, `WorkspaceMember` interface

## Key Decisions
- Parallel data fetching on dashboard load (workspace + members + projects in one `Promise.all`)
- Auto-gen project key from project name (slugified, uppercased prefix) — reduces friction in onboarding
- Workspace dashboard scoped under `/w/:slug` namespace to support multi-workspace switching later

## Bug Fixed
- Missing route `/w/:slug` caused hard redirect failure after workspace creation; added to `App.tsx`

## TypeScript
All checks: PASS
