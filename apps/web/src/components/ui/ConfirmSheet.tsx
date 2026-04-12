import { Show, createEffect, onCleanup } from 'solid-js'
import { Portal } from 'solid-js/web'

import { ActionButton } from './ActionButton'
import { FormActions } from './Form'
import { Pill } from './Pill'

type ConfirmSheetProps = {
  open: boolean
  title: string
  description: string
  confirmLabel: string
  cancelLabel?: string
  busy?: boolean
  confirmTone?: 'primary' | 'danger'
  eyebrow?: string
  onConfirm: () => void | Promise<void>
  onCancel: () => void
}

export function ConfirmSheet(props: ConfirmSheetProps) {
  createEffect(() => {
    if (!props.open) return

    const handleKeydown = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && !props.busy) {
        props.onCancel()
      }
    }

    document.body.style.overflow = 'hidden'
    window.addEventListener('keydown', handleKeydown)

    onCleanup(() => {
      document.body.style.overflow = ''
      window.removeEventListener('keydown', handleKeydown)
    })
  })

  return (
    <Show when={props.open}>
      <Portal>
        <div class="fixed inset-0 z-50">
          <button
            type="button"
            aria-label="Close confirmation"
            class="absolute inset-0 bg-background/80 backdrop-blur-sm"
            disabled={props.busy}
            onClick={props.onCancel}
          />

          <div class="absolute inset-x-0 bottom-0 mx-auto w-full max-w-3xl px-4 pb-4 sm:pb-6">
            <section class="animate-in slide-in-from-bottom-6 duration-300 rounded-[1.75rem] border border-border bg-surface p-5 shadow-2xl shadow-black/40 sm:ml-auto sm:max-w-md">
              <div class="mb-4 flex items-start justify-between gap-3">
                <div class="space-y-2">
                  <Pill tone={props.confirmTone === 'danger' ? 'danger' : 'accent'}>
                    {props.eyebrow ?? 'Confirm action'}
                  </Pill>
                  <div>
                    <h2 class="text-lg font-bold tracking-tight text-text-main">{props.title}</h2>
                    <p class="mt-2 text-sm leading-6 text-text-muted">{props.description}</p>
                  </div>
                </div>
              </div>

              <FormActions>
                <ActionButton tone="ghost" onClick={props.onCancel} disabled={props.busy}>{props.cancelLabel ?? 'Cancel'}</ActionButton>
                <ActionButton tone={props.confirmTone === 'danger' ? 'danger' : 'primary'} busy={props.busy} onClick={props.onConfirm}>
                  {props.confirmLabel}
                </ActionButton>
              </FormActions>
            </section>
          </div>
        </div>
      </Portal>
    </Show>
  )
}
