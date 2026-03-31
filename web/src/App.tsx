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

function DashboardPlaceholder() {
  return (
    <div className="flex h-screen items-center justify-center">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-text-1">Kanho</h1>
        <p className="mt-2 text-text-2">Kanban project management</p>
      </div>
    </div>
  );
}

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
          <Route path="/*" element={<DashboardPlaceholder />} />
        </Route>
      </Routes>
      <SearchDialog isOpen={searchOpen} onClose={closeSearch} />
      <ToastProvider />
    </>
  );
}
