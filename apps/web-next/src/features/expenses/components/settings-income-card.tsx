import { For } from 'solid-js'

import { AmountHeroField } from '../../../components/ui/AmountHeroField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { ProgressBar } from '../../../components/ui/ProgressBar'
import type { ExpenseBudgetPreviewItem, ExpenseRatioSummary } from '../display'

export function SettingsIncomeCard(props: {
  income: string
  ratioSummary: ExpenseRatioSummary
  budgetPreview: ExpenseBudgetPreviewItem[]
  onIncomeInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
}) {
  return (
    <Panel title="Income and split" description={<Pill tone={props.ratioSummary.totalTone}>{props.ratioSummary.totalLabel}</Pill>}>
      <AmountHeroField
        inputId="expense-income"
        name="income"
        label="Monthly net income"
        unit="EUR"
        value={props.income}
        placeholder="0.00"
        badge={`${Math.round(props.ratioSummary.total)}% allocated`}
        badgeTone={props.ratioSummary.totalTone === 'danger' ? 'danger' : props.ratioSummary.totalTone === 'warning' ? 'warning' : 'success'}
        inputMode="decimal"
        min="0"
        step="0.01"
        onInput={props.onIncomeInput}
      />

      <div class="border-t border-border/50 pt-4">
        <div class="flex items-start justify-between gap-3">
          <div>
            <div class="text-sm font-semibold text-text-main">Budget split</div>
            <p class="mt-1 text-xs leading-5 text-text-muted">
              {props.ratioSummary.total === 100
                ? 'Clean split. Nothing extra hanging around.'
                : props.ratioSummary.total > 100
                  ? `You are over by ${Math.round(props.ratioSummary.total - 100)}%.`
                  : `${Math.round(100 - props.ratioSummary.total)}% is still unassigned.`}
            </p>
          </div>
        </div>
      </div>

      <ProgressBar segments={props.budgetPreview.map((item) => ({ percent: item.percent, color: item.barColor }))} />

      <div class="divide-y divide-border/40 rounded-2xl border border-border/50 bg-panel/50">
        <For each={props.budgetPreview}>
          {(item) => (
            <div class="flex items-center justify-between gap-3 px-4 py-3">
              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <span class="text-sm font-semibold text-text-main">{item.label}</span>
                  <Pill tone={item.tone} size="sm">{Math.round(item.percent)}%</Pill>
                </div>
                <p class="mt-1 text-xs leading-5 text-text-muted">{item.description}</p>
              </div>
              <div class="text-right text-sm font-mono text-text-main">{item.value}</div>
            </div>
          )}
        </For>
      </div>
    </Panel>
  )
}
