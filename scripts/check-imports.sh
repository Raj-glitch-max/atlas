#!/usr/bin/env bash
# check-imports.sh — dependency-rule lint for the Atlas module graph.
#
# Executable form of the architecture's dependency rules R1-R7
# (SYSTEM_ARCHITECTURE.md §2) and the per-module forbidden-import tables
# (MODULE_SPECIFICATION.md). A violation fails the build (AD-010): the
# load-bearing properties — offline verification, volatility isolation,
# no cross-domain coupling — must not depend on reviewer vigilance (FM11).
#
# Scope: direct imports of every package in this Go module. Module-internal
# transitivity is covered because every module package is checked; external
# behavioral guarantees (zero egress during verification) are asserted by
# AT16's instrumentation, not by this lint.
#
# Changing a rule here requires amending MODULE_SPECIFICATION.md with a
# frozen trace first (ENGINEERING_DECISION_RECORD.md closing rule).
#
# Usage:
#   scripts/check-imports.sh              lint the module
#   scripts/check-imports.sh --self-test  verify the lint's own logic

set -u

MODULE="github.com/Raj-glitch-max/atlas"
violations=0

# ---------------------------------------------------------------------------
# check_pkg <import-path> <comma-separated-direct-imports>
# Prints one line per violation; increments $violations.
# ---------------------------------------------------------------------------
check_pkg() {
    local pkg="$1"
    local imports_csv="$2"
    local imp

    IFS=',' read -ra imps <<< "$imports_csv"
    for imp in "${imps[@]}"; do
        [ -z "$imp" ] && continue

        # Global rule: no product package (internal/*, cmd/*) imports the
        # test tree (SR-4: fake leakage toward production wiring).
        case "$pkg" in
            "$MODULE"/internal/*|"$MODULE"/cmd/*)
                case "$imp" in
                    "$MODULE"/tests/*)
                        echo "VIOLATION: $pkg imports test tree $imp (SR-4)"
                        violations=$((violations + 1))
                        continue
                        ;;
                esac
                ;;
        esac

        case "$pkg" in
            # M1 record: depends on nothing in the module (R1); pure — no
            # network, no process execution, no SPIRE API surface.
            "$MODULE"/internal/record)
                case "$imp" in
                    "$MODULE"/*)
                        echo "VIOLATION: record imports module package $imp (R1)"
                        violations=$((violations + 1)) ;;
                    net|net/*|os/exec)
                        echo "VIOLATION: record imports $imp (M1 purity)"
                        violations=$((violations + 1)) ;;
                    github.com/spiffe/spire-api-sdk/*)
                        echo "VIOLATION: record imports $imp (M1 forbidden dependency)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M2 issuance: module imports limited to record (R2); no network.
            "$MODULE"/internal/issuance)
                case "$imp" in
                    "$MODULE"/internal/record) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: issuance imports module package $imp (R2/R3)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: issuance imports $imp (M2 forbidden dependency)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M3 verify: module imports limited to record (R3); structurally
            # offline (AP1) — no network; integrity goes through record, so
            # no direct jose; no SPIRE API surface.
            "$MODULE"/internal/verify)
                case "$imp" in
                    "$MODULE"/internal/record) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: verify imports module package $imp (R3)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: verify imports $imp (AP1/R6 structural offline)"
                        violations=$((violations + 1)) ;;
                    github.com/go-jose/*|gopkg.in/go-jose/*|gopkg.in/square/go-jose*)
                        echo "VIOLATION: verify imports $imp directly (integrity goes through record)"
                        violations=$((violations + 1)) ;;
                    github.com/spiffe/spire-api-sdk/*)
                        echo "VIOLATION: verify imports $imp (M3 forbidden dependency)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M4 truststore: module imports limited to record; structurally
            # incapable of fetching (FM9) — no network.
            "$MODULE"/internal/truststore)
                case "$imp" in
                    "$MODULE"/internal/record) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: truststore imports module package $imp (R2/R3)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: truststore imports $imp (FM9 never-fetch, structural)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M5 contracttest: exercises the revstatus contract, so it may
            # import revstatus and record; nothing else in the module; no
            # network. (Checked before the revstatus/* pattern below.)
            "$MODULE"/internal/revstatus/contracttest)
                case "$imp" in
                    "$MODULE"/internal/record|"$MODULE"/internal/revstatus) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: contracttest imports module package $imp (R2/R3)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: contracttest imports $imp (pre-mechanism network ban)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M5 revstatus: module imports limited to record; no network in
            # anything shipped pre-spike (a future realization's maintenance
            # path receives an explicit, file-scoped grant with the
            # mechanism decision — AD-D02 — by amending this lint).
            "$MODULE"/internal/revstatus)
                case "$imp" in
                    "$MODULE"/internal/record) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: revstatus imports module package $imp (R2/R3/R5)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: revstatus imports $imp (pre-mechanism network ban, R6)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # M6 revorigin: module imports limited to record; publication is
            # the deferred propagation channel's business — no network.
            "$MODULE"/internal/revorigin)
                case "$imp" in
                    "$MODULE"/internal/record) ;;
                    "$MODULE"/*)
                        echo "VIOLATION: revorigin imports module package $imp (R2/R5)"
                        violations=$((violations + 1)) ;;
                    net|net/*)
                        echo "VIOLATION: revorigin imports $imp (M6 forbidden dependency)"
                        violations=$((violations + 1)) ;;
                esac
                ;;

            # cmd/* composition roots: may wire any internal package; the
            # test-tree ban above still applies. tests/*: unconstrained
            # (instrumentation, not product).
        esac
    done
}

# ---------------------------------------------------------------------------
# Self-test: feed synthetic package/import pairs through the same rule
# function and assert both directions — violations are caught, legitimate
# imports pass. Guards the lint itself against silent decay.
# ---------------------------------------------------------------------------
self_test() {
    local failures=0

    expect_violations() { # <expected-count> <pkg> <imports-csv>
        local expected="$1" pkg="$2" csv="$3"
        violations=0
        check_pkg "$pkg" "$csv" > /dev/null
        if [ "$violations" -ne "$expected" ]; then
            echo "SELF-TEST FAIL: $pkg [$csv] expected $expected violation(s), got $violations"
            failures=$((failures + 1))
        fi
    }

    # Violations that must be caught.
    expect_violations 1 "$MODULE/internal/verify"     "net/http"
    expect_violations 1 "$MODULE/internal/verify"     "$MODULE/internal/revstatus"
    expect_violations 1 "$MODULE/internal/record"     "$MODULE/internal/verify"
    expect_violations 1 "$MODULE/internal/record"     "os/exec"
    expect_violations 1 "$MODULE/internal/truststore" "net"
    expect_violations 1 "$MODULE/internal/truststore" "$MODULE/tests/harness"
    expect_violations 1 "$MODULE/internal/revstatus"  "$MODULE/internal/revorigin"
    expect_violations 1 "$MODULE/cmd/atlas-verify"    "$MODULE/tests/harness"
    expect_violations 2 "$MODULE/internal/verify"     "net,gopkg.in/go-jose/go-jose.v3"

    # Legitimate imports that must pass.
    expect_violations 0 "$MODULE/internal/verify"     "$MODULE/internal/record,fmt,time"
    expect_violations 0 "$MODULE/internal/issuance"   "$MODULE/internal/record,errors"
    expect_violations 0 "$MODULE/internal/record"     "fmt,errors,strings"
    expect_violations 0 "$MODULE/internal/revstatus/contracttest" "$MODULE/internal/revstatus,$MODULE/internal/record,testing"
    expect_violations 0 "$MODULE/cmd/atlas-verify"    "$MODULE/internal/verify,$MODULE/internal/truststore,os"
    expect_violations 0 "$MODULE/tests/harness"       "net/http,os/exec"

    if [ "$failures" -ne 0 ]; then
        echo "check-imports self-test: FAIL ($failures case(s))"
        return 1
    fi
    echo "check-imports self-test: PASS"
    return 0
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
if [ "${1:-}" = "--self-test" ]; then
    self_test
    exit $?
fi

if ! command -v go > /dev/null 2>&1; then
    echo "check-imports: go toolchain not found" >&2
    exit 1
fi

total=0
while IFS='|' read -r pkg imports_csv; do
    [ -z "$pkg" ] && continue
    check_pkg "$pkg" "$imports_csv"
    total=$((total + 1))
done < <(go list -f '{{.ImportPath}}|{{join .Imports ","}}' ./... 2>&1)

if [ "$violations" -ne 0 ]; then
    echo "check-imports: FAIL — $violations violation(s) across $total package(s)"
    exit 1
fi
echo "check-imports: PASS ($total package(s), 0 violations)"
