import { AppShell } from './AppShell'
import { SkeletonPanel, SkeletonText, SkeletonBlock } from './Skeleton'

export function RouteStatus() {
  return (
    <AppShell currentDomain="">
      <div class="space-y-8 animate-in fade-in duration-500">
        <header>
          <SkeletonText class="h-8 w-1/3 mb-2" />
          <SkeletonText class="h-4 w-1/4" />
        </header>

        <SkeletonPanel class="p-5 space-y-6">
          <header class="mb-4">
            <SkeletonText class="h-3 w-16 bg-accent/20" />
          </header>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
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
    </AppShell>
  )
}
