import type { Priority } from "@/types/board";

const PRIORITY_CONFIG: Record<
  Priority,
  { label: string; color: string; bg: string }
> = {
  urgent: { label: "Urgent", color: "#EF4444", bg: "rgba(239,68,68,0.12)" },
  high: { label: "High", color: "#F59E0B", bg: "rgba(245,158,11,0.12)" },
  medium: { label: "Medium", color: "#6366F1", bg: "rgba(99,102,241,0.12)" },
  low: { label: "Low", color: "#38BDF8", bg: "rgba(56,189,248,0.12)" },
  none: { label: "None", color: "#475569", bg: "rgba(71,85,105,0.10)" },
};

export function PriorityBadge({ priority }: { priority: Priority }) {
  if (priority === "none") return null;
  const cfg = PRIORITY_CONFIG[priority];
  return (
    <span
      className="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium uppercase"
      style={{ color: cfg.color, backgroundColor: cfg.bg }}
    >
      {cfg.label}
    </span>
  );
}
