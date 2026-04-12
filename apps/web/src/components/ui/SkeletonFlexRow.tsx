import type { JSX } from 'solid-js'
import { SkeletonText } from './Skeleton'

export function SkeletonFlexRow(props: { children: JSX.Element }) {
  return <div class="mb-4 flex justify-between">{props.children}</div>
}

export function SkeletonTextInline(props: { variant?: 'sm' | 'md' | 'lg' | 'hero' }) {
  let w = 'w-16'
  switch(props.variant) {
    case 'sm': w = 'w-10'; break;
    case 'md': w = 'w-16'; break;
    case 'lg': w = 'w-24'; break;
    case 'hero': w = 'w-32'; break;
  }
  return <SkeletonText width={w} />
}
