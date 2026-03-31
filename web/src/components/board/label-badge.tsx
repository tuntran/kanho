import type { Label } from "@/types/board";

interface LabelBadgeProps {
  label: Label;
}

export function LabelBadge({ label }: LabelBadgeProps) {
  return (
    <span
      className="inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-medium"
      style={{
        backgroundColor: label.color + "20",
        color: label.color,
      }}
    >
      {label.name}
    </span>
  );
}
