import { AppShell } from '../../components/ui/AppShell'
import { ContentDeck } from '../../components/ui/ContentDeck'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { Panel } from '../../components/ui/Panel'
import { SectionTitle } from '../../components/ui/SectionTitle'

const domains = [
  {
    title: 'Auth first',
    description: 'Session bootstrap and route protection belong in core before feature work starts.',
    bullets: ['Typed login + session calls are scaffolded', 'UI shell is ready for auth-aware routing'],
  },
  {
    title: 'Domain slices',
    description: 'Calories, expenses, and heat each have their own route and thin API wrapper.',
    bullets: ['No feature-owned Tailwind', 'Route-level lazy loading stays explicit'],
  },
  {
    title: 'Theme by target',
    description: 'Serverless, VPS, and K8s each get runtime branding from one codebase.',
    bullets: ['`VITE_DEPLOY_TARGET` drives CSS variables', 'PWA shell keeps one build, many moods'],
  },
]

export default function Dashboard() {
  return (
    <AppShell
      currentDomain="home"
      eyebrow="Migration Blueprint"
      title="Solid scaffold in place"
      description="The bootstrap covers shared core utilities, typed API access, reusable UI, and domain routes."
    >
      <SectionTitle
        eyebrow="Control Room"
        title="What is ready"
        description="This page tracks the scaffold, not feature parity."
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
        description="Each route has a placeholder so feature work can land one slice at a time."
      />

      <ContentDeck>
        <FeaturePlaceholder
          title="Expenses"
          route="/expenses"
          description="Settings and current-sheet reads are ready for the first real screen."
          bullets={['`GET /api/expenses/settings`', '`GET /api/expenses/sheet/current`']}
        />
        <FeaturePlaceholder
          title="Calories"
          route="/calories"
          description="Target reads and meal logging can move once the shell is stable."
          bullets={['`GET /api/calories/target`', '`POST /api/calories/log`']}
        />
        <FeaturePlaceholder
          title="Heat"
          route="/heat"
          description="Refills already have a typed contract, so this should be the cleanest first slice."
          bullets={['`GET /api/heat/refills`', '`POST /api/heat/refills`']}
        />
      </ContentDeck>
    </AppShell>
  )
}
