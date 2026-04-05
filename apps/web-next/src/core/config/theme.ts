export type DeployTarget = 'serverless' | 'vps' | 'k8s'

type ThemePalette = {
  background: string
  surface: string
  panel: string
  panelMuted: string
  accent: string
  accentStrong: string
  accentSoft: string
  accentInk: string
  textMain: string
  textMuted: string
  border: string
  success: string
  danger: string
}

export type DeployTheme = {
  id: DeployTarget
  label: string
  palette: ThemePalette
}

const deployThemes: Record<DeployTarget, DeployTheme> = {
  serverless: {
    id: 'serverless',
    label: 'Serverless',
    palette: {
      background: '#0d1117',
      surface: '#101826',
      panel: '#162236',
      panelMuted: '#1b2940',
      accent: '#ffb347',
      accentStrong: '#ff7a45',
      accentSoft: 'rgba(255, 179, 71, 0.16)',
      accentInk: '#2c1500',
      textMain: '#f8f4ed',
      textMuted: '#a4afbf',
      border: 'rgba(255, 255, 255, 0.1)',
      success: '#4ade80',
      danger: '#fb7185',
    },
  },
  vps: {
    id: 'vps',
    label: 'VPS',
    palette: {
      background: '#091217',
      surface: '#10232c',
      panel: '#12313c',
      panelMuted: '#18414d',
      accent: '#5eead4',
      accentStrong: '#22d3ee',
      accentSoft: 'rgba(94, 234, 212, 0.18)',
      accentInk: '#032b2f',
      textMain: '#eefcf9',
      textMuted: '#96bdb9',
      border: 'rgba(124, 242, 226, 0.14)',
      success: '#86efac',
      danger: '#fda4af',
    },
  },
  k8s: {
    id: 'k8s',
    label: 'K8s',
    palette: {
      background: '#120d16',
      surface: '#22152a',
      panel: '#2c1b37',
      panelMuted: '#3a2249',
      accent: '#f59e0b',
      accentStrong: '#ef4444',
      accentSoft: 'rgba(245, 158, 11, 0.18)',
      accentInk: '#331400',
      textMain: '#fff4ea',
      textMuted: '#c8ab9a',
      border: 'rgba(255, 206, 160, 0.14)',
      success: '#86efac',
      danger: '#fda4af',
    },
  },
}

const deployAliases: Record<string, DeployTarget> = {
  ecs: 'k8s',
  eks: 'k8s',
  kubernetes: 'k8s',
  lambda: 'serverless',
}

export const resolveTheme = (input?: string): DeployTheme => {
  const normalized = (input ?? 'serverless').trim().toLowerCase()
  const target = deployAliases[normalized] ?? normalized

  if (target in deployThemes) {
    return deployThemes[target as DeployTarget]
  }

  return deployThemes.serverless
}

export const applyTheme = (root: HTMLElement, theme: DeployTheme) => {
  root.dataset.deployTarget = theme.id
  root.style.setProperty('--color-background', theme.palette.background)
  root.style.setProperty('--color-surface', theme.palette.surface)
  root.style.setProperty('--color-panel', theme.palette.panel)
  root.style.setProperty('--color-panel-muted', theme.palette.panelMuted)
  root.style.setProperty('--color-accent', theme.palette.accent)
  root.style.setProperty('--color-accent-strong', theme.palette.accentStrong)
  root.style.setProperty('--color-accent-soft', theme.palette.accentSoft)
  root.style.setProperty('--color-accent-ink', theme.palette.accentInk)
  root.style.setProperty('--color-text-main', theme.palette.textMain)
  root.style.setProperty('--color-text-muted', theme.palette.textMuted)
  root.style.setProperty('--color-border', theme.palette.border)
  root.style.setProperty('--color-success', theme.palette.success)
  root.style.setProperty('--color-danger', theme.palette.danger)
}
