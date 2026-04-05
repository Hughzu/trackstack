import type { JSX } from 'solid-js'

type ActionButtonProps = {
  type?: 'button' | 'submit'
  disabled?: boolean
  busy?: boolean
  tone?: 'primary' | 'ghost'
  onClick?: JSX.EventHandlerUnion<HTMLButtonElement, MouseEvent>
  children: JSX.Element
}

export function ActionButton(props: ActionButtonProps) {
  const toneClass =
    props.tone === 'ghost'
      ? 'border border-border bg-panel text-text-main hover:border-accent/40 hover:text-accent'
      : 'bg-accent text-background hover:opacity-90'

  return (
    <button
      type={props.type ?? 'button'}
      disabled={props.disabled || props.busy}
      onClick={props.onClick}
      class={`inline-flex items-center justify-center rounded-xl px-4 py-2 text-sm font-semibold transition disabled:cursor-not-allowed disabled:opacity-60 ${toneClass}`}
    >
      {props.children}
    </button>
  )
}
