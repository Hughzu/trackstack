import { lazy } from 'solid-js'

export const routes = [
  {
    path: '/',
    component: lazy(() => import('./features/dashboard')),
  },
  {
    path: '/login',
    component: lazy(() => import('./features/auth')),
  },
  {
    path: '/calories',
    component: lazy(() => import('./features/calories')),
  },
  {
    path: '/expenses',
    component: lazy(() => import('./features/expenses')),
  },
  {
    path: '/heat',
    component: lazy(() => import('./features/heat')),
  },
]
