import { SkeletonBudgetBreakdown } from '../../../components/ui/BudgetBreakdown'
import { List, SkeletonListItem } from '../../../components/ui/List'
import { SkeletonPanel } from '../../../components/ui/Skeleton'

export function DashboardSkeleton() {
  return (
    <>
      <SkeletonBudgetBreakdown />
      <SkeletonPanel titleVariant="md">
        <List variant="flush">
          <SkeletonListItem />
          <SkeletonListItem />
        </List>
      </SkeletonPanel>
    </>
  )
}
