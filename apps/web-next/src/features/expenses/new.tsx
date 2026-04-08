import type { JSX } from 'solid-js'
import { createMemo, createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ChoiceCards, type ChoiceCardOption } from '../../components/ui/ChoiceCards'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormFieldRow, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { TextField } from '../../components/ui/TextField'
import { createExpenseEntry } from './api/client'
import { getExpenseCategoryMeta } from './display'

type CategoryId = 'fund' | 'fun' | 'future'

const categoryOptions: Array<{ id: CategoryId, title: string, short: string, icon: JSX.Element }> = [
  {
    id: 'fund',
    title: 'Fundamentals',
    short: 'Bills and essentials',
    icon: <ExpenseCategoryIcon kind="fund" />,
  },
  {
    id: 'fun',
    title: 'Fun',
    short: 'Pleasure and chaos',
    icon: <ExpenseCategoryIcon kind="fun" />,
  },
  {
    id: 'future',
    title: 'Future',
    short: 'Long-game money',
    icon: <ExpenseCategoryIcon kind="future" />,
  },
]

const today = new Date().toISOString().slice(0, 10)

function ExpenseCategoryIcon(props: { kind: CategoryId }) {
  if (props.kind === 'fund') {
    return (
      <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.7">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3.5 8.5 10 3l6.5 5.5V16a1 1 0 0 1-1 1h-3.5v-4.5h-4V17H4.5a1 1 0 0 1-1-1Z" />
      </svg>
    )
  }

  if (props.kind === 'fun') {
    return (
      <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.7">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 3.5v3m8-3v3M4 8.5h12M5.5 6.5h9a1.5 1.5 0 0 1 1.5 1.5v6A2.5 2.5 0 0 1 13.5 16.5h-7A2.5 2.5 0 0 1 4 14V8a1.5 1.5 0 0 1 1.5-1.5Z" />
      </svg>
    )
  }

  return (
    <svg viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.7">
      <path stroke-linecap="round" stroke-linejoin="round" d="M10 16.5c3.314 0 6-2.686 6-6 0-1.657-1.343-3-3-3-.63 0-1.213.194-1.694.526C10.974 5.986 9.615 4.5 8 4.5c-2.21 0-4 1.79-4 4 0 4.418 6 8 6 8Z" />
    </svg>
  )
}

export default function NewExpensePage() {
  const navigate = useNavigate()
  const [amount, setAmount] = createSignal('')
  const [title, setTitle] = createSignal('')
  const [date, setDate] = createSignal(today)
  const [category, setCategory] = createSignal<CategoryId>('fund')
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  const selectedCategory = createMemo(() => getExpenseCategoryMeta(category()))
  const categoryBadgeTone = createMemo(() => selectedCategory().tone)
  const choices = createMemo<ChoiceCardOption[]>(() => categoryOptions.map((option) => ({
    value: option.id,
    title: option.title,
    description: option.short,
    icon: option.icon,
    tone: getExpenseCategoryMeta(option.id).tone,
  })))

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSubmitting()) return

    const parsedAmount = Number.parseFloat(amount())
    if (!Number.isFinite(parsedAmount) || parsedAmount <= 0) {
      setErrorMessage('Amount must be greater than zero')
      return
    }

    setIsSubmitting(true)
    setErrorMessage(undefined)

    try {
      await createExpenseEntry({
        title: title().trim(),
        amount: parsedAmount,
        category: category(),
        date: date(),
      })

      void navigate('/expenses', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to save expense')
      setIsSubmitting(false)
    }
  }

  return (
    <ContentDeck layout="stacked" animate>
      <DataRow variant="header">
        <FormBackLink href="/expenses">Back</FormBackLink>
      </DataRow>

      <FormStack onSubmit={handleSubmit}>
        <FormSection>
          <AmountHeroField
            inputId="expense-amount"
            label="Amount"
            unit="EUR"
            value={amount()}
            placeholder="0.00"
            badge={selectedCategory().compactLabel}
            badgeTone={categoryBadgeTone()}
            onInput={(event) => setAmount(event.currentTarget.value)}
          />

          <ChoiceCards testId="expense-category-choice" options={choices()} value={category()} onChange={(value) => setCategory(value as CategoryId)} />

          <FormFieldRow>
            <TextField
              id="expense-title"
              name="title"
              label="Title"
              value={title()}
              onInput={(event) => setTitle(event.currentTarget.value)}
            />

            <TextField
              id="expense-date"
              name="date"
              label="Date"
              type="date"
              value={date()}
              onInput={(event) => setDate(event.currentTarget.value)}
            />
          </FormFieldRow>

          {errorMessage() ? <Notice tone="error" message={errorMessage()!} /> : null}
        </FormSection>

        <FormActions>
          <ActionButton type="submit" disabled={isSubmitting()} busy={isSubmitting()} ariaLabel="Save Expense">
            {isSubmitting() ? 'Saving...' : 'Save Expense'}
          </ActionButton>
        </FormActions>
      </FormStack>
    </ContentDeck>
  )
}
