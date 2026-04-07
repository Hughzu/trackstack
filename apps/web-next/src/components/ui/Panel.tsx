import type { JSX } from 'solid-js'
import { createSignal, Show } from 'solid-js'
import { A } from '@solidjs/router'

type PanelProps = {
  title: string
  description?: string | JSX.Element
  headerAction?: JSX.Element
  href?: string
  collapsibleId?: string
  children?: JSX.Element
}

const frameClass =
  'relative flex h-full flex-col overflow-hidden rounded-2xl border border-border bg-surface p-5 text-text-main transition-all duration-300'

function PanelHeader(props: { title: string, description?: string | JSX.Element, headerAction?: JSX.Element }) {
  return (
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xs font-bold uppercase tracking-widest text-accent">{props.title}</h2>
      <div class="flex items-center gap-2 text-[0.68rem] text-text-muted">
        {props.description ? <div>{props.description}</div> : null}
        {props.headerAction}
      </div>
    </div>
  )
}

function PanelBody(props: PanelProps) {
  return (
    <>
      <PanelHeader title={props.title} description={props.description} headerAction={props.headerAction} />
      <div class="flex flex-1 flex-col gap-3 text-sm leading-6 text-text-main">
        {props.children}
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

  if (props.collapsibleId) {
    const isClient = typeof localStorage !== 'undefined'
    const stored = isClient ? localStorage.getItem(`panel_${props.collapsibleId}`) : null
    const [isOpen, setIsOpen] = createSignal(stored !== 'false')

    const toggle = () => {
      const next = !isOpen()
      setIsOpen(next)
      if (isClient) localStorage.setItem(`panel_${props.collapsibleId}`, String(next))
    }

    return (
      <section class={`${frameClass} shadow-lg shadow-black/20`}>
        <button type="button" class="flex flex-col text-left outline-none" onClick={toggle}>
          <div class="flex items-center justify-between w-full mb-4">
            <h2 class="text-xs font-bold uppercase tracking-widest text-accent flex items-center gap-2">
              {props.title}
              <svg class={`h-3 w-3 transition-transform ${isOpen() ? 'rotate-180' : ''}`} viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
              </svg>
            </h2>
            <div class="flex items-center gap-2 text-[0.68rem] text-text-muted">
              {props.description ? <div>{props.description}</div> : null}
              {props.headerAction}
            </div>
          </div>
        </button>
        <Show when={isOpen()}>
          <div class="flex flex-1 flex-col gap-3 text-sm leading-6 text-text-main animate-in fade-in slide-in-from-top-2 duration-300">
            {props.children}
          </div>
        </Show>
      </section>
    )
  }

  return (
    <section class={`${frameClass} shadow-lg shadow-black/20`}>
      <PanelBody {...props} />
    </section>
  )
}
