import { lazy, type Component } from 'solid-js'

import { ProtectedRoute, PublicOnlyRoute } from './core/auth/guards'

const DashboardPage = lazy(() => import('./features/dashboard'))
const AuthPage = lazy(() => import('./features/auth'))
const CaloriesPage = lazy(() => import('./features/calories'))
const ExpensesPage = lazy(() => import('./features/expenses'))
const NewExpensePage = lazy(() => import('./features/expenses/new'))
const ExpenseSettingsPage = lazy(() => import('./features/expenses/settings'))
const HeatPage = lazy(() => import('./features/heat'))

const withProtectedRoute = (Page: Component): Component => {
  return () => (
    <ProtectedRoute>
      <Page />
    </ProtectedRoute>
  )
}

const withPublicOnlyRoute = (Page: Component): Component => {
  return () => (
    <PublicOnlyRoute>
      <Page />
    </PublicOnlyRoute>
  )
}

export const routes = [
  {
    path: '/',
    component: withProtectedRoute(DashboardPage),
  },
  {
    path: '/login',
    component: withPublicOnlyRoute(AuthPage),
  },
  {
    path: '/calories',
    component: withProtectedRoute(CaloriesPage),
  },
  {
    path: '/expenses',
    component: withProtectedRoute(ExpensesPage),
  },
  {
    path: '/expenses/new',
    component: withProtectedRoute(NewExpensePage),
  },
  {
    path: '/expenses/settings',
    component: withProtectedRoute(ExpenseSettingsPage),
  },
  {
    path: '/heat',
    component: withProtectedRoute(HeatPage),
  },
]
