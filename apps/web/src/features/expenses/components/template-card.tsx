import { For } from 'solid-js'

import { ActionButton, IconButton } from '../../../components/ui/ActionButton'
import { ChoiceCards } from '../../../components/ui/ChoiceCards'
import { FormFieldRow } from '../../../components/ui/Form'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../../components/ui/List'
import { MetricField } from '../../../components/ui/MetricField'
import { Panel } from '../../../components/ui/Panel'
import { Pill } from '../../../components/ui/Pill'
import { TextField } from '../../../components/ui/TextField'
import { TrashIcon } from '../../../components/ui/TrashIcon'
import type { ExpenseTemplate } from '../../../core/api/types'
import { createExpenseCategoryChoices, getExpenseCategoryMeta, type ExpenseCategoryId } from '../display'

type TemplateCardProps = {
  kind: 'checklist' | 'recurring'
  title: string
  description: string
  badgeTone: 'danger' | 'success'
  badgeLabel: string
  inputTitle: string
  inputAmount: string
  inputCategory: ExpenseCategoryId
  templates: ExpenseTemplate[]
  busy?: boolean
  submitLabel: string
  emptyMessage: string
  helperCopy: string
  titlePlaceholder: string
  onTitleInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
  onAmountInput: (event: InputEvent & { currentTarget: HTMLInputElement }) => void
  onCategoryChange: (value: ExpenseCategoryId) => void
  onSubmit: () => void
  onDelete: (template: ExpenseTemplate) => void
}

export function TemplateCard(props: TemplateCardProps) {
  const choices = createExpenseCategoryChoices()

  return (
    <Panel title={props.title} description={<Pill tone={props.badgeTone}>{props.badgeLabel}</Pill>}>
      <div class="flex flex-col gap-4">
        <div class="space-y-3 border-b border-border/50 pb-4">
          <div class="flex items-center justify-between gap-3">
            <div>
              <div class="text-sm font-bold tracking-tight text-text-main">{props.description}</div>
              <div class="mt-1 text-xs text-text-muted">{props.helperCopy}</div>
            </div>
            <Pill tone="neutral">{props.templates.length} item{props.templates.length === 1 ? '' : 's'}</Pill>
          </div>

          <FormFieldRow>
            <TextField
              id={`${props.kind}-title`}
              name="title"
              label="Title"
              value={props.inputTitle}
              placeholder={props.titlePlaceholder}
              onInput={props.onTitleInput}
            />

            <MetricField
              id={`${props.kind}-amount`}
              name="amount"
              label="Amount"
              unit="EUR"
              value={props.inputAmount}
              min="0"
              step="0.01"
              inputMode="decimal"
              onInput={props.onAmountInput}
            />
          </FormFieldRow>

          <ChoiceCards
            options={choices.map((choice) => ({
              value: choice.value,
              title: choice.title,
              description: choice.description,
              tone: choice.tone === 'danger' ? 'danger' : choice.tone === 'warning' ? 'warning' : 'success',
            }))}
            value={props.inputCategory}
            onChange={(value) => props.onCategoryChange(value as ExpenseCategoryId)}
          />

          <div class="flex justify-end">
            <ActionButton busy={props.busy} onClick={props.onSubmit}>
              {props.busy ? 'Adding...' : props.submitLabel}
            </ActionButton>
          </div>
        </div>

        <List variant="flush" emptyMessage={props.emptyMessage}>
          <For each={props.templates}>
            {(item) => (
              <ListItem
                id={item.id}
                title={item.title}
                subtitle={
                  <ListMeta>
                    <Pill tone={getExpenseCategoryMeta(item.category).tone} size="sm">{getExpenseCategoryMeta(item.category).compactLabel}</Pill>
                    <ListMetaDivider />
                    <span>{props.kind === 'checklist' ? 'Checked off into real expenses later' : 'Preloaded into every new monthly sheet'}</span>
                  </ListMeta>
                }
                value={item.amount.toLocaleString('en-IE', { style: 'currency', currency: 'EUR', minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                valueStyle="mono"
                action={
                  <IconButton
                    ariaLabel={`Delete ${item.title}`}
                    textDanger
                    onClick={() => props.onDelete(item)}
                    icon={<TrashIcon />}
                  />
                }
              />
            )}
          </For>
        </List>
      </div>
    </Panel>
  )
}
