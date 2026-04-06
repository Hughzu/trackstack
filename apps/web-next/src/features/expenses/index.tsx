import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function ExpensesPage() {
  return (
    <>
      <SectionTitle eyebrow="Contract" title="Typed expenses reads" />

      <FeaturePlaceholder
        title="Expenses domain"
        route="/expenses"
        bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`', '`POST /api/expenses/entries`']}
      />
    </>
  )
}
