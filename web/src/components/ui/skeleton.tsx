export function CardSkeleton() {
  return (
    <div className="animate-pulse rounded-card border border-border bg-surface p-3">
      <div className="flex items-center justify-between">
        <div className="h-4 w-12 rounded bg-surface-2" />
        <div className="h-3 w-16 rounded bg-surface-2" />
      </div>
      <div className="mt-2 h-4 w-3/4 rounded bg-surface-2" />
      <div className="mt-3 flex gap-1">
        <div className="h-5 w-10 rounded bg-surface-2" />
        <div className="h-5 w-14 rounded bg-surface-2" />
      </div>
    </div>
  );
}

export function DrawerSkeleton() {
  return (
    <div className="animate-pulse space-y-4 p-5">
      <div className="h-6 w-2/3 rounded bg-surface-2" />
      <div className="h-4 w-1/3 rounded bg-surface-2" />
      <div className="mt-6 space-y-2">
        <div className="h-3 w-full rounded bg-surface-2" />
        <div className="h-3 w-5/6 rounded bg-surface-2" />
        <div className="h-3 w-4/6 rounded bg-surface-2" />
      </div>
      <div className="mt-6 space-y-3">
        <div className="h-8 w-full rounded bg-surface-2" />
        <div className="h-8 w-full rounded bg-surface-2" />
      </div>
    </div>
  );
}
