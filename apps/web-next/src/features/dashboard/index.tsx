import { AppShell } from '../../components/ui/AppShell'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { Panel } from '../../components/ui/Panel'
import { SectionTitle } from '../../components/ui/SectionTitle'

const domains = [
  {
    title: 'Auth first',
    description: 'Session bootstrap, token storage, and route protection belong in core before any domain gets clever.',
    bullets: ['Typed login + session calls are scaffolded', 'UI shell is ready for auth-aware routing'],
  },
  {
    title: 'Domain slices',
    description: 'Calories, expenses, and heat now have isolated entry pages plus domain-local API wrappers.',
    bullets: ['No feature-owned Tailwind', 'Route-level lazy loading stays explicit'],
  },
  {
    title: 'Theme by target',
    description: 'Serverless, VPS, and K8s each get distinct runtime branding from one codebase.',
    bullets: ['`VITE_DEPLOY_TARGET` drives CSS variables', 'PWA shell keeps one build, many moods'],
  },
]

export default function Dashboard() {
  return (
    <AppShell
      currentDomain="home"
      eyebrow="Migration Blueprint"
      title="Solid scaffold, minus the bullshit"
      description="The new frontend now has the right bones: shared core, typed API client, reusable UI, and domain-first routes."
    >
      <SectionTitle
        eyebrow="Control Room"
        title="What is ready"
        description="This is the migration dashboard, not the finished product. It exists to keep the rewrite disciplined."
      />

      <ContentDeck>
        {domains.map((domain) => (
          <Panel title={domain.title} description={domain.description}>
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
        description="Each route gets a dedicated placeholder page so the migration can move one feature at a time without turning into soup."
      />

      <ContentDeck>
        <FeaturePlaceholder
          title="Expenses"
          route="/expenses"
          description="Budget settings and current sheet reads have typed wrappers ready for the first real screen."
          bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`']}
        />
        <FeaturePlaceholder
          title="Calories"
          route="/calories"
          description="Target reads and meal logging can move over once the UI shell is locked in."
          bullets={['`GET /api/calories/target`', '`POST /api/calories/log`']}
        />
        <FeaturePlaceholder
          title="Heat"
          route="/heat"
          description="Refills already have a typed contract, so this one should be the cleanest migration of the lot."
          bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`']}
        />
      </ContentDeck>
    </AppShell>
  )
}
