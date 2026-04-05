import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function ExpensesPage() {
  return (
    <AppShell
      currentDomain="expenses"
      eyebrow="Expenses"
      title="Expenses scaffold is ready"
      description="The route is isolated, settings transport is typed, and the UI shell is ready for the first screen."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed expenses reads"
        description="Settings and current-sheet reads are ready before write flows move over."
      />

        <FeaturePlaceholder
          title="Expenses domain"
          route="/expenses"
          description="Port the dashboard first, then forms and actions once auth state is in place."
          bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`', '`POST /api/expenses/entries`']}
        />
    </AppShell>
  )
}
