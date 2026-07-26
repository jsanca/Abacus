#!/bin/sh
set -eu

cleanup() {
  docker compose down --remove-orphans
}

diagnose_readiness_failure() {
  echo "Abacus stack did not become ready." >&2
  docker compose ps >&2 || true
  docker compose logs server >&2 || true
  docker compose logs client >&2 || true
}

wait_for_url() {
  url="$1"
  attempt=1
  max_attempts=45
  while [ "$attempt" -le "$max_attempts" ]; do
    if curl --fail --silent --show-error "$url" >/dev/null; then
      return 0
    fi
    attempt=$((attempt + 1))
    sleep 1
  done
  diagnose_readiness_failure
  return 1
}

trap cleanup EXIT INT TERM

docker compose down --remove-orphans
docker compose up --build -d
wait_for_url http://127.0.0.1:8080/health
wait_for_url http://127.0.0.1:3000/

set +e
(
  cd client
  npm run test:e2e
)
playwright_status=$?
set -e

exit "$playwright_status"
