import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'

import { ActionButton } from '../../components/ui/ActionButton'
import { Notice } from '../../components/ui/Notice'
import { Panel } from '../../components/ui/Panel'
import { FormStack, FormActions } from '../../components/ui/Form'
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
    <Panel title="Access">
      <FormStack onSubmit={handleSubmit}>
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

        <FormActions>
          <ActionButton type="submit" disabled={isSubmitting()} busy={isSubmitting()}>
            {isSubmitting() ? 'Signing in...' : 'Sign in'}
          </ActionButton>
        </FormActions>
      </FormStack>
    </Panel>
  )
}
