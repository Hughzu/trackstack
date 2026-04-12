import type { JSX } from 'solid-js'

// No local SkeletonProps type needed.

export function SkeletonPanel(props: { titleVariant?: 'sm' | 'md' | 'lg', children?: JSX.Element }) {
  let w = 'w-16'
  switch(props.titleVariant) {
    case 'sm': w = 'w-12'; break;
    case 'md': w = 'w-16'; break;
    case 'lg': w = 'w-20'; break;
  }

  return (
    <div class="relative flex h-full flex-col overflow-hidden rounded-2xl border border-border/50 bg-surface/60 p-5 shadow-sm animate-pulse transition-all duration-300">
      <div class="mb-4">
        <SkeletonText width={w} height="h-4" />
      </div>
      <div class="flex flex-1 flex-col gap-3">
        <div class="flex flex-1 flex-col gap-3 text-sm leading-6">
          {props.children}
        </div>
      </div>
    </div>
  )
}

type SkeletonPrimitiveProps = {
  width?: string
  height?: string
  class?: string
}

export function SkeletonText(props: SkeletonPrimitiveProps) {
  const base = `${props.height || 'h-3'} ${props.width || 'w-16'} ${props.class || ''}`
  return (
    <div class={`animate-pulse rounded bg-surface border border-border/20 ${base}`} />
  )
}

export function SkeletonBlock(props: SkeletonPrimitiveProps) {
  const base = `${props.height || 'h-8'} ${props.width || 'w-full'} ${props.class || ''}`
  return (
    <div class={`animate-pulse rounded-xl bg-surface/80 border border-border/30 ${base}`} />
  )
}

export function SkeletonProgressBar() {
  return <SkeletonBlock class="mb-1 h-1.5 w-full rounded-full" />
}
