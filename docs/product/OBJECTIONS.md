# Objections — and how Atlas answers them

Every hard question a skeptic, a security reviewer, an investor, or a FAANG
interviewer will actually ask. Answers are honest: where Atlas is strong it
says so with evidence; where it's limited it points at `LIMITATIONS.md` rather
than bluffing. Bluffing is how you lose these rooms.

---

## Security objections

**"Your crypto is young — why should I trust it?"**
It isn't young; it's deliberately boring. Records are ES256/P-256 JWS via Go's
standard library — the same primitives in your TLS stack. Atlas invents no
cryptography. What's novel is the *verification semantics* (ordered checks,
fail-closed staleness, verifiable-freshness revocation), and those are pinned
by 28 language-neutral conformance vectors incl. 18 adversarial ones. Caveat
carried openly: no third-party audit yet (`LIMITATIONS.md` §9).

**"If someone steals a record, they can use it — it's a bearer token."**
Correct, within its window (FM8, `LIMITATIONS.md` §3). Mitigations today:
short TTLs, narrow scope, revocation. The real fix — proof-of-possession
binding the record to the delegate's key — is roadmapped. We state this
instead of hiding it; a reviewer who finds a hidden bearer problem stops
trusting everything else.

**"alg=none / algorithm confusion / signature transplant?"**
All refused, and it's a *test that fails the build* if any is accepted:
`tests/vectors/negative-vectors.json` has stripped-sig alg=none, HS256
key-confusion (public key as HMAC secret), ES384 substitution, cross-record
signature and payload transplants, forged kid, truncations. Plus a
coverage-guided fuzzer with zero silent acceptances over 800k+ executions.

**"Why should I trust your revocation freshness? Offline means stale."**
Yes — and that staleness is the *point made explicit*. The verifier's
`--max-staleness` budget R is its revocation blind spot, chosen by the
operator, and past it verification **fails closed** (`Inconclusive`→reject),
never open. A tampered bundle that strips revocations is refused because the
snapshot signature covers the revoked set (`THREAT_MODEL.md` C7–C8, both
tested). We don't claim to beat physics; we make the bound honest.

**"What if the issuer key is compromised?"**
Game over for *new* issuance, as with any signature system — bounded by short
expiry and by `kid`-pinned rotation (retire the key in the next bundle).
Automated rotation isn't shipped (`LIMITATIONS.md` §4). We don't pretend key
compromise is survivable; we bound its blast radius.

**"Partition behavior — can an attacker who controls the network forge an
ACCEPT?"** No. Network control buys *unavailability of fresh knowledge*, never
a false accept. In-partition revocations are invisible until recovery (INV12),
stated and tested. The failure mode is "fails closed," which is the safe one.

## "Isn't this just X?" objections

**"…just OAuth?"** OAuth = client↔server auth with a live server. Atlas =
agent↔agent delegation verifiable offline, with attenuation. They compose;
Atlas doesn't replace your IdP (`WHY.md` Q3).

**"…just JWT?"** JWT answers "who signed this." Atlas answers "is this
delegated authority still trustworthy here, now" — attenuation + expiry +
independent revocation + fail-closed staleness (Q4). A record *is* a JWS; the
value is the semantics around it.

**"…just Biscuit / Macaroons?"** They solved offline attenuation first, and we
say so. Atlas adds SPIFFE identity binding, verifiable-freshness revocation
with honest partition semantics, and a full runtime + conformance surface for
the agent workflow (Q5). If Biscuit fits, use it.

**"…just SPIFFE/SPIRE?"** SPIFFE issues *identity*. Atlas rides on top to add
*delegation between* those identities. It consumes SPIFFE, never replaces it
(INV10/INV11).

## Investor objections

**"Who actually buys this?"** AI-platform companies, enterprise AI teams, AI
infra vendors, and regulated buyers (fintech/defense/robotics) building
multi-agent systems that cross trust boundaries. First user, sharply: the AI
engineer wiring agents that delegate work.

**"Big Tech will just build it."** Probably — and that's the Kubernetes
pattern, not the death knell. Open standards win as the *neutral* layer
precisely because no single vendor owns them. Atlas's moat isn't secrecy;
it's the conformance suite + protocol maturity that make it the reference
everyone tests against.

**"Is it replaceable?"** The token isn't the moat; the accumulated
protocol/verifier/revocation/conformance/ecosystem work is. Replacing it means
redoing all of that *and* re-earning interop trust. That's expensive on
purpose.

**"Is the market real or hype?"** Evidence-based framing only: *if*
multi-agent systems keep growing (every framework roadmap bets they will),
secure delegation becomes infrastructure. We don't promise inevitability; we
position for the specific bet.

## Engineering / interview objections

**"Does it even run? Reproduce a result."** `make ci` is green (lint, docs,
frozen-doc integrity, import-boundary, ~20 packages of tests). `bash
examples/unforgettable.sh` runs the whole thesis live. Nothing here is a
screenshot.

**"Benchmarks — real or vibes?"** Measured against the real engine: verify
~94µs / ~10.6k/s, issue ~30µs, 403-byte proof, published with the harness and
a Prometheus latency histogram. Substrate-dependent numbers (partition,
cross-domain) are explicitly *not* claimed until the SPIRE lab runs on real
infra — `LIMITATIONS.md` §8.

**"Where's the boundary between demo and proof?"** Drawn explicitly in
`THREAT_MODEL.md` (every claim → the test that proves it) and `LIMITATIONS.md`
(what isn't proven yet). That boundary being *drawn at all* is the signal.
