import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { AppShell } from '../../components/ui/AppShell'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { SectionTitle } from '../../components/ui/SectionTitle'
import { TextField } from '../../components/ui/TextField'
import { login } from '../../core/auth/store'

export default function AuthPage() {
  const navigate = useNavigate()
  const [email, setEmail] = createSignal('')
  const [password, setPassword] = createSignal('')
  const [isSubmitting, setIsSubmitting] = createSignal(false)
  const [errorMessage, setErrorMessage] = createSignal<string | null>(null)

  const handleSubmit = async (event: SubmitEvent) => {
    event.preventDefault()

    if (isSubmitting()) {
      return
    }

    setIsSubmitting(true)
    setErrorMessage(null)

    try {
      await login({
        email: email().trim(),
        password: password(),
      })

      void navigate('/', { replace: true })
    } catch (error) {
      setErrorMessage(error instanceof Error ? error.message : 'Unable to log in')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <AppShell
      currentDomain="auth"
      eyebrow="Auth"
      title="Sign in"
      description="Short-lived access tokens stay in session storage. Protected routes revalidate on boot and try a refresh before kicking you out."
    >

      <Panel title="Access" description="Use your account credentials to unlock the protected routes.">
        <form class="flex flex-col gap-4" onSubmit={handleSubmit}>
          <TextField
            id="email"
            name="email"
            label="Email"
            type="email"
            value={email()}
            autocomplete="email"
            required
            onInput={(event) => setEmail(event.currentTarget.value)}
          />
          <TextField
            id="password"
            name="password"
            label="Password"
            type="password"
            value={password()}
            autocomplete="current-password"
            required
            onInput={(event) => setPassword(event.currentTarget.value)}
          />

          {errorMessage() ? <Notice tone="error" message={errorMessage()!} /> : null}

          <div class="flex items-center justify-between gap-3">
            <div class="text-sm leading-6 text-text-muted">Routes stay locked until `/api/auth/session` confirms the JWT.</div>
            <ActionButton type="submit" disabled={isSubmitting()} busy={isSubmitting()}>
              {isSubmitting() ? 'Signing in...' : 'Sign in'}
            </ActionButton>
          </div>
        </form>
      </Panel>
    </AppShell>
  )
}
