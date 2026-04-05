import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function CaloriesPage() {
  return (
    <AppShell
      currentDomain="calories"
      eyebrow="Calories"
      title="Nutrition moves over when the shell stops wobbling"
      description="The domain route exists, the API wrapper exists, and the UI layer is ready for a real dashboard instead of another Astro copy-paste job."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed calories entry points"
        description="Targets and logs already have schema coverage, which means this feature can migrate without guesswork."
      />

      <FeaturePlaceholder
        title="Calories domain"
        route="/calories"
        description="Start with the target + dashboard shell, then port logging flows after auth bootstrap lands."
        bullets={['`GET /api/calories/target`', '`POST /api/calories/target`', '`POST /api/calories/log`']}
      />
    </AppShell>
  )
}
