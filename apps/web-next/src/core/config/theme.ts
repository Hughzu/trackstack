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
      background: '#09090b',     // Deep zinc
      surface: '#18181b',
      panel: '#27272a',
      panelMuted: '#3f3f46',
      accent: '#f97316',         // Sharp orange
      accentStrong: '#ea580c',
      accentSoft: 'rgba(249, 115, 22, 0.12)',
      accentInk: '#fff7ed',
      textMain: '#fafafa',
      textMuted: '#a1a1aa',
      border: 'rgba(255, 255, 255, 0.08)',
      success: '#4ade80',
      danger: '#fb7185',
    },
  },
  vps: {
    id: 'vps',
    label: 'VPS',
    palette: {
      background: '#020617',     // Slate void
      surface: '#0f172a',
      panel: '#1e293b',
      panelMuted: '#334155',
      accent: '#2dd4bf',         // Neon teal
      accentStrong: '#14b8a6',
      accentSoft: 'rgba(45, 212, 191, 0.12)',
      accentInk: '#f0fdfa',
      textMain: '#f8fafc',
      textMuted: '#94a3b8',
      border: 'rgba(255, 255, 255, 0.08)',
      success: '#34d399',
      danger: '#f43f5e',
    },
  },
  k8s: {
    id: 'k8s',
    label: 'K8s',
    palette: {
      background: '#0a0a0f',     // Abyssal dim
      surface: '#15131d',
      panel: '#282338',
      panelMuted: '#3d3455',
      accent: '#c084fc',         // Hyper purple
      accentStrong: '#a855f7',
      accentSoft: 'rgba(192, 132, 252, 0.12)',
      accentInk: '#faf5ff',
      textMain: '#fdfbfe',
      textMuted: '#a7a1b5',
      border: 'rgba(255, 255, 255, 0.08)',
      success: '#4ade80',
      danger: '#f43f5e',
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
