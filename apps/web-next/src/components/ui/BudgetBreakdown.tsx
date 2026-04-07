import { For, type JSX } from 'solid-js'

import { DataCell, DataRow } from './DataRow'
import { ProgressBar, type ProgressBarProps } from './ProgressBar'
import { SkeletonPanel, SkeletonProgressBar } from './Skeleton'
import { SkeletonStat, Stat, type StatProps } from './Stat'

export type BudgetBreakdownItem = {
  label: string
  value: string | JSX.Element
  percent?: number
  subtext?: string | JSX.Element
  color?: StatProps['color']
  barColor?: ProgressBarProps['segments'][number]['color']
}

type BudgetBreakdownProps = {
  remaining: string | JSX.Element
  income: string | JSX.Element
  items: BudgetBreakdownItem[]
}

export function BudgetBreakdown(props: BudgetBreakdownProps) {
  return (
    <>
      <DataRow variant="header">
        <Stat label="Remaining" value={props.remaining} variant="lg" />
        <Stat label="Income" value={props.income} variant="mono" align="right" />
      </DataRow>

      <ProgressBar
        segments={props.items.map((item) => ({
          percent: item.percent || 0,
          color: item.barColor || 'main',
        }))}
      />

      <DataRow variant="divided">
        <For each={props.items}>
          {(item) => (
            <DataCell flex>
              <Stat
                label={item.label}
                value={item.value}
                subtext={item.subtext}
                color={item.color || 'main'}
                variant="sm"
                align="center"
              />
            </DataCell>
          )}
        </For>
      </DataRow>
    </>
  )
}

export function SkeletonBudgetBreakdown() {
  return (
    <SkeletonPanel titleVariant="md">
      <DataRow variant="header">
        <SkeletonStat variant="lg" />
        <SkeletonStat variant="mono" align="right" />
      </DataRow>
      <SkeletonProgressBar />
      <DataRow variant="divided">
        <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
        <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
        <DataCell flex><SkeletonStat variant="sm" align="center" /></DataCell>
      </DataRow>
    </SkeletonPanel>
  )
}
