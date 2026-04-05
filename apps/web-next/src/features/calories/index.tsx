import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function CaloriesPage() {
  return (
    <AppShell
      currentDomain="calories"
      eyebrow="Calories"
      title="Calories scaffold is ready"
      description="The route, wrapper, and UI shell are in place for the first real calories screen."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed calories endpoints"
        description="Targets and logs already have schema coverage."
      />

        <FeaturePlaceholder
          title="Calories domain"
          route="/calories"
          description="Start with the target view, then move logging once auth bootstrap lands."
          bullets={['`GET /api/calories/target`', '`POST /api/calories/target`', '`POST /api/calories/log`']}
        />
    </AppShell>
  )
}
