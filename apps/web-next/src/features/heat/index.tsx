import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function HeatPage() {
  return (
    <AppShell
      currentDomain="heat"
      eyebrow="Heat"
      title="Pellet tracking should be the easy win"
      description="The refill contract is already explicit and clean, so this domain can prove the Solid migration without ceremonial pain."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed heat transport"
        description="Refill reads and writes are ready to wrap real screens once the auth flow is connected."
      />

      <FeaturePlaceholder
        title="Heat domain"
        route="/heat"
        description="Port the refill history and seasonal snapshot after the dashboard shell settles."
        bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`', '`DELETE /api/heat/refills/{id}`']}
      />
    </AppShell>
  )
}
