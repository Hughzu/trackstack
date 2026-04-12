import { SkeletonActionChipRow } from '../../../components/ui/ActionChipRow'
import { SkeletonPanel } from '../../../components/ui/Skeleton'
import { SkeletonAmountHeroField } from '../../../components/ui/AmountHeroField'
import { List, SkeletonListItem } from '../../../components/ui/List'
import { SkeletonMetricField } from '../../../components/ui/MetricField'
import { DataCell, DataRow } from '../../../components/ui/DataRow'
import { SkeletonProgressBar } from '../../../components/ui/Skeleton'
import { SkeletonStat } from '../../../components/ui/Stat'

export function CaloriesDashboardSkeleton() {
  return (
    <>
      <SkeletonPanel titleVariant="md">
        <DataRow variant="header-end">
          <SkeletonStat variant="hero" />
          <SkeletonStat variant="md" align="right" />
        </DataRow>
        <SkeletonProgressBar />
        <DataRow variant="footer">
          <SkeletonStat variant="inline" />
          <SkeletonStat variant="inline" />
        </DataRow>
        <DataRow variant="divided">
          <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
          <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
          <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
        </DataRow>
      </SkeletonPanel>

      <SkeletonPanel titleVariant="sm">
        <SkeletonActionChipRow />
      </SkeletonPanel>

      <SkeletonPanel titleVariant="md">
        <List variant="flush">
          <SkeletonListItem />
          <SkeletonListItem />
          <SkeletonListItem />
        </List>
      </SkeletonPanel>
    </>
  )
}

export function CaloriesSettingsSkeleton() {
  return (
    <>
      <SkeletonAmountHeroField />
      <SkeletonPanel titleVariant="md">
        <SkeletonMetricField variant="inline" />
        <SkeletonMetricField variant="inline" />
        <SkeletonMetricField variant="inline" />
      </SkeletonPanel>
    </>
  )
}
