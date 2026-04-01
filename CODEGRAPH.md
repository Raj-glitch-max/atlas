# CodeGraph Index — Atlas

This project uses CodeGraph to build a local, semantic index of the codebase
for AI coding agents (such as Claude Code and Cursor).

---

## How to Rebuild

To rebuild the index from scratch or sync changes, run:

```bash
# Rebuild the full index from scratch
codegraph index

# Or initialize / re-index
codegraph init

# Sync incremental changes since the last index
codegraph sync
```

---

## What is Indexed

* All source code files (Go, JavaScript/TypeScript, Python, etc.) once
  implementation begins.
* Structural metadata, AST symbols, function call trees, import statements,
  and local references.
* System configurations and workflows (YAML files under `.github/`).

---

## What is Excluded

* Any files or folders matching rules in the root `.gitignore` (such as
  `node_modules/`, `dist/`, `.env`, and build targets).
* Local agent cache and session state (`.claude/`, `.git/`, `.vscode/`).
* The local CodeGraph index database itself (`.codegraph/`).

<!-- checkpoint: repo(architecture-draft): extend architecture draft -->
