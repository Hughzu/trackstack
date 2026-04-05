type RouteStatusProps = {
  title: string
  description: string
}

export function RouteStatus(props: RouteStatusProps) {
  return (
    <div class="flex min-h-screen items-center justify-center bg-background px-4 text-text-main">
      <section class="w-full max-w-md rounded-3xl border border-border bg-surface p-8 shadow-xl shadow-black/20">
        <div class="mb-3 text-xs font-bold uppercase tracking-[0.28em] text-accent">TrackStack</div>
        <h1 class="text-2xl font-bold tracking-tight">{props.title}</h1>
        <p class="mt-3 text-sm leading-6 text-text-muted">{props.description}</p>
      </section>
    </div>
  )
}
