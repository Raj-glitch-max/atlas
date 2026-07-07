#!/usr/bin/env bash
# Frozen-docs integrity guard.
# Verifies that the planning documents in scripts/frozen-docs.list have not
# changed from the SHA-256 baseline in FROZEN.sha256. A mismatch halts the run.
# Editing the baseline to silence the alarm is the violation being guarded.
#
# Exit 0 = unchanged; exit 1 = missing baseline, missing file, or hash mismatch.
set -euo pipefail

cd "$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

LIST="scripts/frozen-docs.list"
BASELINE="FROZEN.sha256"

if [ ! -f "$BASELINE" ]; then
  echo "FAIL: $BASELINE missing. Run 'make frozen-baseline' to capture it." >&2
  exit 1
fi
if [ ! -f "$LIST" ]; then
  echo "FAIL: $LIST missing." >&2
  exit 1
fi

status=0
while IFS= read -r f; do
  [ -z "$f" ] && continue
  case "$f" in \#*) continue;; esac   # allow comments in the list
  if [ ! -f "$f" ]; then
    echo "MISSING: $f"
    status=1
    continue
  fi
  actual=$(sha256sum "$f" | awk '{print $1}')
  expected=$(awk -v p="$f" '$2==p {print $1}' "$BASELINE" | head -n1)
  if [ -z "$expected" ]; then
    echo "UNTRACKED: $f (in list, not in baseline)"
    status=1
    continue
  fi
  if [ "$actual" != "$expected" ]; then
    echo "CHANGED: $f"
    status=1
  fi
done < "$LIST"

if [ "$status" -ne 0 ]; then
  cat >&2 <<'MSG'

FAIL: one or more frozen planning documents changed or are missing.
      This is the planning-system-is-frozen guard (see CONTRIBUTING.md §4).
      To amend a frozen document you must:
        1. record a journal entry at agents/journal/<YYYY-MM-DD>-<slug>.md,
        2. add a dated, reasoned change note to the document itself,
        3. re-baseline with 'make frozen-baseline' and commit FROZEN.sha256.
      Editing FROZEN.sha256 to silence this alarm is the violation being guarded.
MSG
  exit 1
fi

echo "OK: frozen planning documents unchanged ($(grep -cvE '^\s*(#|$)' "$LIST") files)."
