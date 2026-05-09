#!/usr/bin/env bash
set -euo pipefail

skip_frontend=0
skip_admin=0
skip_client=0
skip_backend=0
skip_tests=0

for arg in "$@"; do
  case "$arg" in
    --skip-frontend) skip_frontend=1 ;;
    --skip-admin) skip_admin=1 ;;
    --skip-client) skip_client=1 ;;
    --skip-backend) skip_backend=1 ;;
    --skip-tests) skip_tests=1 ;;
    *)
      echo "Unknown argument: $arg" >&2
      exit 1
      ;;
  esac
done

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"

ensure_command() {
  local name="$1"
  command -v "$name" >/dev/null 2>&1 || { echo "Required command '$name' not found in PATH." >&2; exit 1; }
}

detect_node_tool() {
  if command -v pnpm >/dev/null 2>&1; then
    echo "pnpm"
    return
  fi
  if command -v npm >/dev/null 2>&1; then
    echo "npm"
    return
  fi
  echo "Neither 'pnpm' nor 'npm' found. Please install one of them." >&2
  exit 1
}

build_frontend() {
  local dir="$1"
  local name="$2"
  local tool="$3"

  echo "Building frontend ($name)"
  pushd "$dir" >/dev/null
  if [[ "$tool" == "pnpm" ]]; then
    echo "Using pnpm"
    if [[ -f "pnpm-lock.yaml" ]]; then
      pnpm install --frozen-lockfile
    else
      pnpm install
    fi
    pnpm run build
  else
    echo "Using npm"
    if [[ -f "package-lock.json" ]]; then
      npm ci
    else
      npm install
    fi
    npm run build
  fi
  popd >/dev/null
}

echo "Repo: $repo_root"

admin_dir="$repo_root/frontend/admin"
client_dir="$repo_root/frontend/client"
bin_dir="$repo_root/bin"
backend_out="$bin_dir/mou1ght"

if [[ $skip_frontend -eq 0 ]]; then
  node_tool="$(detect_node_tool)"

  if [[ $skip_admin -eq 0 ]]; then
    [[ -d "$admin_dir" ]] || { echo "Admin frontend directory not found: $admin_dir" >&2; exit 1; }
    build_frontend "$admin_dir" "frontend/admin" "$node_tool"
  fi

  if [[ $skip_client -eq 0 ]]; then
    if [[ -d "$client_dir" ]]; then
      build_frontend "$client_dir" "frontend/client" "$node_tool"
    else
      echo "Skipping frontend/client (directory not found)"
    fi
  fi
else
  echo "Skipping all frontend builds"
fi

if [[ $skip_backend -eq 0 ]]; then
  echo "Building backend (Go)"
  ensure_command go
  mkdir -p "$bin_dir"
  pushd "$repo_root" >/dev/null
  if [[ $skip_tests -eq 0 ]]; then
    echo "Running go test ./..."
    go test ./...
  fi
  go build -o "$backend_out" ./cmd
  popd >/dev/null
  echo "Backend output: $backend_out"
else
  echo "Skipping backend build"
fi

echo "Build complete."
