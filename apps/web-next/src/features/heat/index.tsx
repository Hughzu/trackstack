import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function HeatPage() {
  return (
    <>
      <SectionTitle eyebrow="Contract" title="Typed heat transport" />

      <FeaturePlaceholder
        title="Heat domain"
        route="/heat"
        bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`', '`DELETE /api/heat/refills/{id}`']}
      />
    </>
  )
}
