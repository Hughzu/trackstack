import type { JSX } from 'solid-js'

type ContentDeckProps = {
  children: JSX.Element
}

export function ContentDeck(props: ContentDeckProps) {
  return <section class="grid gap-5 md:grid-cols-2">{props.children}</section>
}
