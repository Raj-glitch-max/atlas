# Who wrote this

Atlas was built by one maintainer, using Claude Code as an implementation tool,
under a spec-first process: RFCs and interface specifications first, code
second, conformance vectors as the definition of correct.

I state that plainly because you would infer it anyway — 29 commits carry a
`Co-Authored-By: Claude` trailer — and because a project about verifiable trust
that is cagey about its own provenance has already failed its first test.

## Why the process exists

The structure in this repository is not there to look rigorous. It is there
because unconstrained AI-assisted implementation fails in a specific,
recognizable way: documentation drifts confidently away from the code, claims
accumulate that nobody re-checks, and planning artifacts multiply faster than
working software.

I know that failure mode because this repository had it. An audit in July 2026,
run against the actual code rather than the documentation, found published
benchmarks that did not reproduce, a "coverage-guided fuzzing" claim that CI
never actually ran, a conformance-vector count that was wrong in six places, and
two dozen markdown files at the root burying the engineering underneath process
ceremony. Those are fixed.

They are also the reason the guardrails exist: hash-pinned frozen documents, an
import-boundary lint, language-neutral conformance vectors, and a benchmark
script that records the machine it ran on. Each was added after something drifted
— not in anticipation of it.

## What I claim, and what I don't

I direct the design, I review what lands, and I am accountable for all of it.

What I will not tell you is that every line received equal human attention. That
is exactly the kind of unfalsifiable assurance this project tries not to trade
in, and you have no way to check it.

So do not take the review on faith. Take the machinery, which you can run
yourself in about two minutes:

- **139 tests** across 14 packages
- **30 language-neutral conformance vectors**, 20 of them adversarial —
  `alg:none`, HS256 key confusion, signature and payload transplants, forged
  key ids, truncation, duplicate JSON keys
- **coverage-guided fuzzing on every pull request**
  (`-fuzz=FuzzVerify -fuzztime=60s`). A local run of the same target reached
  ~1.9M mutated executions with zero panics and zero silent acceptances
- an **import-boundary lint** that fails the build rather than a review comment
- **benchmarks published with the CPU, Go version, and iteration counts** that
  produced them, regenerable with one script

```sh
make ci
```

That command is the honest answer to "who wrote this." Not the author field.

## What this does not mean

It does not mean the code is safe by virtue of being tested. No third party has
audited the cryptography. [`LIMITATIONS.md`](LIMITATIONS.md) lists what Atlas
deliberately does not do, and where it would fail outright.

AI assistance changes neither of those facts. It only means the distance between
what this project claims and what it actually does had to be closed by
machinery, because it could not be closed by trust.
