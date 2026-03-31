import { Toaster } from "sonner";

export function ToastProvider() {
  return (
    <Toaster
      position="bottom-right"
      theme="dark"
      toastOptions={{
        className:
          "bg-surface border border-border text-text-1 text-sm shadow-lg",
        duration: 4000,
      }}
    />
  );
}

export { toast } from "sonner";
