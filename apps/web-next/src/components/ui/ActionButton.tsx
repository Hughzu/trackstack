import type { JSX } from 'solid-js'

export type ActionButtonProps = {
  type?: 'button' | 'submit'
  disabled?: boolean
  busy?: boolean
  tone?: 'primary' | 'ghost'
  block?: boolean
  onClick?: JSX.EventHandlerUnion<HTMLButtonElement, MouseEvent>
  children: JSX.Element
}

export function ActionButton(props: ActionButtonProps) {
  const toneClass =
    props.tone === 'ghost'
      ? 'border border-border bg-panel text-text-main hover:border-accent/40 hover:text-accent'
      : 'bg-accent text-background hover:opacity-90'

  const blockClass = props.block ? 'w-full' : ''

  return (
    <button
      type={props.type ?? 'button'}
      disabled={props.disabled || props.busy}
      onClick={props.onClick}
      class={`inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-60 ${toneClass} ${blockClass}`}
    >
      {props.children}
    </button>
  )
}

export function ActionGroup(props: { children: JSX.Element }) {
  return <div class="flex gap-2 w-full pt-2 pb-6">{props.children}</div>
}

import { Portal } from 'solid-js/web'

export function FloatingActionGroup(props: { children: JSX.Element }) {
  return (
    <Portal>
      <div class="fixed bottom-6 left-1/2 z-50 flex w-[calc(100%-2rem)] max-w-sm -translate-x-1/2 gap-2 drop-shadow-2xl">
        {props.children}
      </div>
    </Portal>
  )
}


export function IconButton(props: { icon: JSX.Element, textDanger?: boolean, onClick?: JSX.EventHandlerUnion<HTMLButtonElement, MouseEvent> }) {
  const c = props.textDanger 
   ? "h-8 w-8 rounded-full flex flex-shrink-0 items-center justify-center text-text-muted transition-colors hover:text-danger hover:bg-danger/10 outline-none focus-visible:ring-2 focus-visible:ring-accent"
   : "h-8 w-8 rounded-full flex flex-shrink-0 items-center justify-center text-text-muted transition-colors hover:text-text-main hover:bg-surface outline-none focus-visible:ring-2 focus-visible:ring-accent"
  
   return (
    <button type="button" class={c} onClick={props.onClick}>
      {props.icon}
    </button>
  )
}

