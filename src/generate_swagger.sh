#!/bin/sh

set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)" || exit 1
TMP_DIR="$(mktemp -d)"

cleanup() {
	rm -rf "$TMP_DIR"
}

trap cleanup EXIT INT TERM

cd "$TMP_DIR"
printf 'module swagrunner\n\ngo 1.25.6\n' > go.mod

go run github.com/swaggo/swag/cmd/swag@v1.16.6 init \
	-g main.go \
	-d "$SCRIPT_DIR/cmd","$SCRIPT_DIR/internal/presentation" \
	-o "$SCRIPT_DIR/docs"
