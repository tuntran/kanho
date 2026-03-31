export function LabelChip({ name, color }: { name: string; color: string }) {
  return (
    <span
      className="inline-flex h-5 items-center rounded px-1.5 text-[11px] font-medium"
      style={{
        backgroundColor: `${color}20`,
        color,
      }}
    >
      {name}
    </span>
  );
}
