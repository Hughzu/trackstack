import { DataCell, DataRow } from '../../../components/ui/DataRow'
import { List, SkeletonListItem } from '../../../components/ui/List'
import { SkeletonAmountHeroField } from '../../../components/ui/AmountHeroField'
import { SkeletonMetricField } from '../../../components/ui/MetricField'
import { SkeletonPanel } from '../../../components/ui/Skeleton'
import { SkeletonStat } from '../../../components/ui/Stat'

export function HeatDashboardSkeleton() {
  return (
    <>
      <SkeletonPanel titleVariant="md">
        <DataRow variant="header-spaced">
          <SkeletonStat variant="hero" labelPosition="bottom" />
          <SkeletonStat variant="mono" align="right" />
        </DataRow>
        <DataRow variant="footer">
          <SkeletonStat variant="inline" />
          <SkeletonStat variant="inline" />
        </DataRow>
        <DataRow variant="divided">
          <DataCell flex><SkeletonStat variant="lg" align="center" /></DataCell>
          <DataCell flex><SkeletonStat variant="lg" align="center" /></DataCell>
        </DataRow>
        <DataRow variant="footer">
          <SkeletonStat variant="inline" />
          <SkeletonStat variant="inline" />
        </DataRow>
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

export function HeatFormSkeleton() {
  return (
    <>
      <SkeletonAmountHeroField />
      <SkeletonPanel titleVariant="md">
        <SkeletonMetricField />
        <SkeletonMetricField />
        <SkeletonMetricField />
      </SkeletonPanel>
    </>
  )
}
