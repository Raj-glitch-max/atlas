# tests/

This directory will hold automated tests. It is intentionally **empty of any stack-specific test runner** because the project has not yet selected a language/stack — that selection is deferred until after the C4 feasibility spike (see `../CONTRIBUTING.md` §5 and `../lab/EXP-001-EXECUTION-PLAN.md`).

When a stack is chosen:

- add the stack's test-runner configuration here (or in the language's conventional location),
- replace the placeholder `test` target in `../Makefile` with the real one,
- add a build/test job to `../.github/workflows/`.

Until then, `make test` is a documented no-op placeholder asserting that no stack is selected.
