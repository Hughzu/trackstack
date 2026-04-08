import type { JSX } from 'solid-js'
import { A } from '@solidjs/router'

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

export function FormFieldRow(props: { children: JSX.Element }) {
  return <div class="grid gap-2.5 sm:grid-cols-[1.4fr_0.9fr]">{props.children}</div>
}

export function FormSection(props: { children: JSX.Element }) {
  return <div class="flex flex-col gap-3">{props.children}</div>
}

export function FormBackLink(props: { href: string, children: JSX.Element }) {
  return (
    <A
      href={props.href}
      class="inline-flex items-center rounded-full border border-border bg-panel px-3 py-1.5 text-[0.68rem] font-semibold uppercase tracking-[0.24em] text-text-muted transition hover:border-accent/40 hover:text-accent"
    >
      {props.children}
    </A>
  )
}
