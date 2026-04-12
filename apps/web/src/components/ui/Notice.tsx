type NoticeProps = {
  tone?: 'error' | 'info'
  message: string
}

export function Notice(props: NoticeProps) {
  const toneClass = props.tone === 'error' ? 'border-red-400/30 bg-red-500/10 text-red-100' : 'border-border bg-panel text-text-muted'

  return <div class={`rounded-xl border px-4 py-3 text-sm leading-6 ${toneClass}`}>{props.message}</div>
}
