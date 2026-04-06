import { AppShell } from '../../components/ui/AppShell'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { Panel } from '../../components/ui/Panel'
import { SectionTitle } from '../../components/ui/SectionTitle'

const domains = [
  {
    title: 'Auth first',
    bullets: ['Typed login + session calls are scaffolded', 'UI shell is ready for auth-aware routing'],
  },
  {
    title: 'Domain slices',
    bullets: ['No feature-owned Tailwind', 'Route-level lazy loading stays explicit'],
  },
  {
    title: 'Theme by target',
    bullets: ['`VITE_DEPLOY_TARGET` drives CSS variables', 'PWA shell keeps one build, many moods'],
  },
]

export default function Dashboard() {
  return (
    <AppShell currentDomain="home">
      <SectionTitle
        eyebrow="Control Room"
        title="What is ready"
      />

      <ContentDeck>
        {domains.map((domain) => (
          <Panel title={domain.title}>
            <ul>
              {domain.bullets.map((bullet) => (
                <li>{bullet}</li>
              ))}
            </ul>
          </Panel>
        ))}
      </ContentDeck>

      <SectionTitle
        eyebrow="Domains"
        title="Rewrite lanes"
      />

      <ContentDeck>
        <FeaturePlaceholder
          title="Expenses"
          route="/expenses"
          bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`']}
        />
        <FeaturePlaceholder
          title="Calories"
          route="/calories"
          bullets={['`GET /api/calories/target`', '`POST /api/calories/log`']}
        />
        <FeaturePlaceholder
          title="Heat"
          route="/heat"
          bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`']}
        />
      </ContentDeck>
    </AppShell>
  )
}
