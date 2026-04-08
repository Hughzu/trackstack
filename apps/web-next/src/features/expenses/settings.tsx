import { createEffect, createMemo, createResource, createSignal, For, Show, Suspense } from 'solid-js'

import { ActionButton, IconButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ConfirmSheet } from '../../components/ui/ConfirmSheet'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { FormBackLink } from '../../components/ui/Form'
import { List, ListItem, ListMeta, ListMetaDivider } from '../../components/ui/List'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { Pill } from '../../components/ui/Pill'
import { ProgressBar } from '../../components/ui/ProgressBar'
import { authState } from '../../core/auth/state'
import {
  deleteChecklistTemplate,
  deleteRecurringTemplate,
  readExpensesSettings,
  updateExpensesSettings,
  upsertChecklistTemplate,
  upsertRecurringTemplate,
} from './api/client'
import { getExpenseCategoryMeta } from './display'

type TemplateKind = 'checklist' | 'recurring'
type CategoryId = 'fund' | 'fun' | 'future'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

const categoryOptions: Array<{ id: CategoryId, label: string, description: string }> = [
  { id: 'fund', label: 'Fundamentals', description: 'Bills, rent, insurance, the annoying grown-up stuff.' },
  { id: 'fun', label: 'Fun', description: 'Leisure, meals out, random little dopamine hits.' },
  { id: 'future', label: 'Future', description: 'Savings, investing, buffers, not screwing over future-you.' },
]

const formatEuro = (value: number) =>
  new Intl.NumberFormat('en-IE', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 0,
    maximumFractionDigits: value % 1 === 0 ? 0 : 2,
  }).format(value)

const formatRatio = (value: number) => `${Math.round(value)}%`

const toNumberOrZero = (value: string) => {
  const parsed = Number.parseFloat(value)
  return Number.isFinite(parsed) ? parsed : 0
}

const clampRatio = (value: string) => {
  const parsed = Number.parseFloat(value)
  if (!Number.isFinite(parsed)) return '0'
  if (parsed < 0) return '0'
  if (parsed > 100) return '100'
  return String(parsed)
}

export default function ExpenseSettingsPage() {
  const [settingsView, { refetch }] = createResource(readyKey, readExpensesSettings)
  const [income, setIncome] = createSignal('0')
  const [ratioFund, setRatioFund] = createSignal('50')
  const [ratioFun, setRatioFun] = createSignal('30')
  const [ratioFuture, setRatioFuture] = createSignal('20')
  const [checklistTitle, setChecklistTitle] = createSignal('')
  const [checklistAmount, setChecklistAmount] = createSignal('')
  const [checklistCategory, setChecklistCategory] = createSignal<CategoryId>('fund')
  const [recurringTitle, setRecurringTitle] = createSignal('')
  const [recurringAmount, setRecurringAmount] = createSignal('')
  const [recurringCategory, setRecurringCategory] = createSignal<CategoryId>('fund')
  const [isSaving, setIsSaving] = createSignal(false)
  const [isAddingChecklist, setIsAddingChecklist] = createSignal(false)
  const [isAddingRecurring, setIsAddingRecurring] = createSignal(false)
  const [isDeletingTemplate, setIsDeletingTemplate] = createSignal(false)
  const [feedbackMessage, setFeedbackMessage] = createSignal<string | undefined>()
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()
  const [deleteTarget, setDeleteTarget] = createSignal<{ id: string, title: string, kind: TemplateKind } | undefined>()

  createEffect(() => {
    const view = settingsView()
    if (!view) return

    setIncome(String(view.settings.income ?? 0))
    setRatioFund(String(view.settings.ratioFund ?? 0))
    setRatioFun(String(view.settings.ratioFun ?? 0))
    setRatioFuture(String(view.settings.ratioFuture ?? 0))
  })

  const ratioSummary = createMemo(() => {
    const fund = toNumberOrZero(ratioFund())
    const fun = toNumberOrZero(ratioFun())
    const future = toNumberOrZero(ratioFuture())
    const total = fund + fun + future

    return {
      fund,
      fun,
      future,
      total,
      totalTone: total === 100 ? 'success' : total > 100 ? 'danger' : 'warning' as 'success' | 'danger' | 'warning',
      totalLabel: total === 100 ? 'Allocation locked' : total > 100 ? 'Over allocated' : 'Unassigned budget',
    }
  })

  const budgetPreview = createMemo(() => {
    const monthlyIncome = toNumberOrZero(income())
    const ratio = ratioSummary()

    return categoryOptions.map((option) => {
      const percent = option.id === 'fund' ? ratio.fund : option.id === 'fun' ? ratio.fun : ratio.future
      const amount = (monthlyIncome * percent) / 100
      const meta = getExpenseCategoryMeta(option.id)

      return {
        id: option.id,
        title: option.label,
        description: option.description,
        compactLabel: meta.compactLabel,
        tone: meta.tone,
        percent,
        amount,
      }
    })
  })

  const saveSettings = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSaving()) return

    const nextIncome = Number.parseFloat(income())
    const fund = Number.parseFloat(ratioFund())
    const fun = Number.parseFloat(ratioFun())
    const future = Number.parseFloat(ratioFuture())

    if (!Number.isFinite(nextIncome) || nextIncome < 0) {
      setErrorMessage('Income needs to be a valid positive number.')
      setFeedbackMessage(undefined)
      return
    }

    if (![fund, fun, future].every((value) => Number.isFinite(value) && value >= 0 && value <= 100)) {
      setErrorMessage('Ratios need to stay between 0 and 100.')
      setFeedbackMessage(undefined)
      return
    }

    setIsSaving(true)
    setErrorMessage(undefined)
    setFeedbackMessage(undefined)

    try {
      await updateExpensesSettings({
        income: nextIncome,
        ratioFund: fund,
        ratioFun: fun,
        ratioFuture: future,
      })

      await refetch()
      setFeedbackMessage('Settings saved. Your budget math now has a plan instead of vibes.')
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save expenses settings')
    } finally {
      setIsSaving(false)
    }
  }

  const addTemplate = async (kind: TemplateKind) => {
    if (kind === 'checklist' && isAddingChecklist()) return
    if (kind === 'recurring' && isAddingRecurring()) return

    const title = kind === 'checklist' ? checklistTitle().trim() : recurringTitle().trim()
    const amountText = kind === 'checklist' ? checklistAmount() : recurringAmount()
    const category = kind === 'checklist' ? checklistCategory() : recurringCategory()
    const amount = Number.parseFloat(amountText)

    if (!title) {
      setErrorMessage('Template title cannot be empty.')
      setFeedbackMessage(undefined)
      return
    }

    if (!Number.isFinite(amount) || amount <= 0) {
      setErrorMessage('Template amount needs to be above zero.')
      setFeedbackMessage(undefined)
      return
    }

    setErrorMessage(undefined)
    setFeedbackMessage(undefined)

    if (kind === 'checklist') {
      setIsAddingChecklist(true)
    } else {
      setIsAddingRecurring(true)
    }

    try {
      if (kind === 'checklist') {
        await upsertChecklistTemplate({ title, amount, category })
        setChecklistTitle('')
        setChecklistAmount('')
        setChecklistCategory('fund')
      } else {
        await upsertRecurringTemplate({ title, amount, category })
        setRecurringTitle('')
        setRecurringAmount('')
        setRecurringCategory('fund')
      }

      await refetch()
      setFeedbackMessage(kind === 'checklist' ? 'Checklist template added.' : 'Recurring template added.')
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save template')
    } finally {
      if (kind === 'checklist') {
        setIsAddingChecklist(false)
      } else {
        setIsAddingRecurring(false)
      }
    }
  }

  const confirmDeleteTemplate = async () => {
    const target = deleteTarget()
    if (!target || isDeletingTemplate()) return

    setIsDeletingTemplate(true)
    setErrorMessage(undefined)
    setFeedbackMessage(undefined)

    try {
      if (target.kind === 'checklist') {
        await deleteChecklistTemplate(target.id)
      } else {
        await deleteRecurringTemplate(target.id)
      }

      setDeleteTarget(undefined)
      await refetch()
      setFeedbackMessage(target.kind === 'checklist' ? 'Checklist template removed.' : 'Recurring template removed.')
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to delete template')
    } finally {
      setIsDeletingTemplate(false)
    }
  }

  const ratioSegments = createMemo(() => {
    const ratio = ratioSummary()
    return [
      { percent: ratio.fund, color: 'danger' as const },
      { percent: ratio.fun, color: 'warning' as const },
      { percent: ratio.future, color: 'success' as const },
    ]
  })

  const checklistTemplates = createMemo(() => settingsView()?.checklist ?? [])
  const recurringTemplates = createMemo(() => settingsView()?.recurring ?? [])

  return (
    <ContentDeck layout="stacked" animate>
      <div class="flex flex-wrap items-center justify-between gap-3">
        <FormBackLink href="/expenses">Back</FormBackLink>
        <Show when={ratioSummary().total !== 100}>
          <Pill tone={ratioSummary().total > 100 ? 'danger' : 'warning'}>
            Total {formatRatio(ratioSummary().total)}
          </Pill>
        </Show>
      </div>

      <Show when={errorMessage()}>
        {(message) => <Notice tone="error" message={message()} />}
      </Show>

      <Show when={feedbackMessage()}>
        {(message) => <Notice tone="info" message={message()} />}
      </Show>

      <Suspense fallback={<ExpenseSettingsLoadingState />}>
        <Show when={settingsView()} fallback={settingsView.error ? <Notice tone="error" message="Expenses settings failed to load." /> : <ExpenseSettingsLoadingState />}>
          <form class="flex flex-col gap-4" onSubmit={saveSettings}>
              <Panel title="Income and split" description={<Pill tone={ratioSummary().totalTone}>{ratioSummary().totalLabel}</Pill>}>
                <div class="flex flex-col gap-5">
                  <AmountHeroField
                    inputId="expense-income"
                    label="Monthly net income"
                    unit="EUR"
                    value={income()}
                    placeholder="0.00"
                    badge={`${formatRatio(ratioSummary().total)} allocated`}
                    badgeTone={ratioSummary().totalTone === 'danger' ? 'danger' : ratioSummary().totalTone === 'warning' ? 'warning' : 'success'}
                    onInput={(event) => setIncome(event.currentTarget.value)}
                  />

                    <div class="space-y-3 border-t border-border/50 pt-4">
                      <div class="flex items-start justify-between gap-3">
                        <div>
                          <div class="text-sm font-semibold text-text-main">Budget split</div>
                        <p class="mt-1 text-xs leading-5 text-text-muted">
                          {ratioSummary().total === 100
                            ? 'Clean split. Nothing extra hanging around.'
                            : ratioSummary().total > 100
                              ? `You are over by ${formatRatio(ratioSummary().total - 100)}.`
                              : `${formatRatio(100 - ratioSummary().total)} is still unassigned.`}
                        </p>
                      </div>
                    </div>

                    <ProgressBar segments={ratioSegments()} />

                    <div class="divide-y divide-border/40 rounded-2xl border border-border/50 bg-panel/50">
                      <For each={budgetPreview()}>
                        {(item) => (
                          <div class="flex items-center justify-between gap-3 px-4 py-3">
                            <div class="min-w-0">
                              <div class="flex items-center gap-2">
                                <span class="text-sm font-semibold text-text-main">{item.title}</span>
                                <Pill tone={item.tone} size="sm">{formatRatio(item.percent)}</Pill>
                              </div>
                              <p class="mt-1 text-xs leading-5 text-text-muted">{item.description}</p>
                            </div>
                            <div class="text-right text-sm font-mono text-text-main">{formatEuro(item.amount)}</div>
                          </div>
                        )}
                      </For>
                    </div>
                  </div>
                </div>
              </Panel>

              <Panel title="Ratio architecture" description="Move the sliders until the month stops looking stupid.">
                <div class="divide-y divide-border/50">
                  <RatioEditorCard
                    title="Fundamentals"
                    caption="The fixed and boring stuff that still matters."
                    tone="danger"
                    value={ratioFund()}
                    amount={budgetPreview()[0]?.amount ?? 0}
                    onInput={setRatioFund}
                  />
                  <RatioEditorCard
                    title="Fun"
                    caption="Quality-of-life spending without the guilt opera."
                    tone="warning"
                    value={ratioFun()}
                    amount={budgetPreview()[1]?.amount ?? 0}
                    onInput={setRatioFun}
                  />
                  <RatioEditorCard
                    title="Future"
                    caption="Savings, buffers, and that smug prepared feeling."
                    tone="success"
                    value={ratioFuture()}
                    amount={budgetPreview()[2]?.amount ?? 0}
                    onInput={setRatioFuture}
                  />
                </div>
              </Panel>

                <Panel title="Monthly checklist" description={<Pill tone="danger">Paid manually each month</Pill>}>
                  <div class="flex flex-col gap-4">
                    <div class="space-y-3 border-b border-border/50 pb-4">
                      <div class="flex items-center justify-between gap-3">
                        <div>
                          <div class="text-sm font-bold tracking-tight text-text-main">Add a checklist item</div>
                          <div class="mt-1 text-xs text-text-muted">Stuff you deliberately check off before closing the sheet.</div>
                        </div>
                        <Pill tone="neutral">{checklistTemplates().length} item{checklistTemplates().length === 1 ? '' : 's'}</Pill>
                      </div>

                      <div class="grid gap-3 sm:grid-cols-[1.4fr_0.8fr]">
                        <label class="flex flex-col gap-2">
                          <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Title</span>
                          <input
                            value={checklistTitle()}
                            onInput={(event) => setChecklistTitle(event.currentTarget.value)}
                            placeholder="Netflix, internet, electricity"
                            class="rounded-2xl border border-border bg-panel px-4 py-3 text-sm text-text-main outline-none transition placeholder:text-text-muted/40 focus:border-accent"
                          />
                        </label>

                        <label class="flex flex-col gap-2">
                          <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Amount</span>
                          <input
                            value={checklistAmount()}
                            onInput={(event) => setChecklistAmount(event.currentTarget.value)}
                            inputmode="decimal"
                            type="number"
                            min="0"
                            step="0.01"
                            placeholder="0.00"
                            class="rounded-2xl border border-border bg-panel px-4 py-3 text-sm text-text-main outline-none transition placeholder:text-text-muted/40 focus:border-accent"
                          />
                        </label>
                      </div>

                      <CategorySegmentedPicker value={checklistCategory()} onChange={setChecklistCategory} />

                      <div class="flex justify-end">
                        <ActionButton busy={isAddingChecklist()} onClick={() => void addTemplate('checklist')}>
                          {isAddingChecklist() ? 'Adding...' : 'Add checklist item'}
                        </ActionButton>
                      </div>
                    </div>

                    <List variant="flush" emptyMessage="No checklist items yet. If your month is that clean, I am suspicious.">
                      <For each={checklistTemplates()}>
                        {(item) => (
                          <ListItem
                            id={item.id}
                            title={item.title}
                            subtitle={
                              <ListMeta>
                                <Pill tone={getExpenseCategoryMeta(item.category).tone} size="sm">{getExpenseCategoryMeta(item.category).compactLabel}</Pill>
                                <ListMetaDivider />
                                <span>Checked off into real expenses later</span>
                              </ListMeta>
                            }
                            value={formatEuro(item.amount)}
                            valueStyle="mono"
                            action={
                              <IconButton
                                ariaLabel={`Delete ${item.title}`}
                                textDanger
                                onClick={() => setDeleteTarget({ id: item.id, title: item.title, kind: 'checklist' })}
                                icon={<TrashIcon />}
                              />
                            }
                          />
                        )}
                      </For>
                    </List>
                  </div>
                </Panel>

                <Panel title="Recurring expenses" description={<Pill tone="success">Injected into each new sheet</Pill>}>
                  <div class="flex flex-col gap-4">
                    <div class="space-y-3 border-b border-border/50 pb-4">
                      <div class="flex items-center justify-between gap-3">
                        <div>
                          <div class="text-sm font-bold tracking-tight text-text-main">Add a recurring draft</div>
                          <div class="mt-1 text-xs text-text-muted">Stable costs that should already exist when a new sheet opens.</div>
                        </div>
                        <Pill tone="neutral">{recurringTemplates().length} item{recurringTemplates().length === 1 ? '' : 's'}</Pill>
                      </div>

                      <div class="grid gap-3 sm:grid-cols-[1.4fr_0.8fr]">
                        <label class="flex flex-col gap-2">
                          <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Title</span>
                          <input
                            value={recurringTitle()}
                            onInput={(event) => setRecurringTitle(event.currentTarget.value)}
                            placeholder="Rent, salary transfer, emergency stash"
                            class="rounded-2xl border border-border bg-panel px-4 py-3 text-sm text-text-main outline-none transition placeholder:text-text-muted/40 focus:border-accent"
                          />
                        </label>

                        <label class="flex flex-col gap-2">
                          <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Amount</span>
                          <input
                            value={recurringAmount()}
                            onInput={(event) => setRecurringAmount(event.currentTarget.value)}
                            inputmode="decimal"
                            type="number"
                            min="0"
                            step="0.01"
                            placeholder="0.00"
                            class="rounded-2xl border border-border bg-panel px-4 py-3 text-sm text-text-main outline-none transition placeholder:text-text-muted/40 focus:border-accent"
                          />
                        </label>
                      </div>

                      <CategorySegmentedPicker value={recurringCategory()} onChange={setRecurringCategory} />

                      <div class="flex justify-end">
                        <ActionButton busy={isAddingRecurring()} onClick={() => void addTemplate('recurring')}>
                          {isAddingRecurring() ? 'Adding...' : 'Add recurring draft'}
                        </ActionButton>
                      </div>
                    </div>

                    <List variant="flush" emptyMessage="No recurring expenses yet. Start with the painful ones; they are the real bosses here.">
                      <For each={recurringTemplates()}>
                        {(item) => (
                          <ListItem
                            id={item.id}
                            title={item.title}
                            subtitle={
                              <ListMeta>
                                <Pill tone={getExpenseCategoryMeta(item.category).tone} size="sm">{getExpenseCategoryMeta(item.category).compactLabel}</Pill>
                                <ListMetaDivider />
                                <span>Preloaded into every new monthly sheet</span>
                              </ListMeta>
                            }
                            value={formatEuro(item.amount)}
                            valueStyle="mono"
                            action={
                              <IconButton
                                ariaLabel={`Delete ${item.title}`}
                                textDanger
                                onClick={() => setDeleteTarget({ id: item.id, title: item.title, kind: 'recurring' })}
                                icon={<TrashIcon />}
                              />
                            }
                          />
                        )}
                      </For>
                    </List>
                  </div>
                </Panel>

              <div class="flex justify-end border-t border-border/50 pt-2">
                <ActionButton type="submit" busy={isSaving()}>
                  {isSaving() ? 'Saving...' : 'Save settings'}
                </ActionButton>
              </div>
          </form>
        </Show>
      </Suspense>

      <ConfirmSheet
        open={Boolean(deleteTarget())}
        eyebrow={deleteTarget()?.kind === 'checklist' ? 'Delete checklist item' : 'Delete recurring draft'}
        title={`Remove ${deleteTarget()?.title ?? 'template'}?`}
        description={deleteTarget()?.kind === 'checklist'
          ? 'This item will stop showing up in your monthly checklist. Already completed history stays untouched.'
          : 'This removes the template from future sheets. Existing recorded expenses are left alone.'}
        confirmLabel="Delete"
        confirmTone="danger"
        busy={isDeletingTemplate()}
        onCancel={() => !isDeletingTemplate() && setDeleteTarget(undefined)}
        onConfirm={confirmDeleteTemplate}
      />
    </ContentDeck>
  )
}

function RatioEditorCard(props: { title: string, caption: string, tone: 'danger' | 'warning' | 'success', value: string, amount: number, onInput: (value: string) => void }) {
  const accentClass = props.tone === 'danger'
    ? 'accent-red-500'
    : props.tone === 'warning'
      ? 'accent-orange-500'
      : 'accent-emerald-500'

  const badgeTone = props.tone === 'danger' ? 'danger' : props.tone === 'warning' ? 'warning' : 'success'

  return (
    <div class="flex flex-col gap-4 py-4 first:pt-0 last:pb-0">
      <div class="flex items-start justify-between gap-3">
        <div>
          <div class="text-base font-bold tracking-tight text-text-main">{props.title}</div>
          <p class="mt-1 text-xs leading-5 text-text-muted">{props.caption}</p>
        </div>
        <Pill tone={badgeTone}>{formatRatio(toNumberOrZero(props.value))}</Pill>
      </div>

      <input
        type="range"
        min="0"
        max="100"
        step="1"
        value={props.value}
        onInput={(event) => props.onInput(clampRatio(event.currentTarget.value))}
        class={`h-2 w-full cursor-pointer appearance-none rounded-full bg-border/60 ${accentClass}`}
      />

      <div class="flex items-end justify-between gap-3">
        <div>
          <div class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Projected amount</div>
          <div class="mt-1 text-xl font-bold tracking-tight text-text-main">{formatEuro(props.amount)}</div>
        </div>

        <label class="flex w-24 flex-col gap-2">
          <span class="text-[0.68rem] font-bold uppercase tracking-[0.24em] text-text-muted">Manual</span>
          <input
            type="number"
            min="0"
            max="100"
            step="1"
            value={props.value}
            onInput={(event) => props.onInput(clampRatio(event.currentTarget.value))}
            class="rounded-xl border border-border bg-background/70 px-3 py-2 text-right text-sm font-semibold text-text-main outline-none transition focus:border-accent"
          />
        </label>
      </div>
    </div>
  )
}

function CategorySegmentedPicker(props: { value: CategoryId, onChange: (value: CategoryId) => void }) {
  return (
    <div class="grid gap-2 sm:grid-cols-3">
      <For each={categoryOptions}>
        {(option) => {
          const meta = getExpenseCategoryMeta(option.id)
          const active = () => props.value === option.id

          return (
            <button
              type="button"
              aria-pressed={active()}
              onClick={() => props.onChange(option.id)}
              class={`rounded-2xl border px-4 py-3 text-left transition-all ${active()
                ? meta.tone === 'danger'
                  ? 'border-danger/30 bg-danger/5 text-text-main'
                  : meta.tone === 'warning'
                    ? 'border-warning/30 bg-warning/5 text-text-main'
                    : 'border-success/30 bg-success/5 text-text-main'
                : 'border-border/70 bg-panel/40 text-text-muted hover:border-accent/30 hover:text-text-main'}`}
            >
              <div class="text-sm font-bold tracking-tight">{option.label}</div>
              <div class="mt-1 text-[0.68rem] uppercase tracking-[0.18em] opacity-75">{meta.compactLabel}</div>
            </button>
          )
        }}
      </For>
    </div>
  )
}

function ExpenseSettingsLoadingState() {
  return (
    <div class="grid gap-4">
      <div class="h-96 animate-pulse rounded-2xl border border-border bg-panel/70" />
      <div class="h-72 animate-pulse rounded-2xl border border-border bg-panel/70" />
      <div class="h-96 animate-pulse rounded-2xl border border-border bg-panel/70" />
      <div class="h-96 animate-pulse rounded-2xl border border-border bg-panel/70" />
    </div>
  )
}

function TrashIcon() {
  return (
    <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.7" class="h-4 w-4">
      <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 5.5h11" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M8 5.5v-1A1.5 1.5 0 0 1 9.5 3h1A1.5 1.5 0 0 1 12 4.5v1" />
      <path stroke-linecap="round" stroke-linejoin="round" d="m6.5 5.5.6 9.2A1.5 1.5 0 0 0 8.59 16h2.82a1.5 1.5 0 0 0 1.49-1.3l.6-9.2" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M8.75 8.5v4.25" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 8.5v4.25" />
    </svg>
  )
}
