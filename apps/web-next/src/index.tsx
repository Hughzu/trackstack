/* @refresh reload */
import { render } from 'solid-js/web'
import { Router } from '@solidjs/router'

import { bootstrapAuth } from './core/auth/store'
import './styles/global.css'
import { applyTheme, resolveTheme } from './core/config/theme'
import { routes } from './routes'

const root = document.getElementById('root')

applyTheme(document.documentElement, resolveTheme(import.meta.env.VITE_DEPLOY_TARGET))

void bootstrapAuth().finally(() => {
  render(() => <Router>{routes}</Router>, root!)
})
