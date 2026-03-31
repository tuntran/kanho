import { useBoardStore } from "@/stores/board-store";

export function LiveDot() {
  const connected = useBoardStore((s) => s.wsConnected);

  return (
    <span
      className={`inline-block h-2 w-2 rounded-full ${
        connected
          ? "bg-success animate-[pulse-dot_2s_infinite]"
          : "bg-text-3"
      }`}
      title={connected ? "Connected" : "Disconnected"}
    />
  );
}
