import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function HeatPage() {
  return (
    <AppShell
      currentDomain="heat"
      eyebrow="Heat"
      title="Heat scaffold is ready"
      description="The refill contract is already clear, so this domain is a good first feature slice."
    >
      <SectionTitle
        eyebrow="Contract"
        title="Typed heat transport"
        description="Refill reads and writes are ready for real screens once auth is connected."
      />

        <FeaturePlaceholder
          title="Heat domain"
          route="/heat"
          description="Port refill history and seasonal snapshot after the dashboard shell settles."
          bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`', '`DELETE /api/heat/refills/{id}`']}
        />
    </AppShell>
  )
}
