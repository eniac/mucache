
PORT=5002
## If we don't use the --workers option, then it
## spawns workers equal to the $WEB_CONCURRENCY variable.
WORKERS=${1?Number of workers not given}
## We need to tweak this number to avoid any limits being hit
CONCURRENT_REQUESTS=10000
## Set the log level to info by default and to critical when in production
LOG_LEVEL="${2:-info}"

## --reload is only used when debugging
uvicorn frontend:app --port "${PORT}" --log-level "${LOG_LEVEL}" --workers "${WORKERS}" --limit-concurrency "${CONCURRENT_REQUESTS}"

