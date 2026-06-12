// Package conformance is the executable conformance kit and property/
// differential fuzzer for the Verification Core (M3).
//
// It exists because a trust primitive meant to have competing implementations
// dies of verifier differentials — implementations that disagree on edge
// cases, letting an attacker pick the one that accepts (the Frankencerts
// failure that has burned X.509, TLS, and JWT ecosystems). OMEGA-04 makes
// this precise: the verifier IS the discriminating channel, so if two
// verifiers disagree the observation is ill-defined and I(C;E) is not what a
// relying party thinks. Exact, machine-checked conformance is therefore a
// precondition for the trust guarantee to hold across implementations.
//
// Contents:
//   - conformance.go — the reusable kit: a Verifier interface, a Factory any
//     implementation supplies, the executable Scenario corpus (exhaustive
//     over the verdict space), and Run(t, factory) which any implementation
//     calls to prove itself conformant.
//   - conformance_test.go — runs the corpus against the V1 verifier.
//   - properties_test.go — property/differential fuzzing: the invariants that
//     must hold for all inputs (tamper never accepts, garbage never accepts,
//     Accept implies all stages pass, routing consistency, absent material
//     never accepts, determinism).
//
// This is test infrastructure: it lives under tests/, imports no product
// package's internals beyond their public surfaces, and is imported by no
// product package.
package conformance
