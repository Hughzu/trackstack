import type { JSX } from 'solid-js'

export type DataRowProps = {
  children: JSX.Element
  variant?: 'header' | 'header-end' | 'header-spaced' | 'divided' | 'footer'
}

export function DataRow(props: DataRowProps) {
  let containerClass = 'flex justify-between'

  switch (props.variant) {
    case 'header':
      containerClass += ' mb-2 items-start'
      break
    case 'header-end':
      containerClass += ' mb-4 items-end'
      break
    case 'header-spaced':
      containerClass += ' mb-6 items-start'
      break
    case 'divided':
      containerClass += ' divide-x divide-border/50 text-center'
      break
    case 'footer':
      containerClass += ' border-t border-border/50 pt-3 text-xs text-text-muted'
      break
  }

  return <div class={containerClass}>{props.children}</div>
}

export function DataCell(props: { children: JSX.Element, flex?: boolean }) {
  return <div class={props.flex ? 'flex flex-1 flex-col px-1 items-center' : ''}>{props.children}</div>
}
