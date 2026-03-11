/// <reference path="../.astro/types.d.ts" />

interface ImportMetaEnv {
    readonly PUBLIC_API_BASE_URL?: string;
    readonly DATA_DIR?: string;
    readonly TURSO_EXPENSES_URL?: string;
    readonly TURSO_EXPENSES_TOKEN?: string;
    readonly TURSO_CALORIES_URL?: string;
    readonly TURSO_CALORIES_TOKEN?: string;
    readonly TURSO_HEAT_URL?: string;
    readonly TURSO_HEAT_TOKEN?: string;
    readonly TURSO_USERS_URL?: string;
    readonly TURSO_USERS_TOKEN?: string;
    readonly E2E_TEST_EMAIL?: string;
    readonly E2E_TEST_PASSWORD?: string;
    readonly TURSO_NBA_URL?: string;
    readonly TURSO_NBA_TOKEN?: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}

declare namespace App {
    interface Locals {
        userId?: string;
        sessionId?: string;
    }
}
