type SectionTitleProps = {
  eyebrow: string
  title: string
  description?: string
}

export function SectionTitle(props: SectionTitleProps) {
  return (
    <section class="space-y-2 border-b border-border/70 pb-3">
      <div class="text-xs font-bold uppercase tracking-widest text-accent">{props.eyebrow}</div>
      <h2 class="text-2xl font-bold tracking-tight text-text-main sm:text-3xl">{props.title}</h2>
      {props.description ? <p class="max-w-3xl text-sm leading-6 text-text-muted sm:text-base">{props.description}</p> : null}
    </section>
  )
}
