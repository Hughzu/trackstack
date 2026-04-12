import { ConfirmSheet } from './ConfirmSheet'

type DeleteConfirmSheetProps = {
  open: boolean
  title: string
  description: string
  eyebrow?: string
  confirmLabel?: string
  busy?: boolean
  onCancel: () => void
  onConfirm: () => void | Promise<void>
}

export function DeleteConfirmSheet(props: DeleteConfirmSheetProps) {
  return (
    <ConfirmSheet
      open={props.open}
      title={props.title}
      description={props.description}
      eyebrow={props.eyebrow}
      confirmLabel={props.confirmLabel ?? 'Delete'}
      confirmTone="danger"
      busy={props.busy}
      onCancel={props.onCancel}
      onConfirm={props.onConfirm}
    />
  )
}
