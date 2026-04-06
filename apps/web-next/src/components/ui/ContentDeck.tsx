import type { JSX } from 'solid-js'

type ContentDeckProps = {
  children: JSX.Element
  layout?: 'grid' | 'stacked'
  animate?: boolean
}

export function ContentDeck(props: ContentDeckProps) {
  let baseClass = props.layout === 'stacked' 
    ? 'flex flex-col gap-6' 
    : 'grid gap-5 md:grid-cols-2'
    
  if (props.animate) {
    baseClass += ' animate-in fade-in slide-in-from-bottom-2 duration-700'
  }

  return <section class={baseClass}>{props.children}</section>
}
