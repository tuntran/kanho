import { useState, useCallback } from "react";
import { Routes, Route } from "react-router-dom";
import { useSilentRefresh } from "@/hooks/use-silent-refresh";
import { useCmdK } from "@/hooks/use-cmd-k";
import { ProtectedRoute } from "@/components/protected-route";
import { LoginPage } from "@/pages/auth/login-page";
import { RegisterPage } from "@/pages/auth/register-page";
import { BoardPage } from "@/pages/board/board-page";
import { CardDrawer } from "@/components/drawer/card-drawer";
import { SearchDialog } from "@/components/search/search-dialog";
import { ToastProvider } from "@/components/ui/toast";
import { WorkspaceSelectorPage } from "@/pages/workspace/workspace-selector-page";
import { CreateWorkspacePage } from "@/pages/onboarding/create-workspace-page";
import { WorkspaceDashboardPage } from "@/pages/workspace/workspace-dashboard-page";

export default function App() {
  useSilentRefresh();
  const [searchOpen, setSearchOpen] = useState(false);

  const openSearch = useCallback(() => setSearchOpen(true), []);
  const closeSearch = useCallback(() => setSearchOpen(false), []);

  useCmdK(openSearch);

  return (
    <>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/w/:workspaceSlug/b/:boardId" element={<BoardPage />}>
            <Route path="cards/:cardNumber" element={<CardDrawer />} />
          </Route>
          <Route path="/w/:workspaceSlug" element={<WorkspaceDashboardPage />} />
          <Route path="/" element={<WorkspaceSelectorPage />} />
          <Route path="/onboarding/create-workspace" element={<CreateWorkspacePage />} />
        </Route>
      </Routes>
      <SearchDialog isOpen={searchOpen} onClose={closeSearch} />
      <ToastProvider />
    </>
  );
}
