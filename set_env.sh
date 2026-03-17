#!/bin/bash

if [ -n "${BASH_VERSION:-}" ]; then
	SCRIPT_PATH="${BASH_SOURCE[0]}"
elif [ -n "${ZSH_VERSION:-}" ]; then
	SCRIPT_PATH="$(eval 'printf "%s" "${(%):-%x}"')"
else
	SCRIPT_PATH="$0"
fi

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$SCRIPT_PATH")" && pwd)" || {
	echo "failed to resolve script directory" >&2
	return 1 2>/dev/null || exit 1
}

ENV_FILE="$SCRIPT_DIR/.env"
if [ ! -f "$ENV_FILE" ]; then
	echo "environment file not found: $ENV_FILE" >&2
	return 1 2>/dev/null || exit 1
fi

set -a
. "$ENV_FILE"
set +a
