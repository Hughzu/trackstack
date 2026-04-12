/* @refresh reload */
import { onMount } from 'solid-js'
import { render } from 'solid-js/web'
import { Router } from '@solidjs/router'

import { bootstrapAuth } from './core/auth/store'
import { AppRoot } from './components/ui/AppRoot'
import './styles/global.css'
import { applyTheme, resolveTheme } from './core/config/theme'
import { routes } from './routes'

const root = document.getElementById('root')

applyTheme(document.documentElement, resolveTheme(import.meta.env.VITE_DEPLOY_TARGET))

function App() {
  onMount(() => {
    void bootstrapAuth()
  })

  return (
    <Router root={AppRoot} preload={true}>
      {routes}
    </Router>
  )
}

render(() => <App />, root!)
