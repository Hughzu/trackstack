interface ImportMetaEnv {
  readonly VITE_DEPLOY_TARGET?: string
  readonly VITE_API_BASE_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
