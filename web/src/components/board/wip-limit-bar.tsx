export function WipLimitBar({
  count,
  limit,
}: {
  count: number;
  limit: number;
}) {
  const ratio = count / limit;
  const pct = Math.min(ratio * 100, 100);

  let barColor = "bg-success";
  if (ratio >= 1) barColor = "bg-danger";
  else if (ratio >= 0.8) barColor = "bg-warning";

  return (
    <div className="h-[3px] w-full rounded-sm bg-border">
      <div
        className={`h-full rounded-sm transition-all ${barColor} ${ratio >= 1 ? "animate-pulse" : ""}`}
        style={{ width: `${pct}%` }}
      />
    </div>
  );
}
