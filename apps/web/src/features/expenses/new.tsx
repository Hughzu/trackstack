import { createMemo, createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AmountHeroField } from '../../components/ui/AmountHeroField'
import { ChoiceCards } from '../../components/ui/ChoiceCards'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { DataRow } from '../../components/ui/DataRow'
import { FormActions, FormBackLink, FormFieldRow, FormSection, FormStack } from '../../components/ui/Form'
import { Notice } from '../../components/ui/Notice'
import { TextField } from '../../components/ui/TextField'
import { createExpenseEntry } from './api/client'
import { createExpenseCategoryChoices, getExpenseCategoryMeta, type ExpenseCategoryId } from './display'

const today = new Date().toISOString().slice(0, 10)

export default function NewExpensePage() {
  const navigate = useNavigate()
  const [amount, setAmount] = createSignal('')
  const [title, setTitle] = createSignal('')
  const [date, setDate] = createSignal(today)
  const [category, setCategory] = createSignal<ExpenseCategoryId>('fund')
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | undefined>()

  const selectedCategory = createMemo(() => getExpenseCategoryMeta(category()))
  const categoryBadgeTone = createMemo(() => selectedCategory().tone)
  const choices = createMemo(() => createExpenseCategoryChoices().map((option) => ({
    value: option.value,
    title: option.title,
    description: option.description,
    tone: option.tone,
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
    } finally {
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

          <ChoiceCards testId="expense-category-choice" options={choices()} value={category()} onChange={(value) => setCategory(value as ExpenseCategoryId)} />

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
