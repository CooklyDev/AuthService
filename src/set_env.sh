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

ROOT_SCRIPT="$SCRIPT_DIR/../set_env.sh"
if [ ! -f "$ROOT_SCRIPT" ]; then
	echo "environment loader not found: $ROOT_SCRIPT" >&2
	return 1 2>/dev/null || exit 1
fi

. "$ROOT_SCRIPT"
