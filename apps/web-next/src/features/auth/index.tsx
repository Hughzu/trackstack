import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function AuthPage() {
  return (
    <AppShell
      currentDomain="auth"
      eyebrow="Auth"
      title="Bearer auth gets first-class treatment"
      description="The Solid app now has a shared token helper and typed auth transport, which saves us from inventing new weirdness later."
    >
      <SectionTitle
        eyebrow="Ready"
        title="Auth transport scaffold"
        description="Login, logout, and session calls can hang off the generated OpenAPI client instead of hand-rolled fetch soup."
      />

      <FeaturePlaceholder
        title="Session bootstrap"
        route="/"
        description="Use this lane to wire guarded routes and client-side session hydration next."
        bullets={['`POST /api/auth/login`', '`GET /api/auth/session`', '`POST /api/auth/logout`']}
      />
    </AppShell>
  )
}
