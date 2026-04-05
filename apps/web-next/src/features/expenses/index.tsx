import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function ExpensesPage() {
  return (
    <AppShell
      currentDomain="expenses"
      eyebrow="Expenses"
      title="Budget screens get their own clean lane"
      description="Settings transport is typed, the route is isolated, and the UI shell is reusable enough to stop every screen from becoming custom furniture."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed expenses reads"
        description="The migration can start with settings and current-sheet reads before touching any write flows."
      />

      <FeaturePlaceholder
        title="Expenses domain"
        route="/expenses"
        description="Port the dashboard first, then forms, then delete and checklist actions once auth state is in place."
        bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`', '`POST /api/expenses/entries`']}
      />
    </AppShell>
  )
}
