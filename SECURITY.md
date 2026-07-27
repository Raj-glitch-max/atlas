# Security Policy

Atlas is an offline-verifiable cryptographic delegation runtime. A verification
library that gets security wrong fails silently, so we take reports seriously and
respond on a stated timeline.

Read this alongside [`THREAT_MODEL.md`](docs/architecture/THREAT_MODEL.md) (what Atlas defends and
how each claim is proven) and [`LIMITATIONS.md`](LIMITATIONS.md) (what Atlas does
**not** defend, stated plainly). Neither is a substitute for reporting something
new.

## Reporting a vulnerability

**Do not open a public issue, pull request, or discussion for a suspected
vulnerability.** Public disclosure before a fix puts every user at risk.

Report privately, by either channel:

- **GitHub Security Advisory** (preferred): open a private report at
  <https://github.com/Raj-glitch-max/atlas/security/advisories/new>.
- **Email**: <rpdinkar92260@gmail.com> with subject line `ATLAS SECURITY`.

Please include, as far as you can:

- the affected component (e.g. `internal/verify`, the JWS envelope, a snapshot),
- a description of the issue and its impact (forge / widen / replay / DoS / info-leak),
- a reproduction: a record, an input sequence, or a failing test/vector,
- the version or commit, and your environment.

A working proof-of-concept — a record that verifies and should not, or a
sequence of legal operations composing into an illegal capability — is the most
useful thing you can send.

## Our commitment (response SLA)

| Stage | Target |
|---|---|
| Acknowledge receipt | within **3 business days** |
| Initial assessment (severity + validity) | within **7 business days** |
| Fix or mitigation plan communicated | within **30 days** of a confirmed report |
| Public disclosure | coordinated with you, after a fix is available |

Atlas has a single maintainer today (see [`MAINTAINERS.md`](.github/MAINTAINERS.md)), so
these are honest targets rather than a staffed guarantee. We will keep you
informed if a timeline slips.

## Disclosure policy

We follow **coordinated disclosure**. We ask that you give us a reasonable window
(default **90 days**, or sooner once a fix ships) before public disclosure. We
will credit reporters who wish to be named. We will not pursue legal action
against good-faith research that respects this policy and does not access, modify,
or exfiltrate data beyond what is needed to demonstrate the issue.

There is no paid bug-bounty program at this time.

## What we consider in scope

In scope — anything that breaks a claim in [`THREAT_MODEL.md`](docs/architecture/THREAT_MODEL.md):

- A record that verifies but should be rejected (silent acceptance) — the worst
  class, always high severity.
- Algorithm confusion, signature forgery, or any envelope bypass in `internal/record`.
- Scope widening past the issuance-time strict-subset guard.
- Snapshot rollback that un-revokes, or acceptance of stale revocation knowledge
  beyond the policy bound `R`.
- Verification that performs I/O or otherwise breaks the offline/fail-closed
  contract.
- Cross-implementation conformance divergences (Go / Python / TypeScript) that a
  vector should catch.

Explicitly **out of scope** (documented trust boundaries — see `THREAT_MODEL.md`
§4 and `LIMITATIONS.md`): compromised issuer signing keys, compromised agent
runtime or host, SPIFFE/SPIRE compromise, within-window replay of a valid record
(carried openly as FM8 / `LIMITATIONS.md` §3), and business-policy /
authorization decisions layered above verification. A report against one of these
is welcome as a *documentation* issue if the boundary is unclear, but it is not a
vulnerability in the runtime.

## Supported versions

Atlas is pre-1.0 (`v0.x`, reference implementation). Security fixes are applied to
the `main` branch and the latest tagged release.

| Version | Supported |
|---|---|
| `main` (latest) | ✅ |
| latest `v0.x` tag | ✅ |
| older `v0.x` tags | ❌ (upgrade to latest) |

Once Atlas reaches `v1.0`, this table will list the maintained minor versions.

## Operational note

This repository's pre-commit and CI pipelines run private-key detection
(`detect-private-key`) and secret scanning (gitleaks, `.gitleaks.toml`). A
committed secret triggers CI failure; rotate the secret and rewrite history if
one is exposed — do not rely on deleting the commit alone.
