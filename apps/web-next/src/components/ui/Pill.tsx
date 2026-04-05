import type { JSX } from 'solid-js'

type PillProps = {
  tone?: 'accent' | 'neutral'
  children: JSX.Element
}

export function Pill(props: PillProps) {
  const className =
    props.tone === 'neutral'
      ? 'inline-flex w-fit rounded-full border border-border bg-panel px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-text-muted'
      : 'inline-flex w-fit rounded-full border border-accent/20 bg-accent/10 px-3 py-1 text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-accent'

  return <span class={className}>{props.children}</span>
}
