import { AppShell } from '../../components/ui/AppShell'
import { FeaturePlaceholder } from '../../components/ui/FeaturePlaceholder'
import { SectionTitle } from '../../components/ui/SectionTitle'

export default function AuthPage() {
  return (
    <AppShell
      currentDomain="auth"
      eyebrow="Auth"
      title="Auth scaffold is in place"
      description="Shared token helpers and typed auth calls are wired so the real flow has a clean place to land."
    >
      <SectionTitle
        eyebrow="Ready"
        title="Typed auth transport"
        description="Login, logout, and session calls already hang off the generated API client."
      />

        <FeaturePlaceholder
          title="Session bootstrap"
          route="/"
          description="Next step: route guards, session bootstrap, and logout handling."
          bullets={['`POST /api/auth/login`', '`GET /api/auth/session`', '`POST /api/auth/logout`']}
        />
    </AppShell>
  )
}
