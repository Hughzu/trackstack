import { SkeletonPanel, SkeletonText, SkeletonBlock } from './Skeleton'

export function RouteStatus() {
  return (
    <div class="space-y-8 animate-in fade-in duration-500">
      <header>
        <SkeletonText class="mb-2 h-8 w-1/3" />
        <SkeletonText class="h-4 w-1/4" />
      </header>

      <SkeletonPanel titleVariant="md">
        <div class="mb-4 grid grid-cols-1 gap-4 md:grid-cols-3">
          <SkeletonBlock class="h-24" />
          <SkeletonBlock class="h-24" />
          <SkeletonBlock class="h-24" />
        </div>

        <div class="space-y-3 pt-4">
          <SkeletonText class="h-4 w-full" />
          <SkeletonText class="h-4 w-11/12" />
          <SkeletonText class="h-4 w-4/5" />
        </div>
      </SkeletonPanel>
    </div>
  )
}
