import type { JSX } from 'solid-js'

export function FormStack(props: { children: JSX.Element; onSubmit?: (e: SubmitEvent) => void }) {
  return (
    <form class="flex flex-col gap-4" onSubmit={props.onSubmit}>
      {props.children}
    </form>
  )
}

export function FormActions(props: { children: JSX.Element }) {
  return <div class="flex items-center justify-end gap-3">{props.children}</div>
}
