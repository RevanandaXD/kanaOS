#!/usr/bin/env bash
# Generate a self-contained SVG card with live repository statistics.
#
# Usage:
#   bash scripts/generate-stats.sh [output.svg]
#
# Git metrics (commits, contributors, branches, tags, LOC) are collected from
# the local clone. GitHub metrics (stars, forks, watchers, issues) are read via
# `gh api` when GH_TOKEN is available, and can be overridden with env vars:
#   STARS, FORKS, WATCHERS, OPEN_ISSUES
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT="${1:-$ROOT_DIR/public/stats.svg}"
REPO="${GITHUB_REPOSITORY:-RevanandaXD/kanaOS}"

cd "$ROOT_DIR"
mkdir -p "$(dirname "$OUTPUT")"

count_lines() {
  wc -l | tr -d '[:space:]'
}

github_value() {
  local key="$1" fallback="${2:-0}" value=""
  if [[ -n "${GH_TOKEN:-}" ]] && command -v gh >/dev/null 2>&1; then
    value="$(gh api "repos/$REPO" --jq ".$key" 2>/dev/null || true)"
  fi
  if [[ -z "$value" || "$value" == "null" ]]; then
    value="$fallback"
  fi
  printf '%s' "$value"
}

# --- Git-derived metrics -----------------------------------------------------
commits="$(git rev-list --count HEAD 2>/dev/null || printf '0')"
contributors="$( (git shortlog -sn HEAD 2>/dev/null || true) | count_lines )"
branches="$( ( (git branch -r 2>/dev/null || true) | grep -v 'HEAD' || true) | count_lines )"
tags="$( (git tag 2>/dev/null || true) | count_lines )"
last_commit="$(git log -1 --format=%cd --date=format:'%d %b %Y' 2>/dev/null || printf '-')"
first_commit="$(git log --reverse --format=%cd --date=format:'%d %b %Y' 2>/dev/null | head -n 1 || printf '-')"
loc="$(find app -type f \( -name '*.go' -o -name '*.svelte' -o -name '*.ts' -o -name '*.js' -o -name '*.css' \) \
  -not -path '*/node_modules/*' -not -path '*/.svelte-kit/*' -not -path '*/build/*' -print0 2>/dev/null \
  | xargs -0 -r cat 2>/dev/null | count_lines)"

# --- GitHub-derived metrics --------------------------------------------------
stars="${STARS:-$(github_value stargazers_count)}"
forks="${FORKS:-$(github_value forks_count)}"
watchers="${WATCHERS:-$(github_value subscribers_count)}"
open_issues="${OPEN_ISSUES:-$(github_value open_issues_count)}"

xml_escape() {
  sed -e 's/&/\&amp;/g' -e 's/</\&lt;/g' -e 's/>/\&gt;/g'
}

commits="$(printf '%s' "$commits" | xml_escape)"
contributors="$(printf '%s' "$contributors" | xml_escape)"
stars="$(printf '%s' "$stars" | xml_escape)"
forks="$(printf '%s' "$forks" | xml_escape)"
branches="$(printf '%s' "$branches" | xml_escape)"
tags="$(printf '%s' "$tags" | xml_escape)"
loc="$(printf '%s' "$loc" | xml_escape)"
last_commit="$(printf '%s' "$last_commit" | xml_escape)"
first_commit="$(printf '%s' "$first_commit" | xml_escape)"
repo_escaped="$(printf '%s' "$REPO" | xml_escape)"

cat > "$OUTPUT" <<SVG
<svg xmlns="http://www.w3.org/2000/svg" width="620" height="360" viewBox="0 0 620 360" role="img" aria-label="kanaOS repository statistics">
  <defs>
    <linearGradient id="accent" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="#8B5CF6"/>
      <stop offset="100%" stop-color="#38BDF8"/>
    </linearGradient>
    <linearGradient id="bg" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color="#0D1424"/>
      <stop offset="100%" stop-color="#090E1A"/>
    </linearGradient>
  </defs>

  <rect x="1" y="1" width="618" height="358" rx="18" fill="url(#bg)" stroke="#22304D"/>
  <rect x="24" y="26" width="4" height="52" rx="2" fill="url(#accent)"/>

  <text x="40" y="50" fill="#F8FAFC" font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif" font-size="27" font-weight="700">kanaOS</text>
  <text x="40" y="70" fill="#8B9AB5" font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif" font-size="12">REPOSITORY STATISTICS · github.com/${repo_escaped}</text>

  <g>
    <rect x="514" y="34" width="82" height="24" rx="12" fill="#0F2A24" stroke="#1E5D4C"/>
    <circle cx="529" cy="46" r="4" fill="#34D399">
      <animate attributeName="opacity" values="1;0.35;1" dur="2s" repeatCount="indefinite"/>
    </circle>
    <text x="540" y="50" fill="#6EE7B7" font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif" font-size="11" font-weight="600">AUTO-UPDATE</text>
  </g>

  <line x1="24" y1="94" x2="596" y2="94" stroke="#22304D"/>

  <g font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif">
    <g transform="translate(24,110)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${commits}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">COMMITS</text>
    </g>
    <g transform="translate(218,110)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${contributors}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">CONTRIBUTORS</text>
    </g>
    <g transform="translate(412,110)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${stars}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">STARS</text>
    </g>

    <g transform="translate(24,206)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${forks}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">FORKS</text>
    </g>
    <g transform="translate(218,206)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${branches}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">BRANCHES</text>
    </g>
    <g transform="translate(412,206)">
      <rect width="184" height="86" rx="12" fill="#121B2E" stroke="#233252"/>
      <text x="16" y="45" fill="url(#accent)" font-size="32" font-weight="700">${tags}</text>
      <text x="16" y="68" fill="#8B9AB5" font-size="11" font-weight="600" letter-spacing="1.2">TAGS</text>
    </g>
  </g>

  <line x1="24" y1="308" x2="596" y2="308" stroke="#22304D"/>
  <text x="24" y="329" fill="#74849F" font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif" font-size="11">${loc} lines of code · ${watchers} watchers · ${open_issues} open issues</text>
  <text x="24" y="347" fill="#74849F" font-family="Segoe UI, Roboto, Helvetica, Arial, sans-serif" font-size="11">First commit ${first_commit} · Last commit ${last_commit}</text>
</svg>
SVG

printf 'Stats written to %s\n' "$OUTPUT"
