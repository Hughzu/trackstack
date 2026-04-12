import { createEffect, createMemo, createResource, createSignal, Show, Suspense } from 'solid-js'

import { ActionButton } from '../../components/ui/ActionButton'
import { ConfirmSheet } from '../../components/ui/ConfirmSheet'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { FormActions, FormBackLink, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { Pill } from '../../components/ui/Pill'
import { authState } from '../../core/auth/state'
import {
  deleteChecklistTemplate,
  deleteRecurringTemplate,
  readExpensesSettings,
  updateExpensesSettings,
  upsertChecklistTemplate,
  upsertRecurringTemplate,
} from './api/client'
import { SettingsIncomeCard } from './components/settings-income-card'
import { SettingsRatioCard } from './components/settings-ratio-card'
import { ExpenseSettingsLoadingState } from './components/settings-skeleton'
import { TemplateCard } from './components/template-card'
import {
  clampRatio,
  createExpenseRatioSummary,
  createExpenseSettingsBudgetPreview,
  formatRatio,
  toNumberOrZero,
  type ExpenseCategoryId,
} from './display'

type TemplateKind = 'checklist' | 'recurring'

const readyKey = () => (authState().status === 'authenticated' ? 'ready' : undefined)

const formatEuro = (value: number) =>
  new Intl.NumberFormat('en-IE', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 0,
    maximumFractionDigits: value % 1 === 0 ? 0 : 2,
  }).format(value)

export default function ExpenseSettingsPage() {
  const [settingsView, { refetch }] = createResource(readyKey, readExpensesSettings)
  const [income, setIncome] = createSignal('0')
  const [ratioFund, setRatioFund] = createSignal('50')
  const [ratioFun, setRatioFun] = createSignal('30')
  const [ratioFuture, setRatioFuture] = createSignal('20')
  const [checklistTitle, setChecklistTitle] = createSignal('')
  const [checklistAmount, setChecklistAmount] = createSignal('')
  const [checklistCategory, setChecklistCategory] = createSignal<ExpenseCategoryId>('fund')
  const [recurringTitle, setRecurringTitle] = createSignal('')
  const [recurringAmount, setRecurringAmount] = createSignal('')
  const [recurringCategory, setRecurringCategory] = createSignal<ExpenseCategoryId>('fund')
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

  const ratioSummary = createMemo(() => createExpenseRatioSummary({
    fund: toNumberOrZero(ratioFund()),
    fun: toNumberOrZero(ratioFun()),
    future: toNumberOrZero(ratioFuture()),
  }))

  const budgetPreview = createMemo(() => createExpenseSettingsBudgetPreview(toNumberOrZero(income()), ratioSummary(), formatEuro))

  const checklistTemplates = createMemo(() => settingsView()?.checklist ?? [])
  const recurringTemplates = createMemo(() => settingsView()?.recurring ?? [])

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
          <FormStack onSubmit={saveSettings}>
            <SettingsIncomeCard
              income={income()}
              ratioSummary={ratioSummary()}
              budgetPreview={budgetPreview()}
              onIncomeInput={(event) => setIncome(event.currentTarget.value)}
            />

            <SettingsRatioCard
              items={[
                {
                  id: 'fund',
                  title: 'Fundamentals',
                  caption: 'The fixed and boring stuff that still matters.',
                  tone: 'danger',
                  value: ratioFund(),
                  amount: budgetPreview()[0]?.value ?? formatEuro(0),
                  onInput: (value) => setRatioFund(clampRatio(value)),
                },
                {
                  id: 'fun',
                  title: 'Fun',
                  caption: 'Quality-of-life spending without the guilt opera.',
                  tone: 'warning',
                  value: ratioFun(),
                  amount: budgetPreview()[1]?.value ?? formatEuro(0),
                  onInput: (value) => setRatioFun(clampRatio(value)),
                },
                {
                  id: 'future',
                  title: 'Future',
                  caption: 'Savings, buffers, and that smug prepared feeling.',
                  tone: 'success',
                  value: ratioFuture(),
                  amount: budgetPreview()[2]?.value ?? formatEuro(0),
                  onInput: (value) => setRatioFuture(clampRatio(value)),
                },
              ]}
            />

            <TemplateCard
              kind="checklist"
              title="Monthly checklist"
              description="Add a checklist item"
              badgeTone="danger"
              badgeLabel="Paid manually each month"
              inputTitle={checklistTitle()}
              inputAmount={checklistAmount()}
              inputCategory={checklistCategory()}
              templates={checklistTemplates()}
              busy={isAddingChecklist()}
              submitLabel="Add checklist item"
              emptyMessage="No checklist items yet. If your month is that clean, I am suspicious."
              helperCopy="Stuff you deliberately check off before closing the sheet."
              titlePlaceholder="Netflix, internet, electricity"
              onTitleInput={(event) => setChecklistTitle(event.currentTarget.value)}
              onAmountInput={(event) => setChecklistAmount(event.currentTarget.value)}
              onCategoryChange={setChecklistCategory}
              onSubmit={() => void addTemplate('checklist')}
              onDelete={(item) => setDeleteTarget({ id: item.id, title: item.title, kind: 'checklist' })}
            />

            <TemplateCard
              kind="recurring"
              title="Recurring expenses"
              description="Add a recurring draft"
              badgeTone="success"
              badgeLabel="Injected into each new sheet"
              inputTitle={recurringTitle()}
              inputAmount={recurringAmount()}
              inputCategory={recurringCategory()}
              templates={recurringTemplates()}
              busy={isAddingRecurring()}
              submitLabel="Add recurring draft"
              emptyMessage="No recurring expenses yet. Start with the painful ones; they are the real bosses here."
              helperCopy="Stable costs that should already exist when a new sheet opens."
              titlePlaceholder="Rent, salary transfer, emergency stash"
              onTitleInput={(event) => setRecurringTitle(event.currentTarget.value)}
              onAmountInput={(event) => setRecurringAmount(event.currentTarget.value)}
              onCategoryChange={setRecurringCategory}
              onSubmit={() => void addTemplate('recurring')}
              onDelete={(item) => setDeleteTarget({ id: item.id, title: item.title, kind: 'recurring' })}
            />

            <FormActions>
              <ActionButton type="submit" busy={isSaving()}>
                {isSaving() ? 'Saving...' : 'Save settings'}
              </ActionButton>
            </FormActions>
          </FormStack>
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
