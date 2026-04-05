import type { JSX } from 'solid-js'
import { A } from '@solidjs/router'

type PanelProps = {
  title: string
  description?: string
  href?: string
  children?: JSX.Element
}

const frameClass =
  'relative flex h-full flex-col overflow-hidden rounded-2xl border border-border bg-surface p-5 text-text-main transition-all duration-300'

function PanelBody(props: PanelProps) {
  return (
    <>
      <div class="mb-4">
        <h2 class="text-xs font-bold uppercase tracking-widest text-accent">{props.title}</h2>
      </div>
      <div class="flex flex-1 flex-col gap-3">
        {props.description ? <p class="text-sm leading-6 text-text-muted">{props.description}</p> : null}
        <div class="flex flex-1 flex-col gap-3 text-sm leading-6 text-text-main">{props.children}</div>
      </div>
    </>
  )
}

export function Panel(props: PanelProps) {
  if (props.href) {
    return (
      <A href={props.href} class={`${frameClass} hover:border-accent/30 hover:shadow-lg hover:shadow-black/30 hover:ring-1 hover:ring-accent/20`}>
        <PanelBody {...props} />
      </A>
    )
  }

  return (
    <section class={`${frameClass} shadow-lg shadow-black/20`}>
      <PanelBody {...props} />
    </section>
  )
}
