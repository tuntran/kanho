// TODO: view tạm — implement workspace creation form khi có design chính thức
import { useNavigate } from "react-router-dom";

export function CreateWorkspacePage() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-1 p-6">
      <div className="w-full max-w-md text-center">
        <h1 className="text-2xl font-bold text-text-1">Tạo workspace mới</h1>
        <p className="mt-2 text-text-2">
          Tính năng đang được phát triển.
        </p>
        <button
          onClick={() => navigate("/")}
          className="mt-6 text-sm text-text-2 hover:text-text-1"
        >
          ← Quay lại
        </button>
      </div>
    </div>
  );
}
