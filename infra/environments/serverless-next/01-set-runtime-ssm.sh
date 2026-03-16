#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../bootstrap/bootstrap/config.sh"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
SERVER_ENV_FILE="$REPO_ROOT/apps/server/.env"

SSM_PREFIX="${SSM_PREFIX:-/trackstack/serverless-next}"

if ! command -v aws >/dev/null 2>&1; then
    echo "ERROR: aws CLI is not installed."
    exit 1
fi

echo "=========================================="
echo " SERVERLESS-NEXT: Set Runtime Secrets (SSM)"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Region: $AWS_REGION"
echo "SSM Prefix: $SSM_PREFIX"
echo ""

load_env_defaults() {
    local env_file="$1"

    if [ ! -f "$env_file" ]; then
        return 0
    fi

    while IFS='=' read -r key value; do
        case "$key" in
            ''|'#'*)
                continue
                ;;
            TURSO_USERS_URL_HTTP|TURSO_USERS_TOKEN|TURSO_CALORIES_URL_HTTP|TURSO_CALORIES_TOKEN|TURSO_EXPENSES_URL_HTTP|TURSO_EXPENSES_TOKEN|TURSO_HEAT_URL_HTTP|TURSO_HEAT_TOKEN)
                if [ -z "${!key}" ]; then
                    export "$key=$value"
                fi
                ;;
        esac
    done < "$env_file"
}

load_env_defaults "$SERVER_ENV_FILE"

prompt_value() {
    local name="$1"
    local label="$2"
    local secret="$3"
    local fallback="$4"
    local current_value="${!name}"

    if [ -n "$current_value" ]; then
        echo "$current_value"
        return 0
    fi

    if [ -n "$fallback" ]; then
        label="$label [$fallback]"
    fi

    if [ "$secret" = "true" ]; then
        read -r -s -p "$label: " current_value
        echo ""
    else
        read -r -p "$label: " current_value
    fi

    if [ -z "$current_value" ]; then
        current_value="$fallback"
    fi

    echo "$current_value"
}

put_param() {
    local name="$1"
    local value="$2"

    aws ssm put-parameter \
        --name "$name" \
        --type SecureString \
        --value "$value" \
        --overwrite \
        --region "$AWS_REGION" \
        --profile "$AWS_PROFILE" >/dev/null
}

APP_ENV=$(prompt_value APP_ENV "APP_ENV" false "serverless-next")
LOG_LEVEL=$(prompt_value LOG_LEVEL "LOG_LEVEL" false "info")
DB_CONNECTION_MODE=$(prompt_value DB_CONNECTION_MODE "DB_CONNECTION_MODE" false "HTTP")
CORS_ALLOWED_ORIGIN=$(prompt_value CORS_ALLOWED_ORIGIN "CORS_ALLOWED_ORIGIN" false "")
AUTH_COOKIE_NAME=$(prompt_value AUTH_COOKIE_NAME "AUTH_COOKIE_NAME" false "trackstack_session")
AUTH_COOKIE_SECURE=$(prompt_value AUTH_COOKIE_SECURE "AUTH_COOKIE_SECURE" false "true")
AUTH_COOKIE_SAMESITE=$(prompt_value AUTH_COOKIE_SAMESITE "AUTH_COOKIE_SAMESITE" false "lax")
AUTH_SESSION_IDLE_SECONDS=$(prompt_value AUTH_SESSION_IDLE_SECONDS "AUTH_SESSION_IDLE_SECONDS" false "604800")
AUTH_SESSION_ABSOLUTE_SECONDS=$(prompt_value AUTH_SESSION_ABSOLUTE_SECONDS "AUTH_SESSION_ABSOLUTE_SECONDS" false "7776000")
AUTH_SESSION_ROTATE_AFTER_SECONDS=$(prompt_value AUTH_SESSION_ROTATE_AFTER_SECONDS "AUTH_SESSION_ROTATE_AFTER_SECONDS" false "1800")
AUTH_SESSION_ROTATION_GRACE_SECONDS=$(prompt_value AUTH_SESSION_ROTATION_GRACE_SECONDS "AUTH_SESSION_ROTATION_GRACE_SECONDS" false "300")
AUTH_SESSION_TOUCH_SECONDS=$(prompt_value AUTH_SESSION_TOUCH_SECONDS "AUTH_SESSION_TOUCH_SECONDS" false "300")
TURSO_USERS_URL_HTTP=$(prompt_value TURSO_USERS_URL_HTTP "TURSO_USERS_URL_HTTP" false "")
TURSO_USERS_TOKEN=$(prompt_value TURSO_USERS_TOKEN "TURSO_USERS_TOKEN" true "")
TURSO_CALORIES_URL_HTTP=$(prompt_value TURSO_CALORIES_URL_HTTP "TURSO_CALORIES_URL_HTTP" false "")
TURSO_CALORIES_TOKEN=$(prompt_value TURSO_CALORIES_TOKEN "TURSO_CALORIES_TOKEN" true "")
TURSO_EXPENSES_URL_HTTP=$(prompt_value TURSO_EXPENSES_URL_HTTP "TURSO_EXPENSES_URL_HTTP" false "")
TURSO_EXPENSES_TOKEN=$(prompt_value TURSO_EXPENSES_TOKEN "TURSO_EXPENSES_TOKEN" true "")
TURSO_HEAT_URL_HTTP=$(prompt_value TURSO_HEAT_URL_HTTP "TURSO_HEAT_URL_HTTP" false "")
TURSO_HEAT_TOKEN=$(prompt_value TURSO_HEAT_TOKEN "TURSO_HEAT_TOKEN" true "")

put_param "$SSM_PREFIX/runtime/APP_ENV" "$APP_ENV"
put_param "$SSM_PREFIX/runtime/LOG_LEVEL" "$LOG_LEVEL"
put_param "$SSM_PREFIX/runtime/DB_CONNECTION_MODE" "$DB_CONNECTION_MODE"
put_param "$SSM_PREFIX/runtime/CORS_ALLOWED_ORIGIN" "$CORS_ALLOWED_ORIGIN"
put_param "$SSM_PREFIX/runtime/AUTH_COOKIE_NAME" "$AUTH_COOKIE_NAME"
put_param "$SSM_PREFIX/runtime/AUTH_COOKIE_SECURE" "$AUTH_COOKIE_SECURE"
put_param "$SSM_PREFIX/runtime/AUTH_COOKIE_SAMESITE" "$AUTH_COOKIE_SAMESITE"
put_param "$SSM_PREFIX/runtime/AUTH_SESSION_IDLE_SECONDS" "$AUTH_SESSION_IDLE_SECONDS"
put_param "$SSM_PREFIX/runtime/AUTH_SESSION_ABSOLUTE_SECONDS" "$AUTH_SESSION_ABSOLUTE_SECONDS"
put_param "$SSM_PREFIX/runtime/AUTH_SESSION_ROTATE_AFTER_SECONDS" "$AUTH_SESSION_ROTATE_AFTER_SECONDS"
put_param "$SSM_PREFIX/runtime/AUTH_SESSION_ROTATION_GRACE_SECONDS" "$AUTH_SESSION_ROTATION_GRACE_SECONDS"
put_param "$SSM_PREFIX/runtime/AUTH_SESSION_TOUCH_SECONDS" "$AUTH_SESSION_TOUCH_SECONDS"
put_param "$SSM_PREFIX/runtime/TURSO_USERS_URL_HTTP" "$TURSO_USERS_URL_HTTP"
put_param "$SSM_PREFIX/runtime/TURSO_USERS_TOKEN" "$TURSO_USERS_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_CALORIES_URL_HTTP" "$TURSO_CALORIES_URL_HTTP"
put_param "$SSM_PREFIX/runtime/TURSO_CALORIES_TOKEN" "$TURSO_CALORIES_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_EXPENSES_URL_HTTP" "$TURSO_EXPENSES_URL_HTTP"
put_param "$SSM_PREFIX/runtime/TURSO_EXPENSES_TOKEN" "$TURSO_EXPENSES_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_HEAT_URL_HTTP" "$TURSO_HEAT_URL_HTTP"
put_param "$SSM_PREFIX/runtime/TURSO_HEAT_TOKEN" "$TURSO_HEAT_TOKEN"

echo ""
echo "=========================================="
echo " Runtime secrets stored in SSM"
echo "=========================================="
echo ""
