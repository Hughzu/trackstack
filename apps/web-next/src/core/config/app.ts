export type AppDomain = {
  id: string
  label: string
  href: string
}

export type AppConfig = {
  appName: string
  domains: AppDomain[]
}

export const appConfig: AppConfig = {
  appName: 'TrackStack',
  domains: [
    { id: 'home', label: 'Overview', href: '/' },
    { id: 'expenses', label: 'Expenses', href: '/expenses' },
    { id: 'calories', label: 'Calories', href: '/calories' },
    { id: 'heat', label: 'Heat', href: '/heat' },
    { id: 'auth', label: 'Auth', href: '/login' },
  ],
}
