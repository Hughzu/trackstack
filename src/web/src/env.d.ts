/// <reference path="../.astro/types.d.ts" />

interface ImportMetaEnv {
    readonly DATA_DIR?: string;
    readonly TURSO_EXPENSES_URL?: string;
    readonly TURSO_EXPENSES_TOKEN?: string;
    readonly TURSO_CALORIES_URL?: string;
    readonly TURSO_CALORIES_TOKEN?: string;
    readonly TURSO_HEAT_URL?: string;
    readonly TURSO_HEAT_TOKEN?: string;
    readonly TURSO_NBA_URL?: string;
    readonly TURSO_NBA_TOKEN?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
