import type { JSX } from 'solid-js'

type SkeletonProps = {
  class?: string
  children?: JSX.Element
}

export function SkeletonPanel(props: SkeletonProps) {
  return (
    <div class={`animate-pulse rounded-2xl bg-surface/60 border border-border/50 shadow-sm ${props.class || ''}`}>
      {props.children}
    </div>
  )
}

export function SkeletonText(props: SkeletonProps) {
  return (
    <div class={`animate-pulse rounded bg-surface border border-border/20 ${props.class || ''}`} />
  )
}

export function SkeletonBlock(props: SkeletonProps) {
  return (
    <div class={`animate-pulse rounded-xl bg-surface/80 border border-border/30 ${props.class || ''}`} />
  )
}
