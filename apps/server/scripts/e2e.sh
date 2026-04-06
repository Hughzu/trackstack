#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

if [[ -f "$ENV_FILE" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

BASE_URL="${BASE_URL:-http://localhost:8080}"
EMAIL="${E2E_TEST_EMAIL:-}"
PASSWORD="${E2E_TEST_PASSWORD:-}"

if [[ -z "$EMAIL" || -z "$PASSWORD" ]]; then
  printf 'E2E_TEST_EMAIL and E2E_TEST_PASSWORD must be set.\n' >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  printf 'python3 is required for JSON assertions.\n' >&2
  exit 1
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
COOKIE_JAR="$TMP_DIR/cookies.txt"

TODAY_UTC="$(date -u +%F)"
TIME_UTC="$(date -u +%H:%M)"

LAST_BODY=""
LAST_STATUS=""

log() {
  printf '==> %s\n' "$1"
}

request() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local auth_header="${4:-}"
  local body_file="$TMP_DIR/body.out"

  local -a curl_args=(
    -sS
    -o "$body_file"
    -w "%{http_code}"
    -X "$method"
    -H "Accept: application/json"
    -b "$COOKIE_JAR"
    -c "$COOKIE_JAR"
  )

  if [[ -n "$auth_header" ]]; then
    curl_args+=(-H "Authorization: Bearer $auth_header")
  fi

  if [[ -n "$body" ]]; then
    curl_args+=(-H "Content-Type: application/json" --data "$body")
  fi

  LAST_STATUS="$(curl "${curl_args[@]}" "$BASE_URL$path")"
  LAST_BODY="$(cat "$body_file")"
}

assert_status() {
  local expected="$1"
  if [[ "$LAST_STATUS" != "$expected" ]]; then
    printf 'Expected HTTP %s, got %s\n' "$expected" "$LAST_STATUS" >&2
    printf '%s\n' "$LAST_BODY" >&2
    exit 1
  fi
}

json_get() {
  local expr="$1"
  JSON_INPUT="$LAST_BODY" python3 - "$expr" <<'PY'
import json
import os
import sys

expr = sys.argv[1]
value = json.loads(os.environ["JSON_INPUT"])
for part in expr.split('.'):
    if not part:
        continue
    if isinstance(value, list):
        value = value[int(part)]
    else:
        value = value[part]
if value is None:
    print("")
elif isinstance(value, bool):
    print("true" if value else "false")
else:
    print(value)
PY
}

assert_json_equals() {
  local expr="$1"
  local expected="$2"
  local actual
  actual="$(json_get "$expr")"
  if [[ "$actual" != "$expected" ]]; then
    printf 'Expected JSON %s=%s, got %s\n' "$expr" "$expected" "$actual" >&2
    printf '%s\n' "$LAST_BODY" >&2
    exit 1
  fi
}

assert_json_nonempty() {
  local expr="$1"
  local actual
  actual="$(json_get "$expr")"
  if [[ -z "$actual" ]]; then
    printf 'Expected JSON %s to be non-empty\n' "$expr" >&2
    printf '%s\n' "$LAST_BODY" >&2
    exit 1
  fi
}

assert_json_array_contains() {
  local expr="$1"
  local key="$2"
  local expected="$3"

  if ! JSON_INPUT="$LAST_BODY" python3 - "$expr" "$key" "$expected" <<'PY'
import json
import os
import sys

expr, key, expected = sys.argv[1], sys.argv[2], sys.argv[3]
value = json.loads(os.environ["JSON_INPUT"])
for part in expr.split('.'):
    if not part:
        continue
    value = value[part]

for item in value:
    current = item.get(key)
    if current is None:
        continue
    if str(current) == expected:
        sys.exit(0)

sys.exit(1)
PY
  then

    printf 'Expected JSON array %s to contain %s=%s\n' "$expr" "$key" "$expected" >&2
    printf '%s\n' "$LAST_BODY" >&2
    exit 1
  fi
}

log "health"
request GET /health
assert_status 200
assert_json_equals status ok

log "login"
request POST /api/auth/login "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}"
assert_status 200
assert_json_equals tokenType Bearer
assert_json_equals userId "$(json_get userId)"
assert_json_nonempty accessToken
if ! grep -q 'trackstack_refresh' "$COOKIE_JAR"; then
  printf 'Expected refresh cookie after login\n' >&2
  exit 1
fi
TOKEN="$(json_get accessToken)"
USER_ID="$(json_get userId)"

log "session"
request GET /api/auth/session "" "$TOKEN"
assert_status 200
assert_json_equals userId "$USER_ID"
assert_json_nonempty sessionId

log "refresh"
request POST /api/auth/refresh
assert_status 200
assert_json_equals userId "$USER_ID"
assert_json_equals tokenType Bearer
assert_json_nonempty accessToken
REFRESHED_TOKEN="$(json_get accessToken)"

log "session after refresh"
request GET /api/auth/session "" "$REFRESHED_TOKEN"
assert_status 200
assert_json_equals userId "$USER_ID"
assert_json_nonempty sessionId

TOKEN="$REFRESHED_TOKEN"

log "heat create"
request POST /api/heat/refills '{"date":"2026-03-01","weightKg":12.5,"bags":4,"temperature":7.5}' "$TOKEN"
assert_status 201
assert_json_equals userId "$USER_ID"
assert_json_equals bags 4
HEAT_REFILL_ID="$(json_get id)"

log "heat list"
request GET /api/heat/refills?from=2026-01-01\&to=2026-12-31 "" "$TOKEN"
assert_status 200
assert_json_array_contains "" id "$HEAT_REFILL_ID"

log "heat delete"
request DELETE "/api/heat/refills/$HEAT_REFILL_ID" "" "$TOKEN"
assert_status 204

log "calories target upsert"
request POST /api/calories/target '{"targetCalories":2400,"targetProteinGrams":140,"targetCarbGrams":220,"targetFatGrams":75}' "$TOKEN"
assert_status 200
assert_json_equals userId "$USER_ID"
assert_json_equals targetCalories 2400

log "calories log create"
request POST /api/calories/log "{\"calories\":650,\"proteinGrams\":45,\"carbGrams\":70,\"fatGrams\":20,\"title\":\"E2E Meal\",\"date\":\"$TODAY_UTC\",\"time\":\"$TIME_UTC\"}" "$TOKEN"
assert_status 201
assert_json_equals userId "$USER_ID"
assert_json_equals title "E2E Meal"
CALORIES_LOG_ID="$(json_get id)"

log "calories dashboard"
request GET /api/calories/dashboard "" "$TOKEN"
assert_status 200
assert_json_array_contains logs id "$CALORIES_LOG_ID"
assert_json_equals summary.target.targetCalories 2400

log "calories delete"
request DELETE "/api/calories/logs/$CALORIES_LOG_ID" "" "$TOKEN"
assert_status 204

log "expenses settings read"
request GET /api/expenses/settings "" "$TOKEN"
assert_status 200
assert_json_equals settings.userId "$USER_ID"

log "expenses settings update"
request POST /api/expenses/settings '{"income":3000,"ratioFund":50,"ratioFun":30,"ratioFuture":20}' "$TOKEN"
assert_status 200
assert_json_equals income 3000
assert_json_equals ratioFund 50

log "expenses entry create"
request POST /api/expenses/entries '{"title":"E2E expense","amount":42.5,"category":"fun","date":"2026-03-03"}' "$TOKEN"
assert_status 201
assert_json_equals userId "$USER_ID"
assert_json_equals category fun
EXPENSE_ENTRY_ID="$(json_get id)"

log "expenses checklist create"
request POST /api/expenses/checklists '{"title":"E2E checklist","amount":12.5,"category":"fund"}' "$TOKEN"
assert_status 200
CHECKLIST_ID="$(json_get id)"

log "expenses recurring create"
request POST /api/expenses/recurring '{"title":"E2E recurring","amount":19.0,"category":"future"}' "$TOKEN"
assert_status 200
RECURRING_ID="$(json_get id)"

log "expenses dashboard"
request GET /api/expenses/sheet/current "" "$TOKEN"
assert_status 200
assert_json_array_contains history id "$EXPENSE_ENTRY_ID"
assert_json_array_contains pendingObligations templateId "$CHECKLIST_ID"

log "expenses delete entry"
request DELETE "/api/expenses/entries/$EXPENSE_ENTRY_ID" "" "$TOKEN"
assert_status 204

log "expenses delete checklist"
request DELETE "/api/expenses/checklists/$CHECKLIST_ID" "" "$TOKEN"
assert_status 204

log "expenses delete recurring"
request DELETE "/api/expenses/recurring/$RECURRING_ID" "" "$TOKEN"
assert_status 204

log "logout"
request POST /api/auth/logout "" "$TOKEN"
assert_status 204
if grep -q 'trackstack_refresh' "$COOKIE_JAR"; then
  printf 'Expected refresh cookie to be cleared on logout\n' >&2
  exit 1
fi

log "refresh after logout fails"
request POST /api/auth/refresh
assert_status 401

log "done"
printf 'All backend curl e2e checks passed.\n'
