import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function CaloriesPage() {
  return (
    <AppShell currentDomain="calories">
      <SectionTitle
        eyebrow="Contract"
        title="Typed calories endpoints"
      />

        <FeaturePlaceholder
          title="Calories domain"
          route="/calories"
          bullets={['`GET /api/calories/target`', '`POST /api/calories/target`', '`POST /api/calories/log`']}
        />
    </AppShell>
  )
}
