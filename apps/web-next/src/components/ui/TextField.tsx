import type { JSX } from 'solid-js'

type TextFieldProps = {
  id: string
  name: string
  label: string
  type?: string
  value: string
  autocomplete?: string
  required?: boolean
  onInput: JSX.EventHandlerUnion<HTMLInputElement, InputEvent>
}

export function TextField(props: TextFieldProps) {
  return (
    <label class="flex flex-col gap-2 text-sm text-text-main" for={props.id}>
      <span class="text-xs font-bold uppercase tracking-[0.24em] text-text-muted">{props.label}</span>
      <input
        id={props.id}
        name={props.name}
        type={props.type ?? 'text'}
        value={props.value}
        autocomplete={props.autocomplete}
        required={props.required}
        onInput={props.onInput}
        class="rounded-xl border border-border bg-panel px-4 py-3 text-base text-text-main outline-none transition placeholder:text-text-muted/70 focus:border-accent"
      />
    </label>
  )
}
