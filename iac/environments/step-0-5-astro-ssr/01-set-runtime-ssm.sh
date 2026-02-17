#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/../../bootstrap/config.sh"

SSM_PREFIX="${SSM_PREFIX:-/trackstack/step-0-5}"

if ! command -v aws >/dev/null 2>&1; then
    echo "ERROR: aws CLI is not installed."
    exit 1
fi

echo "=========================================="
echo " STEP 0.5: Set Runtime Secrets (SSM)"
echo "=========================================="
echo ""
echo "Profile: $AWS_PROFILE"
echo "Region: $AWS_REGION"
echo "SSM Prefix: $SSM_PREFIX"
echo ""

prompt_value() {
    local name="$1"
    local label="$2"
    local secret="$3"
    local current_value="${!name}"

    if [ -n "$current_value" ]; then
        echo "$current_value"
        return 0
    fi

    if [ "$secret" = "true" ]; then
        read -r -s -p "$label: " current_value
        echo ""
    else
        read -r -p "$label: " current_value
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

TURSO_USERS_URL=$(prompt_value TURSO_USERS_URL "TURSO_USERS_URL" false)
TURSO_USERS_TOKEN=$(prompt_value TURSO_USERS_TOKEN "TURSO_USERS_TOKEN" true)
TURSO_CALORIES_URL=$(prompt_value TURSO_CALORIES_URL "TURSO_CALORIES_URL" false)
TURSO_CALORIES_TOKEN=$(prompt_value TURSO_CALORIES_TOKEN "TURSO_CALORIES_TOKEN" true)
TURSO_EXPENSES_URL=$(prompt_value TURSO_EXPENSES_URL "TURSO_EXPENSES_URL" false)
TURSO_EXPENSES_TOKEN=$(prompt_value TURSO_EXPENSES_TOKEN "TURSO_EXPENSES_TOKEN" true)
TURSO_HEAT_URL=$(prompt_value TURSO_HEAT_URL "TURSO_HEAT_URL" false)
TURSO_HEAT_TOKEN=$(prompt_value TURSO_HEAT_TOKEN "TURSO_HEAT_TOKEN" true)

put_param "$SSM_PREFIX/runtime/TURSO_USERS_URL" "$TURSO_USERS_URL"
put_param "$SSM_PREFIX/runtime/TURSO_USERS_TOKEN" "$TURSO_USERS_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_CALORIES_URL" "$TURSO_CALORIES_URL"
put_param "$SSM_PREFIX/runtime/TURSO_CALORIES_TOKEN" "$TURSO_CALORIES_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_EXPENSES_URL" "$TURSO_EXPENSES_URL"
put_param "$SSM_PREFIX/runtime/TURSO_EXPENSES_TOKEN" "$TURSO_EXPENSES_TOKEN"
put_param "$SSM_PREFIX/runtime/TURSO_HEAT_URL" "$TURSO_HEAT_URL"
put_param "$SSM_PREFIX/runtime/TURSO_HEAT_TOKEN" "$TURSO_HEAT_TOKEN"

echo ""
echo "=========================================="
echo " Runtime secrets stored in SSM"
echo "=========================================="
echo ""
